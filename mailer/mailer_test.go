package mailer

import (
	"errors"
	"io"
	"testing"
	"time"

	"gopkg.in/gomail.v2"
)

func TestNewSMTPMailer(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:        "smtp.example.com",
			Port:        587,
			Username:    "user@example.com",
			Password:    "password123",
			Encryption:  "tls",
			FromAddress: "test@example.com",
			FromName:    "Test Sender",
			Timeout:     30 * time.Second,
		},
	}

	mailer := NewSMTPMailer(config)

	if mailer == nil {
		t.Fatal("NewSMTPMailer() returned nil")
	}

	if mailer.config != config {
		t.Error("Expected config to match input config")
	}

	if mailer.dialer == nil {
		t.Fatal("Dialer was not initialized")
	}

	// Test dialer configuration
	if mailer.dialer.Host != "smtp.example.com" {
		t.Errorf("Expected dialer host to be 'smtp.example.com', got %s", mailer.dialer.Host)
	}

	if mailer.dialer.Port != 587 {
		t.Errorf("Expected dialer port to be 587, got %d", mailer.dialer.Port)
	}

	if mailer.dialer.Username != "user@example.com" {
		t.Errorf("Expected dialer username to be 'user@example.com', got %s", mailer.dialer.Username)
	}

	if mailer.dialer.Password != "password123" {
		t.Errorf("Expected dialer password to be 'password123', got %s", mailer.dialer.Password)
	}

	// For TLS encryption, SSL should be true
	if !mailer.dialer.SSL {
		t.Error("Expected SSL to be true for TLS encryption")
	}
}

func TestNewSMTPMailer_NilSMTPConfig(t *testing.T) {
	config := &Config{SMTP: nil}
	mailer := NewSMTPMailer(config)

	if mailer == nil {
		t.Fatal("NewSMTPMailer() returned nil")
	}

	// Should use default SMTP config
	if mailer.config.SMTP == nil {
		t.Fatal("SMTP config was not initialized")
	}

	expectedDefaults := &SMTPConfig{
		Host:        "localhost",
		Port:        25,
		Encryption:  "none",
		FromAddress: "no-reply@example.com",
		FromName:    "System Notification",
		Timeout:     10 * time.Second,
	}

	if mailer.config.SMTP.Host != expectedDefaults.Host {
		t.Errorf("Expected default host to be '%s', got %s", expectedDefaults.Host, mailer.config.SMTP.Host)
	}

	if mailer.config.SMTP.Port != expectedDefaults.Port {
		t.Errorf("Expected default port to be %d, got %d", expectedDefaults.Port, mailer.config.SMTP.Port)
	}

	if mailer.config.SMTP.Encryption != expectedDefaults.Encryption {
		t.Errorf("Expected default encryption to be '%s', got %s", expectedDefaults.Encryption, mailer.config.SMTP.Encryption)
	}

	if mailer.config.SMTP.FromAddress != expectedDefaults.FromAddress {
		t.Errorf("Expected default from address to be '%s', got %s", expectedDefaults.FromAddress, mailer.config.SMTP.FromAddress)
	}

	if mailer.config.SMTP.FromName != expectedDefaults.FromName {
		t.Errorf("Expected default from name to be '%s', got %s", expectedDefaults.FromName, mailer.config.SMTP.FromName)
	}

	if mailer.config.SMTP.Timeout != expectedDefaults.Timeout {
		t.Errorf("Expected default timeout to be %v, got %v", expectedDefaults.Timeout, mailer.config.SMTP.Timeout)
	}
}

func TestNewSMTPMailer_SSLEncryption(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:       "smtp.example.com",
			Port:       465,
			Encryption: "ssl",
		},
	}

	mailer := NewSMTPMailer(config)

	if !mailer.dialer.SSL {
		t.Error("Expected SSL to be true for SSL encryption")
	}
}

func TestNewSMTPMailer_NoEncryption(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:       "smtp.example.com",
			Port:       25,
			Encryption: "none",
		},
	}

	mailer := NewSMTPMailer(config)

	if mailer.dialer.SSL {
		t.Error("Expected SSL to be false for no encryption")
	}
}

func TestNewSMTPMailer_EmptyEncryption(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:       "smtp.example.com",
			Port:       25,
			Encryption: "",
		},
	}

	mailer := NewSMTPMailer(config)

	if mailer.dialer.SSL {
		t.Error("Expected SSL to be false for empty encryption")
	}
}

func TestNewSMTPMailer_TLSConfig(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:       "smtp.example.com",
			Port:       587,
			Encryption: "tls",
		},
	}

	mailer := NewSMTPMailer(config)

	if mailer.dialer.TLSConfig == nil {
		t.Fatal("TLS config was not initialized")
	}

	if !mailer.dialer.TLSConfig.InsecureSkipVerify {
		t.Error("Expected InsecureSkipVerify to be true")
	}

	if mailer.dialer.TLSConfig.ServerName != "smtp.example.com" {
		t.Errorf("Expected ServerName to be 'smtp.example.com', got %s", mailer.dialer.TLSConfig.ServerName)
	}
}

