package service

import (
	"auth-go/internal/config"
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendWelcomeEmail(toEmail string, name string) error
	SendResetPasswordEmail(toEmail string, resetLink string) error
}

type emailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{cfg}
}

func (s *emailService) SendWelcomeEmail(toEmail string, name string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Welcome to Auth Go!")
	m.SetBody("text/html", fmt.Sprintf("<h1>Hello %s!</h1><p>Welcome to our platform. We are glad to have you.</p>", name))

	d := gomail.NewDialer(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPEmail, s.cfg.SMTPPassword)
	return d.DialAndSend(m)
}

func (s *emailService) SendResetPasswordEmail(toEmail string, resetLink string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Reset Your Password")
	m.SetBody("text/html", fmt.Sprintf("<p>Click <a href='%s'>here</a> to reset your password.</p>", resetLink))

	d := gomail.NewDialer(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPEmail, s.cfg.SMTPPassword)
	return d.DialAndSend(m)
}
