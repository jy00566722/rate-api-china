package main

import (
	"mihu007/config"
	"mihu007/internal/handler"
	"mihu007/internal/middleware"
	"mihu007/internal/repository"
	"mihu007/internal/service"
	"mihu007/pkg/database"
	"mihu007/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db := database.InitMySQL(cfg.Database)

	// 初始化存储层
	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)

	// 初始化服务层
	jwtUtil := utils.NewJWTUtil(cfg.JWT)
	userService := service.NewUserService(userRepo, jwtUtil)
	deviceService := service.NewDeviceService(deviceRepo)
	membershipService := service.NewMembershipService(membershipRepo)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	deviceHandler := handler.NewDeviceHandler(deviceService)
	membershipHandler := handler.NewMembershipHandler(membershipService)

	// 设置路由
	r := gin.Default()

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/user/register", userHandler.Register)
		public.POST("/user/login", userHandler.Login)
		public.POST("/user/password-reset-code", userHandler.SendPasswordResetVerifyCode)
		public.POST("/user/reset-password", userHandler.ResetPassword)
	}

	// 需要认证的路由
	authorized := r.Group("/api/v1")
	authorized.Use(middleware.AuthMiddleware(jwtUtil))
	{
		authorized.GET("/user/info", userHandler.GetUserInfo)
		authorized.POST("/device/register", deviceHandler.Register)
		authorized.GET("/membership/info", membershipHandler.GetMembershipInfo)
		authorized.GET("/membership/plans", membershipHandler.GetMembershipPlans)
		authorized.POST("/membership/purchase", membershipHandler.PurchaseMembership)
	}

	// 启动服务器
	r.Run(cfg.Server.Port)
}
