package queue

import (
	"github.com/redis/go-redis/v9"
)

// Config chứa cấu hình cho queue package.
type Config struct {
	// Adapter chứa cấu hình cho các queue adapter.
	Adapter AdapterConfig `mapstructure:"adapter"`

	// Server chứa cấu hình cho queue server.
	Server ServerConfig `mapstructure:"server"`

	// Client chứa cấu hình cho queue client.
	Client ClientConfig `mapstructure:"client"`
}

// AdapterConfig chứa cấu hình cho các adapter.
type AdapterConfig struct {
	// Default xác định adapter mặc định sẽ được sử dụng.
	// Các giá trị hợp lệ: "memory", "redis"
	Default string `mapstructure:"default"`

	// Memory chứa cấu hình cho memory adapter.
	Memory MemoryConfig `mapstructure:"memory"`

	// Redis chứa cấu hình cho redis adapter.
	Redis RedisConfig `mapstructure:"redis"`
}

// MemoryConfig chứa cấu hình cho memory adapter.
type MemoryConfig struct {
	// Prefix là tiền tố cho tên của các queue trong bộ nhớ.
	Prefix string `mapstructure:"prefix"`
}

// RedisConfig chứa cấu hình cho redis adapter.
type RedisConfig struct {
	// Address là địa chỉ của Redis server.
	Address string `mapstructure:"address"`

	// Password là mật khẩu của Redis server.
	Password string `mapstructure:"password"`

	// DB là số của Redis database.
	DB int `mapstructure:"db"`

	// TLS xác định liệu có sử dụng TLS khi kết nối đến Redis server hay không.
	TLS bool `mapstructure:"tls"`

	// Prefix là tiền tố cho tên của các queue trong Redis.
	Prefix string `mapstructure:"prefix"`

	// Cluster chứa cấu hình cho Redis cluster.
	Cluster RedisClusterConfig `mapstructure:"cluster"`

	// Client là Redis client được cấu hình trước đó.
	// Nếu được cung cấp, các tùy chọn khác sẽ bị bỏ qua.
	Client redis.UniversalClient `mapstructure:"-"`
}

// RedisClusterConfig chứa cấu hình cho Redis cluster.
type RedisClusterConfig struct {
	// Enabled xác định liệu có sử dụng Redis cluster hay không.
	Enabled bool `mapstructure:"enabled"`

	// Addresses là danh sách các địa chỉ của Redis cluster.
	Addresses []string `mapstructure:"addresses"`
}

// ServerConfig chứa cấu hình cho queue server.
type ServerConfig struct {
	// Concurrency là số lượng worker xử lý tác vụ cùng một lúc.
	Concurrency int `mapstructure:"concurrency"`

	// PollingInterval là khoảng thời gian giữa các lần kiểm tra tác vụ mới (tính bằng mili giây).
	PollingInterval int `mapstructure:"pollingInterval"`

	// DefaultQueue là tên queue mặc định nếu không có queue nào được chỉ định.
	DefaultQueue string `mapstructure:"defaultQueue"`

	// StrictPriority xác định liệu có ưu tiên nghiêm ngặt các queue ưu tiên cao hay không.
	StrictPriority bool `mapstructure:"strictPriority"`

	// Queues là danh sách các queue cần lắng nghe, theo thứ tự ưu tiên.
	Queues []string `mapstructure:"queues"`

	// ShutdownTimeout là thời gian chờ để các worker hoàn tất tác vụ khi dừng server (tính bằng giây).
	ShutdownTimeout int `mapstructure:"shutdownTimeout"`

	// LogLevel xác định mức độ log.
	LogLevel int `mapstructure:"logLevel"`

	// RetryLimit xác định số lần thử lại tối đa cho tác vụ bị lỗi.
	RetryLimit int `mapstructure:"retryLimit"`
}

// ClientConfig chứa cấu hình cho queue client.
type ClientConfig struct {
	// DefaultOptions chứa các tùy chọn mặc định cho tác vụ.
	DefaultOptions ClientDefaultOptions `mapstructure:"defaultOptions"`
}

// ClientDefaultOptions chứa các tùy chọn mặc định cho tác vụ.
type ClientDefaultOptions struct {
	// Queue là tên queue mặc định cho các tác vụ.
	Queue string `mapstructure:"queue"`

	// MaxRetry là số lần thử lại tối đa cho tác vụ bị lỗi.
	MaxRetry int `mapstructure:"maxRetry"`

	// Timeout là thời gian tối đa để tác vụ hoàn thành (tính bằng phút).
	Timeout int `mapstructure:"timeout"`
}

// DefaultConfig trả về cấu hình mặc định cho queue.
func DefaultConfig() Config {
	return Config{
		Adapter: AdapterConfig{
			Default: "memory",
			Memory: MemoryConfig{
				Prefix: "queue:",
			},
			Redis: RedisConfig{
				Address:  "localhost:6379",
				Password: "",
				DB:       0,
				TLS:      false,
				Prefix:   "queue:",
				Cluster: RedisClusterConfig{
					Enabled:   false,
					Addresses: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
				},
			},
		},
		Server: ServerConfig{
			Concurrency:     10,
			PollingInterval: 1000,
			DefaultQueue:    "default",
			StrictPriority:  true,
			Queues:          []string{"critical", "high", "default", "low"},
			ShutdownTimeout: 30,
			LogLevel:        1,
			RetryLimit:      3,
		},
		Client: ClientConfig{
			DefaultOptions: ClientDefaultOptions{
				Queue:    "default",
				MaxRetry: 3,
				Timeout:  30,
			},
		},
	}
}
