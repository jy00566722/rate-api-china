package sms

type SMSSender interface {
	SendPasswordResetCode(to string, code string) error
}
