package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer interface {
	SendOTP(to, name, otp, purpose string) error
}

type Config struct {
	Driver   string // "smtp" or "console"
	From     string
	Host     string
	Port     string
	Username string
	Password string
}

func New(cfg Config) Mailer {
	if cfg.Driver == "smtp" {
		return &smtpMailer{cfg: cfg}
	}
	return &consoleMailer{from: cfg.From}
}

// consoleMailer — logs email to stdout, used in development
type consoleMailer struct {
	from string
}

func (m *consoleMailer) SendOTP(to, name, otp, purpose string) error {
	fmt.Printf("[MAILER] To: %s | Name: %s | OTP: %s | Purpose: %s\n", to, name, otp, purpose)
	return nil
}

// smtpMailer — sends real emails via SMTP
type smtpMailer struct {
	cfg Config
}

func (m *smtpMailer) SendOTP(to, name, otp, purpose string) error {
	subject := otpSubject(purpose)
	body := fmt.Sprintf(
		"Hi %s,\n\nYour verification code is: %s\n\nThis code will expire in 10 minutes.\n\nIf you did not request this, please ignore this email.",
		name, otp,
	)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		m.cfg.From, to, subject, body,
	))
	addr := fmt.Sprintf("%s:%s", m.cfg.Host, m.cfg.Port)
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	return smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg)
}

func otpSubject(purpose string) string {
	switch purpose {
	case "EMAIL_VERIFICATION":
		return "Email Verification Code"
	case "PASSWORD_RESET":
		return "Password Reset Code"
	default:
		return "Verification Code"
	}
}
