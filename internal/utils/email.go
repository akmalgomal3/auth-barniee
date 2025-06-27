package utils

import (
	"auth-barniee/internal/config"
	"fmt"
	"net/smtp"
)

func SendEmail(cfg *config.Config, to, subject, body string) error {
	from := cfg.SenderEmail
	password := cfg.SMTPPassword
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort

	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
