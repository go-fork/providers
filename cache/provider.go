// Package cache cung cấp framework cache với nhiều driver cho ứng dụng Go.
//
// Package này triển khai một hệ thống cache đa driver, thread-safe và extensible cho các ứng dụng Go.
// Nó hỗ trợ các driver như Memory, File, Redis, MongoDB, và có thể mở rộng với các driver tùy chỉnh.
// Package cung cấp các phương thức tiêu chuẩn như Get, Set, Delete, Flush, và các chức năng nâng cao như
// Remember, GetMultiple, SetMultiple để tối ưu hiệu suất tương tác với cache.
package cache

import (
	"go.fork.vn/di"
	"go.fork.vn/providers/cache/config"
	"go.fork.vn/providers/cache/driver"
	configService "go.fork.vn/providers/config"
	"go.fork.vn/providers/mongodb"
	"go.fork.vn/providers/redis"
)

// ServiceProvider là interface cho cache service provider.
//
// ServiceProvider định nghĩa các phương thức cần thiết cho một cache service provider
// và kế thừa từ di.ServiceProvider.
type ServiceProvider interface {
	di.ServiceProvider
}

// serviceProvider là service provider cho module cache.
//
// serviceProvider đảm nhận việc đăng ký và khởi tạo các dịch vụ cache vào container DI.
// Nó cung cấp cơ chế để đăng ký cache manager và các driver mặc định vào ứng dụng.
type serviceProvider struct {
	providers []string
}

// NewServiceProvider tạo một cache service provider mới.
//
// Phương thức này khởi tạo một service provider mới cho module cache.
//
// Returns:
//   - ServiceProvider: Đối tượng service provider đã sẵn sàng đăng ký vào ứng dụng
func NewServiceProvider() ServiceProvider {
	return &serviceProvider{}
}

// Requires trả về danh sách các service provider mà provider này phụ thuộc.
//
// Cache provider phụ thuộc vào config provider để đọc cấu hình.
//
// Trả về:
//   - []string: danh sách các service provider khác mà provider này yêu cầu
func (p *serviceProvider) Requires() []string {
	return []string{"config", "redis", "mongodb"}
}

// Register đăng ký các binding vào container.
//
// Phương thức này đăng ký cache manager vào container DI của ứng dụng.
// Nó khởi tạo một cache manager mới và đăng ký nó dưới các alias "cache" và "cache.manager".
// Cấu hình sẽ được load từ config manager trong quá trình Boot.
//
// Params:
//   - app: Interface của ứng dụng, phải cung cấp phương thức Container() để lấy container DI
func (p *serviceProvider) Register(app interface{}) {
	// Lấy container từ app
	if appWithContainer, ok := app.(interface {
		Container() *di.Container
	}); ok {
		c := appWithContainer.Container()

		// Load cache manager và config manager

		configManager := c.MustMake("config").(configService.Manager)

		var cfg config.Config
		if err := configManager.UnmarshalKey("cache", &cfg); err != nil {
			panic("Cache config unmarshal error: " + err.Error())
		}
		manager := NewManager()
		// Đăng ký cache service
		c.Instance("cache", manager)
		p.providers = append(p.providers, "cache")
		var (
			redisManager   redis.Manager
			mongodbManager mongodb.Manager
		)

		if cfg.Drivers.Memory != nil && cfg.Drivers.Memory.Enabled {
			// Đăng ký Memory Driver vào cache manager
			memoryDriver := driver.NewMemoryDriver(*cfg.Drivers.Memory)
			manager.AddDriver("memory", memoryDriver)
			c.Instance("cache.memory", memoryDriver)
			p.providers = append(p.providers, "cache.memory")
		}

		if cfg.Drivers.File != nil && cfg.Drivers.File.Enabled {
			// Đăng ký File Driver vào cache manager
			fileDriver, err := driver.NewFileDriver(*cfg.Drivers.File)
			if err != nil {
				panic("Failed to create File driver: " + err.Error())
			}
			manager.AddDriver("file", fileDriver)
			c.Instance("cache.file", fileDriver)
			p.providers = append(p.providers, "cache.file")
		}

		if cfg.Drivers.Redis != nil && cfg.Drivers.Redis.Enabled {
			redisManager = c.MustMake("redis").(redis.Manager)
			if redisManager == nil {
				panic("Redis manager is nil, please ensure Redis provider is registered")
			}
			// Đăng ký Redis Driver vào cache manager
			redisDriver, err := driver.NewRedisDriver(*cfg.Drivers.Redis, redisManager)
			if err != nil {
				panic("Failed to create Redis driver: " + err.Error())
			}
			manager.AddDriver("redis", redisDriver)
			c.Instance("cache.redis", redisDriver)
			p.providers = append(p.providers, "cache.redis")
		}

		if cfg.Drivers.MongoDB != nil && cfg.Drivers.MongoDB.Enabled {
			mongodbManager = c.MustMake("mongodb").(mongodb.Manager)
			if mongodbManager == nil {
				panic("MongoDB manager is nil, please ensure MongoDB provider is registered")
			}

			// Đăng ký MongoDB Driver vào cache manager
			mongodbDriver, err := driver.NewMongoDBDriver(*cfg.Drivers.MongoDB, mongodbManager)
			if err != nil {
				panic("Failed to create MongoDB driver: " + err.Error())
			}
			manager.AddDriver("mongodb", mongodbDriver)
			c.Instance("cache.mongodb", mongodbDriver)
			p.providers = append(p.providers, "cache.mongodb")
		}
	}
}

// Boot được gọi sau khi tất cả các service provider đã được đăng ký.
//
// Phương thức này khởi tạo và cấu hình các cache driver dựa trên cấu hình.
// Nó sẽ load cấu hình từ config manager, kiểm tra dependencies cần thiết
// và thiết lập các driver theo cấu hình.
//
// Params:
//   - app: Interface của ứng dụng, phải cung cấp các phương thức Container(), Environment(), và BasePath()
func (p *serviceProvider) Boot(app interface{}) {
	if app == nil {
		return // Không cần xử lý nếu app không implement Container()
	}
}

// Providers trả về danh sách các dịch vụ được cung cấp bởi provider.
//
// Trả về:
//   - []string: danh sách các khóa dịch vụ mà provider này cung cấp
func (p *serviceProvider) Providers() []string {
	return p.providers
}
