// Package scheduler cung cấp giải pháp lên lịch và chạy các task định kỳ cho ứng dụng Go.
package scheduler

import (
	"github.com/go-fork/di"
	"github.com/go-fork/providers/config"
	"github.com/go-fork/providers/redis"
)

// ServiceProvider cung cấp dịch vụ scheduler và tích hợp với DI container.
//
// ServiceProvider là một implementation của interface di.ServiceProvider, cho phép tự động
// đăng ký scheduler manager vào DI container của ứng dụng. ServiceProvider thực hiện công việc:
//   - Tạo một scheduler manager mới sử dụng gocron
//   - Đăng ký scheduler manager vào DI container với key "scheduler"
//
// Việc cấu hình cụ thể và đăng ký các task được thực hiện bởi ứng dụng.
//
// Để sử dụng ServiceProvider, ứng dụng cần:
//   - Implement interface Container() *di.Container để cung cấp DI container
type ServiceProvider struct {
	providers []string
}

// NewServiceProvider trả về một ServiceProvider mới cho module scheduler.
//
// Hàm này khởi tạo và trả về một đối tượng ServiceProvider để sử dụng với DI container.
// ServiceProvider cho phép tự động đăng ký và cấu hình scheduler manager cho ứng dụng.
//
// Returns:
//   - di.ServiceProvider: Interface di.ServiceProvider đã được implement bởi ServiceProvider
//
// Example:
//
//	app.Register(scheduler.NewServiceProvider())
func NewServiceProvider() di.ServiceProvider {
	return &ServiceProvider{}
}

// Register đăng ký scheduler vào DI container.
//
// Register được gọi khi đăng ký ServiceProvider vào ứng dụng. Phương thức này
// tạo một scheduler manager mới và đăng ký vào DI container của ứng dụng.
//
// Params:
//   - app: interface{} - Đối tượng ứng dụng phải implement interface:
//     Container() *di.Container - Trả về DI container
//
// Luồng thực thi:
//  1. Kiểm tra app có implement Container() không, nếu không thì return
//  2. Lấy container từ app, kiểm tra nếu nil thì panic
//  3. Tạo scheduler manager mới
//  4. Đăng ký scheduler manager vào container với key "scheduler"
//
// Việc cấu hình và đăng ký các task sẽ được thực hiện bởi ứng dụng,
// cho phép mỗi ứng dụng tùy chỉnh scheduler theo nhu cầu riêng.
func (p *ServiceProvider) Register(app interface{}) {
	// Lấy container từ app nếu có
	appWithContainer, ok := app.(interface{ Container() *di.Container })
	if !ok {
		return // Không cần xử lý nếu app không implement Container()
	}
	container := appWithContainer.Container()
	if container == nil {
		panic("DI container is nil")
	}

	// Load cấu hình scheduler với default fallback
	cfg := DefaultConfig()

	// Thử lấy cấu hình từ config provider (optional)
	if configInstance, err := container.Make("config"); err == nil {
		if configManager, ok := configInstance.(config.Manager); ok {
			// Load cấu hình từ file config, nếu lỗi thì sử dụng default
			configManager.UnmarshalKey("scheduler", &cfg)
		}
	}

	// Tạo scheduler manager với cấu hình
	manager := NewSchedulerWithConfig(cfg)

	// Cấu hình distributed locking nếu được bật
	if cfg.DistributedLock.Enabled {
		if redisInstance, err := container.Make("redis"); err == nil {
			// Thử lấy redis client từ redis provider
			if redisManager, ok := redisInstance.(redis.Manager); ok {
				if redisClient, err := redisManager.Client(); err == nil {
					if locker, err := NewRedisLocker(redisClient, cfg.Options); err == nil {
						manager = manager.WithDistributedLocker(locker)
					}
				}
			}
		}
	}

	// Đăng ký scheduler manager vào container
	container.Instance("scheduler", manager)
	p.providers = append(p.providers, "scheduler")
}

// Boot được gọi sau khi tất cả các service provider đã được đăng ký.
//
// Boot là một lifecycle hook của di.ServiceProvider mà thực hiện sau khi tất cả
// các service provider đã được đăng ký xong.
//
// Trong trường hợp của SchedulerServiceProvider, có thể dùng Boot để:
// 1. Lấy scheduler manager từ container
// 2. Bắt đầu scheduler trong chế độ async để nó sẵn sàng xử lý các task
//
// Params:
//   - app: interface{} - Đối tượng ứng dụng phải implement interface:
//     Container() *di.Container - Trả về DI container
func (p *ServiceProvider) Boot(app interface{}) {
	// Lấy container từ app nếu có
	appWithContainer, ok := app.(interface{ Container() *di.Container })
	if !ok {
		return // Không cần xử lý nếu app không implement Container()
	}
	container := appWithContainer.Container()
	if container == nil {
		return // Không xử lý nếu container nil
	}

	// Lấy scheduler manager từ container
	instance, err := container.Make("scheduler")
	if err != nil {
		return // Không tìm thấy scheduler trong container
	}
	scheduler, ok := instance.(Manager)
	if !ok {
		return // Không phải loại scheduler manager
	}

	// Kiểm tra xem scheduler đã được start chưa
	if scheduler.IsRunning() {
		return // Scheduler đã được start rồi
	}

	// Lấy cấu hình để kiểm tra AutoStart
	cfg := DefaultConfig()

	// Thử lấy cấu hình từ config provider (optional)
	if configInstance, err := container.Make("config"); err == nil {
		if configManager, ok := configInstance.(config.Manager); ok {
			// Load cấu hình từ file config, nếu lỗi thì sử dụng default
			configManager.UnmarshalKey("scheduler", &cfg)
		}
	}

	// Chỉ auto start nếu được cấu hình để làm vậy
	if cfg.AutoStart {
		scheduler.StartAsync()
	}
}

func (p *ServiceProvider) Requires() []string {
	return []string{
		"config",
		"redis",
	}
}

func (p *ServiceProvider) Providers() []string {
	return p.providers
}
