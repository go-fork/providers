package mailer

import (
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	// Test with nil config
	manager, err := NewManager(nil)
	if err != nil {
		t.Fatalf("NewManager(nil) failed: %v", err)
	}

	if manager == nil {
		t.Fatal("NewManager(nil) returned nil manager")
	}

	// Should use default config
	config := manager.Config()
	if config == nil {
		t.Fatal("Manager config is nil")
	}

	if config.SMTP == nil {
		t.Fatal("Manager SMTP config is nil")
	}

	// Test with valid config
	customConfig := &Config{
		SMTP: &SMTPConfig{
			Host:        "smtp.example.com",
			Port:        587,
			Username:    "user@example.com",
			Password:    "password",
			Encryption:  "tls",
			FromAddress: "test@example.com",
			FromName:    "Test Sender",
			Timeout:     30 * time.Second,
		},
		Queue: &QueueConfig{
			Enabled:      false,
			Name:         "test_mailer",
			Adapter:      "memory",
			MaxRetries:   5,
			RetryDelay:   120,
			DelayTimeout: 60,
			FailFast:     false,
			TrackStatus:  true,
		},
	}

	manager, err = NewManager(customConfig)
	if err != nil {
		t.Fatalf("NewManager() with config failed: %v", err)
	}

	if manager == nil {
		t.Fatal("NewManager() returned nil manager")
	}

	if manager.Config() != customConfig {
		t.Error("Manager config should match input config")
	}
}

