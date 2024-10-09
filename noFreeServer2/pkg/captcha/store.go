package captcha

import (
	"time"

	"github.com/mojocn/base64Captcha"
)

type CaptchaStore interface {
	Get() base64Captcha.Store
	Set(id string, value string, expiration time.Duration) error
	Delete(id string) error
	Verify(id string, value string) (bool, error)
}
