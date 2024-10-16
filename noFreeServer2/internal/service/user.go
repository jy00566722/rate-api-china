package service

import (
	"context"
	"errors"
	"time"

	"mihu007/internal/model"
	"mihu007/internal/repository"
	"mihu007/pkg/captcha"
	"mihu007/pkg/email"
	"mihu007/pkg/sms"
	"mihu007/pkg/utils"
	"mihu007/pkg/wechat"

	"github.com/mojocn/base64Captcha"
)

type UserService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.User, string, error)
	ResetPassword(ctx context.Context, req model.PasswordResetRequest) error
	SendPasswordResetVerifyCode(ctx context.Context, req model.PasswordResetVerifyRequest) error
	WechatLogin(ctx context.Context, code string) (*model.User, string, error)
	WechatRegister(ctx context.Context, wechatInfo model.WechatInfo) (*model.User, error)
	GenerateCaptcha(ctx context.Context) (string, string, string, error)
	VerifyCaptcha(ctx context.Context, captchaID string, captchaCode string) (bool, error)
}

type userService struct {
	userRepo       repository.UserRepository
	membershipRepo repository.MembershipRepository
	verifyCodeRepo repository.VerifyCodeRepository
	jwtUtil        utils.JWTUtil
	wechatClient   wechat.WechatClient
	captchaStore   captcha.CaptchaStore
	emailSender    email.EmailSender
	smsSender      sms.SMSSender
}

func NewUserService(
	userRepo repository.UserRepository,
	membershipRepo repository.MembershipRepository,
	verifyCodeRepo repository.VerifyCodeRepository,
	jwtUtil utils.JWTUtil,
	wechatClient wechat.WechatClient,
	captchaStore captcha.CaptchaStore,
	emailSender email.EmailSender,
	smsSender sms.SMSSender,
) UserService {
	return &userService{
		userRepo:       userRepo,
		membershipRepo: membershipRepo,
		verifyCodeRepo: verifyCodeRepo,
		jwtUtil:        jwtUtil,
		wechatClient:   wechatClient,
		captchaStore:   captchaStore,
		emailSender:    emailSender,
		smsSender:      smsSender,
	}
}

func (s *userService) Register(ctx context.Context, req model.RegisterRequest) (*model.User, error) {
	// 1. 验证用户名是否已存在
	existingUser, _ := s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 2. 验证邮箱或手机号是否已存在
	if req.Email != "" {
		if existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email); existingUser != nil {
			return nil, errors.New("email already registered")
		}
	}
	if req.Phone != "" {
		if existingUser, _ := s.userRepo.GetByPhone(ctx, req.Phone); existingUser != nil {
			return nil, errors.New("phone already registered")
		}
	}

	// 3. 验证验证码
	var verifyID string
	if req.Email != "" {
		verifyID = req.Email
	} else if req.Phone != "" {
		verifyID = req.Phone
	} else {
		return nil, errors.New("email or phone is required")
	}

	isValid, err := s.verifyCodeRepo.Verify(ctx, verifyID, req.VerifyCode)
	if err != nil || !isValid {
		return nil, errors.New("invalid verification code")
	}

	// 4. 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 5. 创建用户
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		Password:     hashedPassword,
		RegisterType: getRegisterType(req),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// 6. 创建免费会员资格
	membership := &model.Membership{
		UserID:   user.ID,
		Level:    model.MemberLevelFree,
		ExpireAt: time.Now().AddDate(0, 1, 0), // 一个月有效期
		// 免费会员只能使用1个设备
	}
	err = s.membershipRepo.Create(ctx, membership)
	if err != nil {
		// 如果创建会员资格失败，需要回滚用户创建
		s.userRepo.Delete(ctx, user.ID)
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, req model.LoginRequest) (*model.User, string, error) {
	// 1. 查找用户
	var user *model.User
	var err error

	// 支持使用用户名、邮箱或手机号登录
	if utils.IsEmail(req.Username) {
		user, err = s.userRepo.GetByEmail(ctx, req.Username)
	} else if utils.IsPhone(req.Username) {
		user, err = s.userRepo.GetByPhone(ctx, req.Username)
	} else {
		user, err = s.userRepo.GetByUsername(ctx, req.Username)
	}

	if err != nil {
		return nil, "", errors.New("user not found")
	}

	// 2. 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, "", errors.New("invalid password")
	}

	// 3. 生成 JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// 辅助函数
func getRegisterType(req model.RegisterRequest) string {
	if req.Email != "" {
		return "email"
	}
	if req.Phone != "" {
		return "phone"
	}
	return "username"
}

func (s *userService) ResetPassword(ctx context.Context, req model.PasswordResetRequest) error {
	var user *model.User
	var err error

	if req.Email != "" {
		user, err = s.userRepo.GetByEmail(ctx, req.Email)
	} else if req.Phone != "" {
		user, err = s.userRepo.GetByPhone(ctx, req.Phone)
	} else {
		return errors.New("email or phone is required")
	}

	if err != nil {
		return err
	}

	// 验证验证码
	valid, err := s.verifyCodeRepo.Verify(ctx, user.Username, req.VerifyCode)
	if err != nil || !valid {
		return errors.New("invalid verify code")
	}

	// 更新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, user.ID, hashedPassword)
}
func (s *userService) SendPasswordResetVerifyCode(ctx context.Context, req model.PasswordResetVerifyRequest) error {
	var user *model.User
	var err error

	if req.Email != "" {
		user, err = s.userRepo.GetByEmail(ctx, req.Email)
	} else if req.Phone != "" {
		user, err = s.userRepo.GetByPhone(ctx, req.Phone)
	} else {
		return errors.New("email or phone is required")
	}

	if err != nil {
		return err
	}

	// 生成验证码
	code := utils.GenerateVerifyCode(6)

	// 存储验证码
	err = s.verifyCodeRepo.Save(ctx, user.Username, code, time.Now().Add(15*time.Minute))
	if err != nil {
		return err
	}

	// 发送验证码
	if req.Email != "" {
		return s.emailSender.SendPasswordResetCode(req.Email, code)
	} else {
		return s.smsSender.SendPasswordResetCode(req.Phone, code)
	}
}

func (s *userService) GenerateCaptcha(ctx context.Context) (string, string, string, error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, s.captchaStore.Get())
	id, b64s, answer, err := captcha.Generate()
	return id, b64s, answer, err
}

func (s *userService) VerifyCaptcha(ctx context.Context, captchaID string, captchaCode string) (bool, error) {
	return s.captchaStore.Verify(captchaID, captchaCode)
}

func (s *userService) WechatLogin(ctx context.Context, code string) (*model.User, string, error) {
	// 在此实现微信登录的逻辑
	return nil, "", nil // 临时返回nil、空字符串和nil，请根据实际需求修改
}

func (s *userService) WechatRegister(ctx context.Context, wechatInfo model.WechatInfo) (*model.User, error) {
	// 在此实现微信注册的逻辑
	return nil, nil // 临时返回nil和nil，请根据实际需求修改
}
