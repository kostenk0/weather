package api

import (
	"fmt"
	"net/smtp"
	"weather/internal/config"
)

func Send(to, subject, body string) error {
	cfg := config.App.Email

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		cfg.EmailFrom, to, subject, body)

	addr := cfg.SMTPHost + ":" + cfg.SMTPPort

	return smtp.SendMail(addr, auth, cfg.SMTPUser, []string{to}, []byte(msg))
}
