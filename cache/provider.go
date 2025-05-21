// Package cache cung cấp framework cache với nhiều driver cho ứng dụng Go.
//
// Package này triển khai một hệ thống cache đa driver, thread-safe và extensible cho các ứng dụng Go.
// Nó hỗ trợ các driver như Memory, File, Redis, MongoDB, và có thể mở rộng với các driver tùy chỉnh.
// Package cung cấp các phương thức tiêu chuẩn như Get, Set, Delete, Flush, và các chức năng nâng cao như
// Remember, GetMultiple, SetMultiple để tối ưu hiệu suất tương tác với cache.
package cache

import (
	"github.com/go-fork/di"
	"github.com/go-fork/providers/cache/driver"
)

// ServiceProvider là service provider cho module cache.
//
// ServiceProvider đảm nhận việc đăng ký và khởi tạo các dịch vụ cache vào container DI.
// Nó cung cấp cơ chế để đăng ký cache manager và các driver mặc định vào ứng dụng.
type ServiceProvider struct{}

// NewServiceProvider tạo một cache service provider mới.
//
// Phương thức này khởi tạo một service provider mới cho module cache.
//
// Returns:
//   - di.ServiceProvider: Đối tượng service provider đã sẵn sàng đăng ký vào ứng dụng
func NewServiceProvider() di.ServiceProvider {
	return &ServiceProvider{}
}

// Register đăng ký các binding vào container.
//
// Phương thức này đăng ký cache manager vào container DI của ứng dụng.
// Nó khởi tạo một cache manager mới và đăng ký nó dưới các alias "cache" và "cache.manager".
//
// Params:
//   - app: Interface của ứng dụng, phải cung cấp phương thức Container() để lấy container DI
func (p *ServiceProvider) Register(app interface{}) {
	// Lấy container từ app
	if appWithContainer, ok := app.(interface {
		Container() *di.Container
	}); ok {
		c := appWithContainer.Container()

		// Tạo cache manager
		manager := NewManager()

		// Đăng ký cache service
		c.Instance("cache", manager)
		c.Instance("cache.manager", manager)
	}
}

// Boot được gọi sau khi tất cả các service provider đã được đăng ký.
//
// Phương thức này khởi tạo và cấu hình các cache driver mặc định cho ứng dụng.
// Nó tự động thiết lập memory driver làm mặc định và cố gắng khởi tạo file driver nếu có thể.
// Boot chỉ được gọi sau khi tất cả các service provider đã được đăng ký thông qua Register().
//
// Params:
//   - app: Interface của ứng dụng, phải cung cấp các phương thức Container(), Environment(), và BasePath()
func (p *ServiceProvider) Boot(app interface{}) {
	// Có thể cấu hình các driver mặc định dựa trên môi trường hoặc cấu hình
	if appWithContainer, ok := app.(interface {
		Container() *di.Container
		Environment() string
		BasePath(paths ...string) string
	}); ok {
		c := appWithContainer.Container()

		if cacheManager, err := c.Make("cache"); err == nil {
			if manager, ok := cacheManager.(Manager); ok {
				// Thêm memory driver mặc định
				manager.AddDriver("memory", driver.NewMemoryDriver())

				// Thêm file driver nếu cần
				storagePath := appWithContainer.BasePath("storage", "cache")
				fileDriver, err := driver.NewFileDriver(storagePath)
				if err == nil {
					manager.AddDriver("file", fileDriver)
				}

				// Thiết lập driver mặc định
				manager.SetDefaultDriver("memory")
			}
		}
	}
}
