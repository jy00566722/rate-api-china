package email

type EmailSender interface {
	SendPasswordResetCode(to string, subject string) error
}
