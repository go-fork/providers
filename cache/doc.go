// Package cache cung cấp một hệ thống quản lý cache hiện đại, linh hoạt và có khả năng mở rộng cao
// cho framework dependency injection go-fork.
//
// Package này được thiết kế theo kiến trúc đa driver, cho phép tích hợp dễ dàng với nhiều loại
// storage backend khác nhau như Memory, File, Redis và MongoDB. Với interface Driver đã được thiết kế
// cẩn thận, package hỗ trợ việc tạo custom driver một cách dễ dàng.
//
// # Tính năng chính
//
//   - Đa dạng driver: Memory, File, Redis, MongoDB với khả năng tùy chỉnh
//   - TTL (Time To Live): Quản lý thời gian sống tự động cho cache entries
//   - Remember Pattern: Lazy computation với caching kết quả tự động
//   - Batch Operations: GetMultiple, SetMultiple, DeleteMultiple để tối ưu hiệu suất
//   - Thread-Safe: An toàn cho môi trường đa luồng với sync.RWMutex
//   - Dependency Injection: ServiceProvider tích hợp với DI container
//   - Monitoring: Stats() method cho metrics và performance tracking
//   - High Performance: Memory driver với automatic cleanup của expired entries
//
// # Cấu trúc Package
//
//	cache/
//	├── manager.go              # Manager interface và DefaultManager implementation
//	├── provider.go             # ServiceProvider cho DI integration
//	├── doc.go                  # Package documentation
//	├── config/
//	│   ├── config.go           # Configuration structs và loading
//	│   └── config_test.go      # Configuration tests
//	├── driver/
//	│   ├── driver.go           # Driver interface definition
//	│   ├── memory.go           # In-memory cache driver
//	│   ├── file.go             # File-based cache driver
//	│   ├── redis.go            # Redis cache driver (v9+)
//	│   └── mongodb.go          # MongoDB cache driver
//	├── mocks/                  # Auto-generated mocks cho testing
//	└── configs/                # Sample configuration files
//
// # Sử dụng với Service Provider (Khuyến nghị)
//
//	package main
//
//	import (
//	    "time"
//	    "go.fork.vn/di"
//	    "go.fork.vn/providers/cache"
//	)
//
//	func main() {
//	    // Khởi tạo DI container
//	    container := di.New()
//
//	    // Đăng ký Cache Service Provider
//	    provider := cache.NewServiceProvider()
//	    provider.Register(container)
//	    provider.Boot(container)
//
//	    // Resolve cache manager từ container
//	    cacheManager := container.MustMake("cache").(cache.Manager)
//
//	    // Sử dụng cache manager
//	    err := cacheManager.Set("user:1", userData, 1*time.Hour)
//	    userData, exists := cacheManager.Get("user:1")
//
//	    defer cacheManager.Close()
//	}
//
// # Sử dụng trực tiếp
//
//	package main
//
//	import (
//	    "time"
//	    "go.fork.vn/providers/cache"
//	    "go.fork.vn/providers/cache/driver"
//	    "go.fork.vn/providers/cache/config"
//	)
//
//	func main() {
//	    // Tạo cache manager
//	    manager := cache.NewManager()
//
//	    // Thêm memory driver với cấu hình
//	    memConfig := config.DriverMemoryConfig{
//	        DefaultTTL:      3600, // 1 hour in seconds
//	        CleanupInterval: 600,  // 10 minutes in seconds
//	        MaxItems:        1000,
//	    }
//	    memDriver := driver.NewMemoryDriver(memConfig)
//	    manager.AddDriver("memory", memDriver)
//	    manager.SetDefaultDriver("memory")
//
//	    // Sử dụng cache
//	    manager.Set("key", "value", 30*time.Minute)
//	    value, exists := manager.Get("key")
//
//	    defer manager.Close()
//	}
//
// # Remember Pattern
//
// Remember pattern giúp tránh tính toán lặp lại cho những operation tốn kém:
//
//	userData, err := manager.Remember("user:123", 1*time.Hour, func() (interface{}, error) {
//	    // Callback chỉ được gọi khi cache miss
//	    return fetchUserFromDatabase(123)
//	})
//
// # Batch Operations
//
// Sử dụng batch operations để giảm số lần gọi API và tối ưu hiệu suất:
//
//	// Batch get
//	userIDs := []string{"user:1", "user:2", "user:3"}
//	users, missingKeys := manager.GetMultiple(userIDs)
//
//	// Batch set
//	userDataMap := map[string]interface{}{
//	    "user:4": user4Data,
//	    "user:5": user5Data,
//	}
//	err := manager.SetMultiple(userDataMap, 1*time.Hour)
//
//	// Batch delete
//	err = manager.DeleteMultiple([]string{"user:1", "user:2"})
//
// # Custom Driver Development
//
// Tạo driver tùy chỉnh bằng cách implement interface Driver:
//
//	type CustomDriver struct {
//	    // ... fields
//	}
//
//	func (d *CustomDriver) Get(ctx context.Context, key string) (interface{}, bool) {
//	    // Implementation
//	}
//
//	func (d *CustomDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
//	    // Implementation
//	}
//
//	// ... implement other methods
//
// # Driver Types
//
// Memory Driver: Lưu cache trong RAM, tốc độ cao nhất, phù hợp cho single instance.
//
//	memConfig := config.DriverMemoryConfig{
//	    DefaultTTL:      3600,
//	    CleanupInterval: 600,
//	    MaxItems:        1000,
//	}
//	memDriver := driver.NewMemoryDriver(memConfig)
//
// File Driver: Lưu cache trong file system, persistence data, ít network overhead.
//
//	fileConfig := config.DriverFileConfig{
//	    Path:            "/tmp/cache",
//	    DefaultTTL:      1800,
//	    FilePermissions: "0644",
//	}
//	fileDriver, err := driver.NewFileDriver(fileConfig)
//
// Redis Driver: Phù hợp cho distributed systems, high performance, persistence.
//
//	import "github.com/redis/go-redis/v9"
//
//	redisClient := redis.NewClient(&redis.Options{
//	    Addr: "localhost:6379",
//	})
//	redisDriver := driver.NewRedisDriver(redisClient)
//
// MongoDB Driver: Phù hợp khi cần query phức tạp và document-based storage.
//
//	import "go.mongodb.org/mongo-driver/mongo"
//
//	mongoDriver, err := driver.NewMongoDBDriver(
//	    mongoClient,
//	    "cache_database",
//	    "cache_collection",
//	)
//
// # Configuration
//
// Package hỗ trợ cấu hình thông qua YAML/JSON:
//
//	cache:
//	  default_driver: "memory"
//	  drivers:
//	    memory:
//	      type: "memory"
//	      default_ttl: 3600
//	      cleanup_interval: 600
//	      max_items: 1000
//	    redis:
//	      type: "redis"
//	      addr: "localhost:6379"
//	      password: ""
//	      db: 0
//	      prefix: "myapp:"
//
// # Testing
//
// Package cung cấp mocks cho testing:
//
//	import "go.fork.vn/providers/cache/mocks"
//
//	func TestCacheFunction(t *testing.T) {
//	    mockManager := mocks.NewMockManager(t)
//	    mockManager.On("Get", "key").Return("value", true)
//	    mockManager.On("Set", "key", mock.Anything, mock.Anything).Return(nil)
//
//	    // Test your function
//	    result := YourFunction(mockManager)
//
//	    mockManager.AssertExpectations(t)
//	}
//
// # Performance Tips
//
//   - Chọn driver phù hợp với use case
//   - Sử dụng batch operations khi có thể
//   - Set TTL hợp lý để tránh memory leak
//   - Monitor cache hit ratio
//   - Sử dụng connection pooling cho Redis/MongoDB
//
// Package cache giúp tối ưu hóa hiệu suất ứng dụng bằng cách giảm thiểu tải trọng
// truy vấn database và cải thiện thời gian phản hồi của hệ thống.
package cache
