package queue

import (
	"github.com/go-fork/di"
	"github.com/go-fork/providers/config"
)

// ServiceProvider triển khai interface di.ServiceProvider cho các dịch vụ queue.
//
// Provider này tự động hóa việc đăng ký các dịch vụ queue trong một container
// dependency injection, thiết lập client và server với các giá trị mặc định hợp lý.
type ServiceProvider struct {
	config Config
}

// NewServiceProvider tạo một provider dịch vụ queue mới với cấu hình mặc định.
//
// Sử dụng hàm này để tạo một provider có thể được đăng ký với
// một instance di.Container.
//
// Trả về:
//   - di.ServiceProvider: một service provider cho queue
//
// Ví dụ:
//
//	app := myapp.New()
//	app.Register(queue.NewServiceProvider())
func NewServiceProvider() di.ServiceProvider {
	return &ServiceProvider{
		config: DefaultConfig(),
	}
}

// NewServiceProviderWithConfig tạo một provider dịch vụ queue mới với cấu hình tùy chỉnh.
//
// Sử dụng hàm này để tạo một provider với cấu hình tùy chỉnh có thể được đăng ký với
// một instance di.Container.
//
// Tham số:
//   - config: Config - cấu hình cho dịch vụ queue
//
// Trả về:
//   - di.ServiceProvider: một service provider cho queue
func NewServiceProviderWithConfig(config Config) di.ServiceProvider {
	return &ServiceProvider{
		config: config,
	}
}

// Register đăng ký các dịch vụ queue với container của ứng dụng.
//
// Phương thức này:
//   - Tạo một queue manager
//   - Đăng ký client, server và các thành phần cần thiết khác
//   - Đăng ký tất cả vào container DI
//
// Tham số:
//   - app: interface{} - instance của ứng dụng cung cấp Container()
//
// Ứng dụng phải triển khai:
//   - Container() *di.Container
func (p *ServiceProvider) Register(app interface{}) {
	// Trích xuất container từ ứng dụng
	if appWithContainer, ok := app.(interface {
		Container() *di.Container
	}); ok {
		c := appWithContainer.Container()

		// Kiểm tra xem container đã có config manager chưa
		configInstance, err := c.Make("config")
		if err == nil {
			// Nếu có config manager, thử lấy cấu hình queue
			if configManager, ok := configInstance.(config.Manager); ok {
				// Cố gắng load cấu hình từ config manager
				var queueConfig Config
				if configManager.Has("queue") {
					if err := configManager.UnmarshalKey("queue", &queueConfig); err == nil {
						// Nếu load thành công, sử dụng cấu hình đã load
						p.config = queueConfig
					}
				}
			}
		}

		// Tạo một queue manager mới
		manager := NewManagerWithConfig(p.config)

		// Đăng ký các thành phần chính
		c.Instance("queue", manager)                 // Dịch vụ queue manager chung
		c.Instance("queue.client", manager.Client()) // Binding đặc biệt cho client
		c.Instance("queue.server", manager.Server()) // Binding đặc biệt cho server
		c.Instance("queue.manager", manager)         // Binding đặc biệt cho manager

		// Đăng ký các adapter và các thành phần phụ thuộc
		if p.config.Adapter.Default == "redis" {
			c.Instance("queue.redis", manager.RedisClient()) // Redis client
		}
	}
}

// Boot thực hiện thiết lập sau đăng ký cho dịch vụ queue.
//
// Đối với provider queue, hiện tại đây là no-op vì tất cả thiết lập
// được thực hiện trong quá trình đăng ký.
//
// Tham số:
//   - app: interface{} - instance của ứng dụng
func (p *ServiceProvider) Boot(app interface{}) {
	// Không yêu cầu thiết lập bổ sung sau khi đăng ký
	if app == nil {
		return
	}
}
