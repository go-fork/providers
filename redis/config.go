package redis

import (
	"fmt"
	"time"
)

// Config là cấu trúc dữ liệu cấu hình chính cho Redis.
type Config struct {
	// Client là cấu hình cho Redis Standard Client.
	Client *ClientConfig `mapstructure:"client"`

	// Universal là cấu hình cho Redis Universal Client (hỗ trợ Cluster, Sentinel, và standalone).
	Universal *UniversalConfig `mapstructure:"universal"`
}

// ClientConfig chứa cấu hình cho Redis Standard Client.
type ClientConfig struct {
	// Host là địa chỉ máy chủ Redis.
	Host string `mapstructure:"host"`

	// Port là cổng kết nối của máy chủ Redis.
	Port int `mapstructure:"port"`

	// Password là mật khẩu xác thực với máy chủ Redis.
	Password string `mapstructure:"password"`

	// DB là số của database Redis sẽ được sử dụng.
	DB int `mapstructure:"db"`

	// Prefix là tiền tố sẽ được thêm vào tất cả các khóa.
	Prefix string `mapstructure:"prefix"`

	// Timeout là thời gian chờ chung cho các hoạt động Redis.
	Timeout int `mapstructure:"timeout"` // seconds

	// DialTimeout là thời gian chờ khi thiết lập kết nối tới Redis.
	DialTimeout int `mapstructure:"dial_timeout"` // seconds

	// ReadTimeout là thời gian chờ khi đọc dữ liệu từ Redis.
	ReadTimeout int `mapstructure:"read_timeout"` // seconds

	// WriteTimeout là thời gian chờ khi ghi dữ liệu vào Redis.
	WriteTimeout int `mapstructure:"write_timeout"` // seconds

	// PoolSize là số lượng kết nối tối đa được giữ trong pool.
	PoolSize int `mapstructure:"pool_size"`

	// MinIdleConns là số lượng kết nối rảnh tối thiểu được giữ trong pool.
	MinIdleConns int `mapstructure:"min_idle_conns"`
}

// UniversalConfig chứa cấu hình cho Redis Universal Client.
type UniversalConfig struct {
	// Addresses là danh sách các địa chỉ máy chủ Redis (host:port)
	Addresses []string `mapstructure:"addresses"`

	// Password là mật khẩu xác thực với máy chủ Redis.
	Password string `mapstructure:"password"`

	// DB là số của database Redis sẽ được sử dụng.
	DB int `mapstructure:"db"`

	// Prefix là tiền tố sẽ được thêm vào tất cả các khóa.
	Prefix string `mapstructure:"prefix"`

	// Timeout là thời gian chờ chung cho các hoạt động Redis.
	Timeout int `mapstructure:"timeout"` // seconds

	// DialTimeout là thời gian chờ khi thiết lập kết nối tới Redis.
	DialTimeout int `mapstructure:"dial_timeout"` // seconds

	// ReadTimeout là thời gian chờ khi đọc dữ liệu từ Redis.
	ReadTimeout int `mapstructure:"read_timeout"` // seconds

	// WriteTimeout là thời gian chờ khi ghi dữ liệu vào Redis.
	WriteTimeout int `mapstructure:"write_timeout"` // seconds

	// MaxRetries là số lần thử lại tối đa cho các thao tác lỗi.
	MaxRetries int `mapstructure:"max_retries"`

	// MinRetryBackoff là thời gian chờ tối thiểu giữa các lần thử lại (milliseconds).
	MinRetryBackoff int `mapstructure:"min_retry_backoff"` // milliseconds

	// MaxRetryBackoff là thời gian chờ tối đa giữa các lần thử lại (milliseconds).
	MaxRetryBackoff int `mapstructure:"max_retry_backoff"` // milliseconds

	// PoolSize là số lượng kết nối tối đa được giữ trong pool.
	PoolSize int `mapstructure:"pool_size"`

	// MinIdleConns là số lượng kết nối rảnh tối thiểu được giữ trong pool.
	MinIdleConns int `mapstructure:"min_idle_conns"`

	// ClusterMode xác định liệu Universal Client có chạy ở chế độ Cluster hay không.
	ClusterMode bool `mapstructure:"cluster_mode"`

	// MaxRedirects là số lần chuyển hướng tối đa cho các thao tác cluster.
	MaxRedirects int `mapstructure:"max_redirects"`

	// SentinelMode xác định liệu Universal Client có chạy ở chế độ Sentinel hay không.
	SentinelMode bool `mapstructure:"sentinel_mode"`

	// MasterName là tên của master đang được giám sát bởi các máy chủ sentinel.
	MasterName string `mapstructure:"master_name"`
}

// DefaultConfig trả về cấu hình mặc định cho Redis.
func DefaultConfig() *Config {
	return &Config{
		Client: &ClientConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			Prefix:       "",
			Timeout:      5,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			PoolSize:     10,
			MinIdleConns: 5,
		},
		Universal: &UniversalConfig{
			Addresses:       []string{"localhost:6379"},
			Password:        "",
			DB:              0,
			Prefix:          "",
			Timeout:         5,
			DialTimeout:     5,
			ReadTimeout:     3,
			WriteTimeout:    3,
			MaxRetries:      3,
			MinRetryBackoff: 8,
			MaxRetryBackoff: 512,
			PoolSize:        10,
			MinIdleConns:    5,
			ClusterMode:     false,
			MaxRedirects:    3,
			SentinelMode:    false,
			MasterName:      "mymaster",
		},
	}
}

// GetClientOptions trả về Redis Client Options từ cấu hình.
func (c *ClientConfig) GetClientOptions() *ClientOptions {
	return &ClientOptions{
		Addr:         c.Host + ":" + fmt.Sprintf("%d", c.Port),
		Password:     c.Password,
		DB:           c.DB,
		DialTimeout:  time.Duration(c.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(c.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.WriteTimeout) * time.Second,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
	}
}

// GetUniversalOptions trả về Redis Universal Client Options từ cấu hình.
func (c *UniversalConfig) GetUniversalOptions() *UniversalOptions {
	return &UniversalOptions{
		Addrs:           c.Addresses,
		Password:        c.Password,
		DB:              c.DB,
		DialTimeout:     time.Duration(c.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(c.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(c.WriteTimeout) * time.Second,
		MaxRetries:      c.MaxRetries,
		MinRetryBackoff: time.Duration(c.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff: time.Duration(c.MaxRetryBackoff) * time.Millisecond,
		PoolSize:        c.PoolSize,
		MinIdleConns:    c.MinIdleConns,
		RouteRandomly:   true,
		MasterName:      c.MasterName,
	}
}

// ClientOptions là wrapper cho Redis Client Options.
type ClientOptions struct {
	Addr         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	MinIdleConns int
}

// UniversalOptions là wrapper cho Redis Universal Client Options.
type UniversalOptions struct {
	Addrs           []string
	MasterName      string
	Password        string
	DB              int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration
	PoolSize        int
	MinIdleConns    int
	RouteRandomly   bool
}
