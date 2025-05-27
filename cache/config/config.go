package config

import "time"

// Config là cấu trúc cấu hình chính cho cache provider.
//
// Config định nghĩa các tùy chọn cấu hình cho cache manager và các driver.
// Nó hỗ trợ nhiều driver khác nhau như memory, file, redis và mongodb.
type Config struct {
	// DefaultDriver chỉ định driver mặc định để sử dụng
	// Options: memory, file, redis, mongodb
	DefaultDriver string `mapstructure:"default_driver" yaml:"default_driver"`

	// DefaultTTL là thời gian sống mặc định cho cache entries (giây)
	// 0 = không hết hạn, -1 = sử dụng mặc định của driver
	DefaultTTL int `mapstructure:"default_ttl" yaml:"default_ttl"`

	// Prefix là tiền tố cache key để tránh xung đột với ứng dụng khác
	Prefix string `mapstructure:"prefix" yaml:"prefix"`

	// Drivers chứa cấu hình cho từng driver
	Drivers DriversConfig `mapstructure:"drivers" yaml:"drivers"`
}

// DriversConfig chứa cấu hình cho tất cả các driver.
type DriversConfig struct {
	// Memory driver configuration
	Memory *DriverMemoryConfig `mapstructure:"memory" yaml:"memory"`

	// File driver configuration
	File *DriverFileConfig `mapstructure:"file" yaml:"file"`

	// Redis driver configuration
	Redis *DriverRedisConfig `mapstructure:"redis" yaml:"redis"`

	// MongoDB driver configuration
	MongoDB *DriverMongodbConfig `mapstructure:"mongodb" yaml:"mongodb"`
}

// DriverMemoryConfig là cấu hình cho memory driver.
type DriverMemoryConfig struct {
	// Enabled xác định có kích hoạt Memory driver không
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
	// DefaultTTL là thời gian hết hạn mặc định cho memory cache (giây)
	DefaultTTL int `mapstructure:"default_ttl" yaml:"default_ttl"`

	// CleanupInterval là khoảng thời gian dọn dẹp các item hết hạn (giây)
	CleanupInterval int `mapstructure:"cleanup_interval" yaml:"cleanup_interval"`

	// MaxItems là số lượng item tối đa trong memory cache (0 = unlimited)
	MaxItems int `mapstructure:"max_items" yaml:"max_items"`
}

// DriverFileConfig là cấu hình cho file driver.
type DriverFileConfig struct {
	// Enabled xác định có kích hoạt File driver không
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
	// Path là đường dẫn thư mục lưu trữ cache files
	Path string `mapstructure:"path" yaml:"path"`

	// DefaultTTL là thời gian hết hạn mặc định cho file cache (giây)
	DefaultTTL int `mapstructure:"default_ttl" yaml:"default_ttl"`

	// Extension là phần mở rộng cho cache files
	Extension string `mapstructure:"extension" yaml:"extension"`

	// CleanupInterval là khoảng thời gian dọn dẹp các file hết hạn (giây)
	CleanupInterval int `mapstructure:"cleanup_interval" yaml:"cleanup_interval"`
}

// DriverRedisConfig là cấu hình cho redis driver.
type DriverRedisConfig struct {
	// Enabled xác định có kích hoạt Redis driver không
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`

	// DefaultTTL là thời gian hết hạn mặc định cho Redis cache (giây)
	DefaultTTL int `mapstructure:"default_ttl" yaml:"default_ttl"`

	// Serializer là định dạng serialization: json, gob, msgpack
	Serializer string `mapstructure:"serializer" yaml:"serializer"`
}

// DriverMongodbConfig là cấu hình cho mongodb driver.
type DriverMongodbConfig struct {
	// Enabled xác định có kích hoạt MongoDB driver không
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`

	// Database là tên database để lưu trữ cache
	Database string `mapstructure:"database" yaml:"database"`

	// Collection là tên collection để lưu trữ cache
	Collection string `mapstructure:"collection" yaml:"collection"`

	// DefaultTTL là thời gian hết hạn mặc định cho MongoDB cache (giây)
	DefaultTTL int `mapstructure:"default_ttl" yaml:"default_ttl"`

	// Hits là số lần cache hit (readonly, được quản lý bởi driver)
	Hits int64 `mapstructure:"hits" yaml:"hits"`

	// Misses là số lần cache miss (readonly, được quản lý bởi driver)
	Misses int64 `mapstructure:"misses" yaml:"misses"`
}

