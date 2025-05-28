package mailer

import (
	"context"

	"github.com/go-fork/di"
	"github.com/go-fork/providers/config"
	"github.com/go-fork/providers/queue"
)

// ServiceProvider triển khai interface di.ServiceProvider cho dịch vụ gửi mail.
//
// Provider này tự động hóa việc đăng ký các dịch vụ mail và queue trong container
// dependency injection, thiết lập các giá trị mặc định hợp lý.
type ServiceProvider struct{}

// NewServiceProvider tạo một provider dịch vụ mail mới với cấu hình mặc định.
//
// Sử dụng hàm này để tạo một provider có thể được đăng ký với
// một instance di.Container.
//
// Trả về:
//   - di.ServiceProvider: một service provider cho mailer
//
// Ví dụ:
//
//	app := myapp.New()
//	app.Register(mailer.NewServiceProvider())
func NewServiceProvider() di.ServiceProvider {
	return &ServiceProvider{}
}

// Register đăng ký các dịch vụ mailer với container của ứng dụng.
//
// Phương thức này sẽ đăng ký:
//   - Manager: quản lý các thành phần của mailer
//   - Mailer: dịch vụ gửi mail
//
// Tham số:
//   - app: interface{} - instance của ứng dụng cung cấp Container()
func (p *ServiceProvider) Register(app interface{}) {
	// Trích xuất container từ ứng dụng
	if appWithContainer, ok := app.(interface {
		Container() *di.Container
	}); ok {
		c := appWithContainer.Container()

		// Kiểm tra xem container đã có config manager chưa
		configInstance, err := c.Make("config")
		if err == nil {
			// Nếu có config manager, thử lấy cấu hình mailer
			if configManager, ok := configInstance.(config.Manager); ok {
				mailerConfig, err := LoadConfig(configManager)
				if err != nil {
					// Nếu có lỗi khi load cấu hình, sử dụng cấu hình mặc định
					panic("Please configure mailer in config: " + err.Error())
				}
				// Đăng ký manager
				c.Bind("mailer.manager", func(c *di.Container) interface{} {
					manager, _ := NewManager(mailerConfig)
					return manager
				})

				// Đăng ký mailer service
				c.Bind("mailer", func(c *di.Container) interface{} {
					manager := c.MustMake("mailer.manager").(Manager)
					return manager.Mailer()
				})
			}
		}
	}
}

// Boot thực hiện các thiết lập cần thiết sau khi đăng ký dịch vụ.
//
// Phương thức này sẽ:
//   - Thiết lập xử lý queue tasks cho mail nếu queue được bật
//
// Tham số:
//   - app: interface{} - instance của ứng dụng cung cấp Container()
func (p *ServiceProvider) Boot(app interface{}) {
	// Trích xuất container từ ứng dụng
	if appWithContainer, ok := app.(interface {
		Container() *di.Container
	}); ok {
		c := appWithContainer.Container()

		// Lấy manager từ container
		managerInstance, err := c.Make("mailer.manager")
		if err != nil {
			return
		}
		manager := managerInstance.(Manager)

		// Nếu queue được bật, đăng ký xử lý mail tasks
		if manager.QueueEnabled() {
			queueManager := manager.QueueManager()
			if queueManager != nil {
				server := queueManager.Server()

				// Đăng ký handler cho task gửi mail
				server.RegisterHandler("mailer:send", func(ctx context.Context, task *queue.Task) error {
					return manager.ProcessMessage(task.Payload)
				})
			}
		}
	}
}

// Providers trả về danh sách các dịch vụ được cung cấp bởi provider.
//
// Trả về:
//   - []string: danh sách các khóa dịch vụ mà provider này cung cấp
func (p *ServiceProvider) Providers() []string {
	return []string{
		"mailer.manager",
		"mailer",
	}
}

// Requires trả về danh sách các provider mà mailer provider phụ thuộc vào.
//
// Mailer provider phụ thuộc vào config provider để đọc cấu hình mail.
// Nếu queue được bật, mailer cũng phụ thuộc vào queue provider.
//
// Trả về:
//   - []string: danh sách các provider mà mailer phụ thuộc vào
func (p *ServiceProvider) Requires() []string {
	return []string{
		"config",
		// queue là optional dependency (chỉ cần khi queue.enabled = true)
	}
}
