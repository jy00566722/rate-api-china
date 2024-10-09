package wechat

import (
	"mihu007/internal/model"
)

type WechatClient interface {
	GetAccessToken(code string) (string, error)
	GetUserInfo(accessToken, openID string) (*model.WechatInfo, error)
}
