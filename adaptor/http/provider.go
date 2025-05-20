package adapter

import (
	"github.com/go-fork/pamm/pkg/di/container"
	"github.com/go-fork/pamm/pkg/infra/config"
)

// ServiceProvider là service provider cho HTTP module
type ServiceProvider struct{}

// NewServiceProvider tạo một HTTP service provider mới
func NewServiceProvider() container.ServiceProvider {
	return &ServiceProvider{}
}

// Register đăng ký các binding vào container
func (p *ServiceProvider) Register(app interface{}) {
	// Lấy container từ app
	if appWithContainer, ok := app.(interface {
		Container() *container.Container
	}); ok {
		c := appWithContainer.Container()

		// Lấy config từ container
		var cfg *Config
		if c.Bound("config") {
			configManager := c.MustMake("config").(config.Manager)
			if configManager.Has(ConfigPrefix) {
				if err := configManager.Unmarshal(ConfigPrefix, cfg); err != nil {
					cfg = DefaultConfig()
				}
			}
		}
		if cfg == nil {
			// Sử dụng cấu hình mặc định
			cfg = DefaultConfig()
		}
		// Tạo HTTP application
		apdater := NewNetHTTPAdapter(cfg)
		// Đăng ký HTTP Adapter vào container
		c.Instance("http.http", apdater)
	}
}

// Boot được gọi sau khi tất cả các service provider đã được đăng ký
func (p *ServiceProvider) Boot(app interface{}) {
	// Không cần thực hiện gì trong Boot
}
