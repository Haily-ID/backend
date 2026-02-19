package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/haily-id/engine/internal/pkg/i18n"
)

type Mailer interface {
	SendOTP(to, name, otp, purpose, lang string) error
}

type Config struct {
	Driver   string // "smtp" or "console"
	FromName string
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
	return &consoleMailer{cfg: cfg}
}

// consoleMailer — logs email to stdout, used in development
type consoleMailer struct {
	cfg Config
}

func (m *consoleMailer) SendOTP(to, name, otp, purpose, lang string) error {
	content := i18n.OTPEmail(name, otp, purpose, lang)
	fmt.Printf("[MAILER] To: %s | Lang: %s | Subject: %s | OTP: %s\n", to, lang, content.Subject, otp)
	return nil
}

// smtpMailer — sends real emails via SMTP
type smtpMailer struct {
	cfg Config
}

func (m *smtpMailer) SendOTP(to, name, otp, purpose, lang string) error {
	content := i18n.OTPEmail(name, otp, purpose, lang)

	fromHeader := fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.From)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		fromHeader, to, content.Subject, content.Body,
	))

	addr := fmt.Sprintf("%s:%s", m.cfg.Host, m.cfg.Port)

	// Port 465 uses implicit TLS (SMTPS)
	if m.cfg.Port == "465" {
		return m.sendSSL(addr, to, msg)
	}

	// Port 587 uses STARTTLS
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	return smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg)
}

func (m *smtpMailer) sendSSL(addr, to string, msg []byte) error {
	tlsCfg := &tls.Config{ServerName: m.cfg.Host}
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return fmt.Errorf("failed to dial TLS: %w", err)
	}

	client, err := smtp.NewClient(conn, m.cfg.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}
	if err := client.Mail(m.cfg.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return w.Close()
}
