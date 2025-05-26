package redis

import (
	"github.com/go-fork/di"
	"github.com/go-fork/providers/config"
)

// ServiceProvider định nghĩa interface cho Redis service provider.
//
// ServiceProvider kế thừa từ di.ServiceProvider và định nghĩa
// các phương thức cần thiết cho một Redis service provider.
type ServiceProvider interface {
	di.ServiceProvider
}

// serviceProvider là implementation của ServiceProvider.
//
// serviceProvider chịu trách nhiệm đăng ký các dịch vụ Redis vào DI container
// và cung cấp Redis client cho các module khác trong ứng dụng.
type serviceProvider struct {
	providers []string
}

// NewServiceProvider tạo một Redis service provider mới.
func NewServiceProvider() ServiceProvider {
	return &serviceProvider{}
}

// Register đăng ký các dịch vụ Redis với DI container.
//
// Phương thức này đăng ký Redis manager vào container DI của ứng dụng.
// Nó khởi tạo một Redis manager mới và đăng ký nó dưới các alias "redis.manager",
// "redis.client" và "redis.universal".
//
// Params:
//   - app: Interface của ứng dụng, phải cung cấp phương thức Container() để lấy container DI
func (p *serviceProvider) Register(app interface{}) {
	// Lấy container từ app
	if appWithContainer, ok := app.(interface {
		Container() *di.Container
	}); ok {
		c := appWithContainer.Container()

		// Kiểm tra xem container có tồn tại không
		if c == nil {
			// Không làm gì khi không có container
			return
		}

		// Kiểm tra xem container đã có config manager chưa
		redisConfig := DefaultConfig()
		configService := c.MustMake("config").(config.Manager)
		if configService == nil {
			panic("Redis provider requires config service to be registered")
		}
		err := configService.UnmarshalKey("redis", &redisConfig)
		if err != nil {
			panic("Redis config unmarshal error: " + err.Error())
		}

		manager := NewManagerWithConfig(redisConfig)

		client, err := manager.Client()
		if err == nil {
			c.Instance("redis.client", client)
			p.providers = append(p.providers, "redis.client")
		}

		universalClient, err := manager.UniversalClient()
		if err == nil {
			c.Instance("redis.universal", universalClient)
			p.providers = append(p.providers, "redis.universal")
		}
		c.Instance("redis.manager", manager)
		p.providers = append(p.providers, "redis")
	}
}

// Boot khởi động Redis provider.
//
// Phương thức này khởi động Redis provider sau khi tất cả các service provider đã được đăng ký.
// Trong trường hợp này, không cần thực hiện thêm tác vụ nào trong Boot vì các cấu hình
// đã được xử lý trong Register.
//
// Params:
//   - app: Interface của ứng dụng
func (p *serviceProvider) Boot(app interface{}) {
	// Không cần thực hiện thêm tác vụ nào trong Boot
	// vì cấu hình đã được xử lý trong Register
	if app == nil {
		return
	}
}

// Providers trả về danh sách các service được cung cấp bởi Redis provider.
//
// Phương thức này trả về danh sách các abstract type mà Redis provider đăng ký với container.
// Danh sách này được sử dụng để kiểm tra dependencies và đảm bảo đúng thứ tự khởi tạo.
//
// Trả về:
//   - []string: danh sách các service được cung cấp
func (p *serviceProvider) Providers() []string {
	return p.providers
}

// Requires trả về danh sách các dependency mà Redis provider phụ thuộc.
//
// Trả về:
//   - []string: danh sách các service provider khác mà provider này yêu cầu
func (p *serviceProvider) Requires() []string {
	return []string{
		// Redis provider có thể làm việc độc lập nhưng sẽ sử dụng config provider nếu có
		"config",
	}
}
