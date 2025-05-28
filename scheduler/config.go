package scheduler

import "time"

// Config là cấu trúc cấu hình chính cho scheduler provider.
//
// Config định nghĩa các tùy chọn cấu hình cho scheduler manager và distributed locking.
// Nó hỗ trợ Redis distributed locking khi chạy trên nhiều instance.
type Config struct {
	// AutoStart xác định có tự động khởi động scheduler khi ứng dụng boot không
	// ServiceProvider sẽ tự động gọi scheduler.StartAsync() trong Boot() method nếu true
	AutoStart bool `mapstructure:"auto_start" yaml:"auto_start"`

	// DistributedLock chứa cấu hình cho distributed locking
	DistributedLock DistributedLockConfig `mapstructure:"distributed_lock" yaml:"distributed_lock"`

	// Options chứa cấu hình RedisLockerOptions cho distributed locking
	Options RedisLockerOptions `mapstructure:"options" yaml:"options"`
}

// DistributedLockConfig chứa cấu hình cho distributed locking.
type DistributedLockConfig struct {
	// Enabled xác định có bật distributed locking không
	// Chỉ cần thiết khi chạy scheduler trên nhiều instance trong môi trường phân tán
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
}

// RedisLockerOptions chứa các tùy chọn cấu hình cho Redis Locker.
type RedisLockerOptions struct {
	// KeyPrefix là tiền tố được thêm vào trước mỗi khóa trong Redis
	KeyPrefix string `mapstructure:"key_prefix" yaml:"key_prefix"`

	// LockDuration là thời gian một khóa sẽ tồn tại trước khi tự động hết hạn (giây)
	LockDuration int `mapstructure:"lock_duration" yaml:"lock_duration"`

	// MaxRetries là số lần thử tối đa khi gặp lỗi khi tương tác với Redis
	MaxRetries int `mapstructure:"max_retries" yaml:"max_retries"`

	// RetryDelay là thời gian chờ giữa các lần thử (milliseconds)
	RetryDelay int `mapstructure:"retry_delay" yaml:"retry_delay"`
}

// DefaultConfig trả về cấu hình mặc định cho scheduler.
func DefaultConfig() Config {
	return Config{
		AutoStart: true,
		DistributedLock: DistributedLockConfig{
			Enabled: false,
		},
		Options: DefaultRedisLockerOptions(),
	}
}

// DefaultRedisLockerOptions trả về các tùy chọn mặc định cho Redis Locker.
func DefaultRedisLockerOptions() RedisLockerOptions {
	return RedisLockerOptions{
		KeyPrefix:    "scheduler_lock:",
		LockDuration: 30, // 30 seconds
		MaxRetries:   3,
		RetryDelay:   100, // 100 milliseconds
	}
}

// ToTimeDuration chuyển đổi các giá trị int trong config thành time.Duration.
func (opts RedisLockerOptions) ToTimeDuration() RedisLockerOptionsTime {
	return RedisLockerOptionsTime{
		KeyPrefix:    opts.KeyPrefix,
		LockDuration: time.Duration(opts.LockDuration) * time.Second,
		MaxRetries:   opts.MaxRetries,
		RetryDelay:   time.Duration(opts.RetryDelay) * time.Millisecond,
	}
}

// RedisLockerOptionsTime chứa các tùy chọn cấu hình với time.Duration.
type RedisLockerOptionsTime struct {
	// KeyPrefix là tiền tố được thêm vào trước mỗi khóa trong Redis
	KeyPrefix string

	// LockDuration là thời gian một khóa sẽ tồn tại trước khi tự động hết hạn
	LockDuration time.Duration

	// MaxRetries là số lần thử tối đa khi gặp lỗi khi tương tác với Redis
	MaxRetries int

	// RetryDelay là thời gian chờ giữa các lần thử
	RetryDelay time.Duration
}