func TestNewManagerWithConfig(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
		Queue: &QueueConfig{
			Enabled: false,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	if manager == nil {
		t.Fatal("NewManagerWithConfig() returned nil manager")
	}

	// Test mailer initialization
	mailer := manager.Mailer()
	if mailer == nil {
		t.Fatal("Manager mailer is nil")
	}

	// Test queue settings
	queueSettings := manager.QueueSettings()
	if queueSettings == nil {
		t.Fatal("Manager queue settings is nil")
	}

	if queueSettings.QueueEnabled != false {
		t.Error("Queue should be disabled")
	}
}

func TestNewManagerWithConfig_QueueEnabled(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
		Queue: &QueueConfig{
			Enabled:      true,
			Name:         "test_queue",
			Adapter:      "memory",
			MaxRetries:   3,
			RetryDelay:   60,
			DelayTimeout: 30,
			FailFast:     false,
			TrackStatus:  true,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() with queue failed: %v", err)
	}

	if !manager.QueueEnabled() {
		t.Error("Queue should be enabled")
	}

	queueManager := manager.QueueManager()
	if queueManager == nil {
		t.Fatal("Queue manager should not be nil when queue is enabled")
	}
}

func TestNewManagerWithConfig_RedisQueue(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
		Queue: &QueueConfig{
			Enabled:      true,
			Name:         "redis_queue",
			Adapter:      "redis",
			MaxRetries:   5,
			RetryDelay:   120,
			DelayTimeout: 60,
			FailFast:     true,
			TrackStatus:  false,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() with Redis queue failed: %v", err)
	}

	if !manager.QueueEnabled() {
		t.Error("Queue should be enabled")
	}

	queueSettings := manager.QueueSettings()
	if queueSettings.QueueAdapter != "redis" {
		t.Errorf("Expected Redis adapter, got %s", queueSettings.QueueAdapter)
	}
}

func TestManager_Mailer(t *testing.T) {
	manager, err := NewManager(nil)
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	mailer := manager.Mailer()
	if mailer == nil {
		t.Fatal("Mailer() returned nil")
	}

	// Test that mailer implements Mailer interface
	var _ Mailer = mailer

	// Test mailer functionality
	message := mailer.NewMessage()
	if message == nil {
		t.Error("Mailer.NewMessage() returned nil")
	}
}

func TestManager_Config(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host: "test.example.com",
			Port: 587,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	returnedConfig := manager.Config()
	if returnedConfig != config {
		t.Error("Config() should return the same config instance")
	}
}

func TestManager_QueueSettings(t *testing.T) {
	config := &Config{
		Queue: &QueueConfig{
			Enabled:      true,
			Name:         "custom_queue",
			Adapter:      "memory",
			MaxRetries:   10,
			RetryDelay:   300,
			DelayTimeout: 120,
			FailFast:     true,
			TrackStatus:  false,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	queueSettings := manager.QueueSettings()
	if queueSettings == nil {
		t.Fatal("QueueSettings() returned nil")
	}

	if queueSettings.QueueEnabled != true {
		t.Error("QueueEnabled should be true")
	}

	if queueSettings.QueueName != "custom_queue" {
		t.Errorf("Expected QueueName to be 'custom_queue', got %s", queueSettings.QueueName)
	}

	if queueSettings.MaxRetries != 10 {
		t.Errorf("Expected MaxRetries to be 10, got %d", queueSettings.MaxRetries)
	}

	if queueSettings.RetryDelay != 300*time.Second {
		t.Errorf("Expected RetryDelay to be 300s, got %v", queueSettings.RetryDelay)
	}
}

func TestManager_QueueEnabled(t *testing.T) {
	// Test with queue disabled
	config := &Config{
		Queue: &QueueConfig{
			Enabled: false,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	if manager.QueueEnabled() {
		t.Error("QueueEnabled() should return false when queue is disabled")
	}

	// Test with queue enabled
	config.Queue.Enabled = true

	manager, err = NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() with enabled queue failed: %v", err)
	}

	if !manager.QueueEnabled() {
		t.Error("QueueEnabled() should return true when queue is enabled")
	}
}

func TestManager_QueueManager(t *testing.T) {
	// Test with queue disabled
	config := &Config{
		Queue: &QueueConfig{
			Enabled: false,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	queueManager := manager.QueueManager()
	if queueManager != nil {
		t.Error("QueueManager() should return nil when queue is disabled")
	}

	// Test with queue enabled
	config.Queue.Enabled = true

	manager, err = NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() with enabled queue failed: %v", err)
	}

	queueManager = manager.QueueManager()
	if queueManager == nil {
		t.Error("QueueManager() should return queue manager when queue is enabled")
	}
}

func TestManager_QueueClient(t *testing.T) {
	// Test with queue disabled
	config := &Config{
		Queue: &QueueConfig{
			Enabled: false,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	queueClient := manager.QueueClient()
	if queueClient != nil {
		t.Error("QueueClient() should return nil when queue is disabled")
	}

	// Test with queue enabled
	config.Queue.Enabled = true

	manager, err = NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() with enabled queue failed: %v", err)
	}

	queueClient = manager.QueueClient()
	if queueClient == nil {
		t.Error("QueueClient() should return queue client when queue is enabled")
	}
}

func TestManager_EnqueueMessage_QueueDisabled(t *testing.T) {
	// Create manager with queue disabled
	config := &Config{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
		Queue: &QueueConfig{
			Enabled: false,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	// Create a message
	message := NewMessage().
		To("recipient@example.com").
		Subject("Test").
		Text("Test body")

	// EnqueueMessage should send directly when queue is disabled
	err = manager.EnqueueMessage(message)

	// We expect a connection error since we're not running a real SMTP server
	if err == nil {
		t.Fatal("Expected connection error when sending directly")
	}

	// Should not be a queue-related error
	if contains(err.Error(), "queue") {
		t.Errorf("Error should not be queue-related when queue is disabled, got: %s", err.Error())
	}
}

func TestManager_EnqueueMessage_QueueEnabled(t *testing.T) {
	// Create manager with queue enabled
	config := &Config{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
		Queue: &QueueConfig{
			Enabled:      true,
			Name:         "test_queue",
			Adapter:      "memory",
			MaxRetries:   3,
			RetryDelay:   60,
			DelayTimeout: 30,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	// Create a valid message
	message := NewMessage().
		To("recipient@example.com").
		Subject("Queue Test").
		Text("Test body for queue")

	// This should attempt to enqueue the message
	err = manager.EnqueueMessage(message)

	// We might get an error here depending on the queue implementation
	// but it shouldn't be a direct SMTP connection error
	if err != nil {
		// The error should be queue-related, not SMTP-related
		if contains(err.Error(), "dial") {
			t.Errorf("Should not get SMTP dial error when using queue, got: %s", err.Error())
		}
	}
}

func TestManager_ProcessMessage(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	// Create a message and marshal it to JSON
	message := NewMessage().
		To("recipient@example.com").
		Subject("Process Test").
		Text("Test body for processing")

	messageData, err := message.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() failed: %v", err)
	}

	// Process the message
	err = manager.ProcessMessage(messageData)

	// We expect a connection error since we're not running a real SMTP server
	if err == nil {
		t.Fatal("Expected connection error when processing message")
	}

	// Should be an SMTP connection error, not a JSON parsing error
	if contains(err.Error(), "json") || contains(err.Error(), "unmarshal") {
		t.Errorf("Should not get JSON error when processing valid message data, got: %s", err.Error())
	}
}

func TestManager_ProcessMessage_InvalidData(t *testing.T) {
	config := &Config{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	// Process invalid message data
	err = manager.ProcessMessage([]byte("invalid json"))

	if err == nil {
		t.Fatal("Expected error when processing invalid message data")
	}

	// The error should be related to JSON parsing
	if err.Error() == "invalid character 'i' looking for beginning of value" {
		// This is the expected JSON parsing error
	} else if !contains(err.Error(), "json") && !contains(err.Error(), "unmarshal") && !contains(err.Error(), "invalid") {
		t.Errorf("Expected JSON parsing error, got: %s", err.Error())
	}
}

func TestManagerInterface(t *testing.T) {
	manager, err := NewManager(nil)
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	// Test that manager implements Manager interface
	var _ Manager = manager

	// Test all interface methods
	if manager.Mailer() == nil {
		t.Error("Mailer() should return a valid mailer")
	}

	if manager.Config() == nil {
		t.Error("Config() should return a valid config")
	}

	if manager.QueueSettings() == nil {
		t.Error("QueueSettings() should return valid queue settings")
	}

	// QueueEnabled should return a boolean
	_ = manager.QueueEnabled()

	// QueueManager and QueueClient might return nil depending on configuration
	_ = manager.QueueManager()
	_ = manager.QueueClient()

	// Test EnqueueMessage with nil message
	err = manager.EnqueueMessage(nil)
	if err == nil {
		t.Error("EnqueueMessage(nil) should return an error")
	}

	// Test ProcessMessage with empty data
	err = manager.ProcessMessage([]byte{})
	if err == nil {
		t.Error("ProcessMessage with empty data should return an error")
	}
}

func TestManager_WithQueueSettings(t *testing.T) {
	// Create a basic manager with queue settings
	config := &Config{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			FromAddress: "test@example.com",
			FromName:    "Test",
		},
		Queue: &QueueConfig{
			Enabled:      true,
			Name:         "test_queue",
			Adapter:      "memory",
			MaxRetries:   3,
			RetryDelay:   60,
			DelayTimeout: 30,
		},
	}

	manager, err := NewManagerWithConfig(config)
	if err != nil {
		t.Fatalf("NewManagerWithConfig() failed: %v", err)
	}

	// Test that queue is properly configured
	if !manager.QueueEnabled() {
		t.Error("Queue should be enabled")
	}

	queueSettings := manager.QueueSettings()
	if queueSettings.QueueName != "test_queue" {
		t.Errorf("Expected queue name to be 'test_queue', got %s", queueSettings.QueueName)
	}

	if queueSettings.QueueAdapter != "memory" {
		t.Errorf("Expected queue adapter to be 'memory', got %s", queueSettings.QueueAdapter)
	}
}
