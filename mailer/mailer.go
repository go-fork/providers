package mailer

import (
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

// Mailer interface định nghĩa các phương thức để gửi email
type Mailer interface {
	// NewMessage tạo một message mới
	NewMessage() *Message

	// Send gửi một email đã được cấu hình
	Send(message *Message) error

	// Raw gửi một email với gomail.Message
	Raw(message *gomail.Message) error
}

// smtpMailer triển khai Mailer interface sử dụng SMTP
type smtpMailer struct {
	config *Config
	dialer *gomail.Dialer
}

// NewSMTPMailer tạo một mailer mới với cấu hình cho SMTP
func NewSMTPMailer(config *Config) *smtpMailer {
	if config.SMTP == nil {
		config.SMTP = &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			Encryption:  "none",
			FromAddress: "no-reply@example.com",
			FromName:    "System Notification",
			Timeout:     10 * time.Second,
		}
	}

	dialer := gomail.NewDialer(config.SMTP.Host, config.SMTP.Port, config.SMTP.Username, config.SMTP.Password)

	// Thiết lập mã hóa
	switch config.SMTP.Encryption {
	case "ssl", "tls":
		dialer.SSL = true
	case "none", "":
		dialer.SSL = false
	}

	// Thiết lập TLS config
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.SMTP.Host,
	}

	return &smtpMailer{
		config: config,
		dialer: dialer,
	}
}

// NewMessage tạo một message mới
func (m *smtpMailer) NewMessage() *Message {
	return NewMessage()
}

// Send gửi một email đã được cấu hình thông qua Message
func (m *smtpMailer) Send(message *Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}

	msg, err := message.BuildGoMailMessage(m.config.SMTP.FromAddress, m.config.SMTP.FromName)
	if err != nil {
		return err
	}
	return m.Raw(msg)
}

// Raw gửi một email trực tiếp với gomail.Message
func (m *smtpMailer) Raw(message *gomail.Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}

	sender, err := m.dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer sender.Close()

	return gomail.Send(sender, message)
}
