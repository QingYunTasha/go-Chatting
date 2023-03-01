package infra

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type SMTPMailer struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewSMTPMailer(host string, port int, username string, password string) *SMTPMailer {
	return &SMTPMailer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (m *SMTPMailer) SendPasswordResetEmail(to, token string) error {
	// Set up email message
	from := "your-email@example.com"
	subject := "Password reset request"
	body := fmt.Sprintf("Hello,\n\nTo reset your password, please click the following link:\n\nhttps://your-website.com/reset-password?token=%s\n\nIf you did not request this password reset, please ignore this email.\n\nThanks,\nThe Your Website Team", token)
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	// Set up SMTP connection
	dialer := gomail.NewDialer(m.Host, m.Port, m.Username, m.Password)

	// Send email
	err := dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
