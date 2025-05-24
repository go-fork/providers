package mailer

import (
	"time"

	"errors"

	"github.com/go-fork/providers/config"
)

// Config cấu hình cho dịch vụ gửi mail
type Config struct {
	// SMTP cấu hình cho SMTP server
	SMTP *SMTPConfig `mapstructure:"smtp"`

	// Queue cấu hình cho việc gửi email qua hàng đợi
	Queue *QueueConfig `mapstructure:"queue"`
}

// SMTPConfig cấu hình cho SMTP server
type SMTPConfig struct {
	// Host là địa chỉ SMTP server
	Host string `mapstructure:"host"`

	// Port là cổng SMTP server
	Port int `mapstructure:"port"`

	// Username là tên người dùng để xác thực với SMTP server
	Username string `mapstructure:"username"`

	// Password là mật khẩu để xác thực với SMTP server
	Password string `mapstructure:"password"`

	// Encryption loại mã hóa cho SMTP (tls, ssl, none)
	Encryption string `mapstructure:"encryption"`

	// FromAddress là địa chỉ email mặc định dùng để gửi mail
	FromAddress string `mapstructure:"from_address"`

	// FromName là tên người gửi mặc định
	FromName string `mapstructure:"from_name"`

	// Timeout là thời gian chờ tối đa khi kết nối đến SMTP server
	Timeout time.Duration `mapstructure:"timeout"`
}

// QueueConfig chứa cấu hình cho việc gửi email qua hàng đợi
type QueueConfig struct {
	// Enabled xác định liệu có sử dụng queue cho việc gửi email hay không
	Enabled bool `mapstructure:"enabled"`

	// Name là tên của queue cho việc gửi email
	Name string `mapstructure:"name"`

	// Adapter xác định adapter được sử dụng cho queue ("memory" hoặc "redis")
	Adapter string `mapstructure:"adapter"`

	// DelayTimeout là thời gian timeout cho việc xử lý một email (tính bằng giây)
	DelayTimeout int `mapstructure:"delay_timeout"`

	// FailFast xác định liệu có dừng xử lý khi gặp lỗi hay không
	FailFast bool `mapstructure:"fail_fast"`

	// TrackStatus xác định liệu có ghi lại trạng thái gửi email hay không
	TrackStatus bool `mapstructure:"track_status"`

	// MaxRetries là số lần thử lại tối đa nếu gửi email thất bại
	MaxRetries int `mapstructure:"max_retries"`

	// RetryDelay là thời gian chờ giữa các lần thử lại (tính bằng giây)
	RetryDelay int `mapstructure:"retry_delay"`
}

// Các cấu trúc WorkersConfig, RedisConfig và ProcessingConfig đã được loại bỏ
// vì cấu trúc đã được đơn giản hóa theo cấu hình mới

// QueueSettings là cấu hình cho việc gửi email qua hàng đợi
// được chuyển đổi từ QueueConfig để sử dụng trong mã nguồn
type QueueSettings struct {
	// QueueEnabled xác định liệu có sử dụng queue cho việc gửi email hay không
	QueueEnabled bool

	// QueueName là tên của queue cho việc gửi email
	QueueName string

	// MaxRetries là số lần thử lại tối đa nếu gửi email thất bại
	MaxRetries int

	// RetryDelay là thời gian chờ giữa các lần thử lại
	RetryDelay time.Duration

	// WorkerConcurrency là số lượng worker xử lý email đồng thời
	WorkerConcurrency int

	// PollingInterval là khoảng thời gian giữa các lần kiểm tra email mới
	PollingInterval time.Duration

	// RedisAddress là địa chỉ của Redis server
	RedisAddress string

	// RedisPassword là mật khẩu của Redis server
	RedisPassword string

	// RedisDB là số của Redis database
	RedisDB int

	// RedisUseTLS xác định liệu có sử dụng TLS khi kết nối đến Redis server hay không
	RedisUseTLS bool

	// QueuePrefix là tiền tố cho tên của các queue
	QueuePrefix string

	// QueueAdapter xác định adapter được sử dụng cho queue ("memory" hoặc "redis")
	QueueAdapter string

	// ProcessTimeout là thời gian timeout cho việc xử lý một email
	ProcessTimeout time.Duration

	// FailFast xác định liệu có dừng xử lý khi gặp lỗi hay không
	FailFast bool

	// TrackStatus xác định liệu có ghi lại trạng thái gửi email hay không
	TrackStatus bool
}

// NewConfig tạo cấu hình mặc định cho mailer
func NewConfig() *Config {
	return &Config{
		SMTP: &SMTPConfig{
			Host:        "localhost",
			Port:        25,
			Username:    "",
			Password:    "",
			Encryption:  "none",
			FromAddress: "no-reply@example.com",
			FromName:    "System Notification",
			Timeout:     10 * time.Second,
		},
		Queue: &QueueConfig{
			Enabled:      false,
			Name:         "mailer",
			MaxRetries:   3,
			RetryDelay:   60,
			Adapter:      "memory",
			DelayTimeout: 60,
			FailFast:     false,
			TrackStatus:  true,
		},
	}
}

// LoadConfig tải cấu hình từ config manager
func LoadConfig(configManager config.Manager) (*Config, error) {
	cfg := NewConfig()
	if configManager == nil || !configManager.Has("mailer") {
		return nil, errors.New("mailer configuration not found")
	}
	err := configManager.UnmarshalKey("mailer", &cfg)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi thời gian từ giây sang time.Duration
	if cfg.SMTP != nil {
		if cfg.SMTP.Timeout == 0 {
			cfg.SMTP.Timeout = 10 * time.Second
		} else {
			cfg.SMTP.Timeout = time.Duration(cfg.SMTP.Timeout) * time.Second
		}
	}
	return cfg, nil
}

// DefaultQueueSettings tạo QueueSettings mặc định cho mailer
func defaultQueueSettings() *QueueSettings {
	return &QueueSettings{
		QueueEnabled:      false,
		QueueName:         "mailer",
		MaxRetries:        3,
		RetryDelay:        60 * time.Second,
		WorkerConcurrency: 5,
		PollingInterval:   1000 * time.Millisecond,
		RedisAddress:      "localhost:6379",
		RedisPassword:     "",
		RedisDB:           0,
		RedisUseTLS:       false,
		QueuePrefix:       "mailer:",
		QueueAdapter:      "memory",
		ProcessTimeout:    60 * time.Second,
		FailFast:          false,
		TrackStatus:       true,
	}
}

// QueueSettingsFromConfig chuyển đổi QueueConfig thành QueueSettings
func queueSettingsFromConfig(queueConfig *QueueConfig) *QueueSettings {
	if queueConfig == nil {
		return defaultQueueSettings()
	}

	settings := &QueueSettings{
		QueueEnabled:   queueConfig.Enabled,
		QueueName:      queueConfig.Name,
		QueueAdapter:   queueConfig.Adapter,
		MaxRetries:     queueConfig.MaxRetries,
		RetryDelay:     time.Duration(queueConfig.RetryDelay) * time.Second,
		ProcessTimeout: time.Duration(queueConfig.DelayTimeout) * time.Second,
		FailFast:       queueConfig.FailFast,
		TrackStatus:    queueConfig.TrackStatus,
		// Các giá trị mặc định cho worker và Redis
		WorkerConcurrency: 5,
		PollingInterval:   1000 * time.Millisecond,
		RedisAddress:      "localhost:6379",
		RedisPassword:     "",
		RedisDB:           0,
		RedisUseTLS:       false,
		QueuePrefix:       "mailer:",
	}

	return settings
}