// DefaultConfig trả về cấu hình mặc định cho cache.
//
// Cấu hình mặc định sử dụng memory driver với TTL 1 giờ.
//
// Returns:
//   - *Config: Cấu hình mặc định
func DefaultConfig() *Config {
	return &Config{
		DefaultDriver: "memory",
		DefaultTTL:    3600, // 1 hour
		Prefix:        "cache:",
		Drivers: DriversConfig{
			Memory: &DriverMemoryConfig{
				DefaultTTL:      3600, // 1 hour
				CleanupInterval: 600,  // 10 minutes
				MaxItems:        10000,
			},
			File: &DriverFileConfig{
				Path:            "./storage/cache",
				DefaultTTL:      3600, // 1 hour
				Extension:       ".cache",
				CleanupInterval: 600, // 10 minutes
			},
			Redis: &DriverRedisConfig{
				Enabled:    true,
				DefaultTTL: 3600, // 1 hour
				Serializer: "json",
			},
			MongoDB: &DriverMongodbConfig{
				Enabled:    true,
				Database:   "cache_db",
				Collection: "cache_items",
				DefaultTTL: 3600, // 1 hour
				Hits:       0,
				Misses:     0,
			},
		},
	}
}

// GetDefaultExpiration trả về thời gian hết hạn mặc định theo kiểu time.Duration.
//
// Returns:
//   - time.Duration: Thời gian hết hạn mặc định
func (c *Config) GetDefaultExpiration() time.Duration {
	return time.Duration(c.DefaultTTL) * time.Second
}

// GetMemoryDefaultExpiration trả về thời gian hết hạn mặc định cho memory driver.
//
// Returns:
//   - time.Duration: Thời gian hết hạn mặc định cho memory driver
func (m *DriverMemoryConfig) GetDefaultExpiration() time.Duration {
	return time.Duration(m.DefaultTTL) * time.Second
}

// GetCleanupInterval trả về khoảng thời gian dọn dẹp cho memory driver.
//
// Returns:
//   - time.Duration: Khoảng thời gian dọn dẹp
func (m *DriverMemoryConfig) GetCleanupInterval() time.Duration {
	return time.Duration(m.CleanupInterval) * time.Second
}

// GetFileDefaultExpiration trả về thời gian hết hạn mặc định cho file driver.
//
// Returns:
//   - time.Duration: Thời gian hết hạn mặc định cho file driver
func (f *DriverFileConfig) GetDefaultExpiration() time.Duration {
	return time.Duration(f.DefaultTTL) * time.Second
}

// GetFileCleanupInterval trả về khoảng thời gian dọn dẹp cho file driver.
//
// Returns:
//   - time.Duration: Khoảng thời gian dọn dẹp
func (f *DriverFileConfig) GetFileCleanupInterval() time.Duration {
	return time.Duration(f.CleanupInterval) * time.Second
}

// GetRedisDefaultExpiration trả về thời gian hết hạn mặc định cho redis driver.
//
// Returns:
//   - time.Duration: Thời gian hết hạn mặc định cho redis driver
func (r *DriverRedisConfig) GetDefaultExpiration() time.Duration {
	return time.Duration(r.DefaultTTL) * time.Second
}

// GetMongoDBDefaultExpiration trả về thời gian hết hạn mặc định cho mongodb driver.
//
// Returns:
//   - time.Duration: Thời gian hết hạn mặc định cho mongodb driver
func (m *DriverMongodbConfig) GetDefaultExpiration() time.Duration {
	return time.Duration(m.DefaultTTL) * time.Second
}