func TestSMTPMailer_NewMessage(t *testing.T) {
	config := &Config{SMTP: &SMTPConfig{}}
	mailer := NewSMTPMailer(config)

	message := mailer.NewMessage()

	if message == nil {
		t.Fatal("NewMessage() returned nil")
	}

	// Verify it's a properly initialized Message
	if message.to == nil {
		t.Error("Message 'to' slice was not initialized")
	}

	if message.cc == nil {
		t.Error("Message 'cc' slice was not initialized")
	}

	if message.bcc == nil {
		t.Error("Message 'bcc' slice was not initialized")
	}
}

func TestSMTPMailer_Send_NilMessage(t *testing.T) {
	config := &Config{SMTP: &SMTPConfig{}}
	mailer := NewSMTPMailer(config)

	err := mailer.Send(nil)

	if err == nil {
		t.Fatal("Expected error when sending nil message")
	}

	if err.Error() != "message cannot be nil" {
		t.Errorf("Expected 'message cannot be nil', got %s", err.Error())
	}
}

func TestSMTPMailer_Send_InvalidMessage(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
	}
	mailer := NewSMTPMailer(config)

	// Create an invalid message (no recipients, no content)
	message := NewMessage()

	err := mailer.Send(message)

	if err == nil {
		t.Fatal("Expected error when sending invalid message")
	}

	// The error should come from message validation
	if err.Error() != "email must have at least one recipient (to, cc, or bcc)" {
		t.Errorf("Expected validation error, got %s", err.Error())
	}
}

func TestSMTPMailer_Send_ValidMessage(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
	}
	mailer := NewSMTPMailer(config)

	// Create a valid message
	message := NewMessage().
		To("recipient@example.com").
		Subject("Test Subject").
		Text("Test Body")

	// Note: This will fail with a real SMTP connection error,
	// but we're testing that the message is properly built and passed to Raw()
	err := mailer.Send(message)

	// We expect a connection error since we're not running a real SMTP server
	if err == nil {
		t.Fatal("Expected connection error")
	}

	// The error should be about connection, not about message format
	if err.Error() == "message cannot be nil" {
		t.Error("Message was nil when it shouldn't be")
	}
}

func TestSMTPMailer_Raw_NilMessage(t *testing.T) {
	config := &Config{SMTP: &SMTPConfig{}}
	mailer := NewSMTPMailer(config)

	err := mailer.Raw(nil)

	if err == nil {
		t.Fatal("Expected error when sending nil message")
	}

	if err.Error() != "message cannot be nil" {
		t.Errorf("Expected 'message cannot be nil', got %s", err.Error())
	}
}

func TestSMTPMailer_Raw_ValidMessage(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host: "localhost",
			Port: 25,
		},
	}
	mailer := NewSMTPMailer(config)

	// Create a gomail message
	msg := gomail.NewMessage()
	msg.SetHeader("From", "test@example.com")
	msg.SetHeader("To", "recipient@example.com")
	msg.SetHeader("Subject", "Test")
	msg.SetBody("text/plain", "Test body")

	// This will fail with connection error, but we're testing message handling
	err := mailer.Raw(msg)

	if err == nil {
		t.Fatal("Expected connection error")
	}

	// Should get a connection error, not a nil message error
	if err.Error() == "message cannot be nil" {
		t.Error("Message was treated as nil when it shouldn't be")
	}

	// The error should be about failed dial
	if !contains(err.Error(), "failed to dial SMTP server") {
		t.Errorf("Expected dial error, got %s", err.Error())
	}
}

// MockDialer is a mock implementation for testing SMTP dialer
type MockDialer struct {
	shouldFailDial bool
	shouldFailSend bool
}

func (m *MockDialer) Dial() (gomail.SendCloser, error) {
	if m.shouldFailDial {
		return nil, errors.New("dial failed")
	}
	return &MockSendCloser{shouldFailSend: m.shouldFailSend}, nil
}

type MockSendCloser struct {
	shouldFailSend bool
	closed         bool
}

func (m *MockSendCloser) Send(from string, to []string, msg io.WriterTo) error {
	if m.shouldFailSend {
		return errors.New("send failed")
	}
	return nil
}

func (m *MockSendCloser) Close() error {
	m.closed = true
	return nil
}

// Helper functions moved to test_helpers.go

func TestMailerInterface(t *testing.T) {
	config := &Config{SMTP: &SMTPConfig{}}
	mailer := NewSMTPMailer(config)

	// Test that smtpMailer implements Mailer interface
	var _ Mailer = mailer

	// Test interface methods
	message := mailer.NewMessage()
	if message == nil {
		t.Error("NewMessage() should return a valid message")
	}

	// Test Send method exists and handles nil properly
	err := mailer.Send(nil)
	if err == nil {
		t.Error("Send(nil) should return error")
	}

	// Test Raw method exists and handles nil properly
	err = mailer.Raw(nil)
	if err == nil {
		t.Error("Raw(nil) should return error")
	}
}
