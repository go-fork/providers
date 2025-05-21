package log

import (
	"os"
	"path/filepath"

	"github.com/go-fork/di"
	"github.com/go-fork/providers/log/handler"
)

// ServiceProvider triển khai interface di.ServiceProvider cho các dịch vụ logging.
//
// Provider này tự động hóa việc đăng ký các dịch vụ logging trong một container
// dependency injection, thiết lập các handlers cho console và file với các giá trị mặc định hợp lý.
type ServiceProvider struct{}

// NewServiceProvider tạo một provider dịch vụ log mới.
//
// Sử dụng hàm này để tạo một provider có thể được đăng ký với
// một instance di.Container.
//
// Trả về:
//   - di.ServiceProvider: một service provider cho logging
//
// Ví dụ:
//
//	app := myapp.New()
//	app.Register(log.NewServiceProvider())
func NewServiceProvider() di.ServiceProvider {
	return &ServiceProvider{}
}

// Register đăng ký các dịch vụ logging với container của ứng dụng.
//
// Phương thức này:
//   - Tạo một log manager
//   - Cấu hình một console handler có màu
//   - Thiết lập một file handler trong thư mục storage/logs
//   - Đăng ký manager trong container DI
//
// Tham số:
//   - app: interface{} - instance của ứng dụng cung cấp Container() và BasePath()
//
// Ứng dụng phải triển khai:
//   - Container() *di.Container
//   - BasePath(...string) string
func (p *ServiceProvider) Register(app interface{}) {
	// Trích xuất container và đường dẫn cơ sở từ ứng dụng
	if appWithContainer, ok := app.(interface {
		Container() *di.Container
		BasePath(paths ...string) string
	}); ok {
		c := appWithContainer.Container()

		// Tạo một log manager mới
		manager := NewManager()

		// Thêm một console handler có màu
		consoleHandler := handler.NewConsoleHandler(true)
		manager.AddHandler("console", consoleHandler)

		// Cấu hình đường dẫn lưu trữ cho các file log
		storagePath := appWithContainer.BasePath("storage", "logs")
		if _, err := os.Stat(storagePath); os.IsNotExist(err) {
			// Tạo thư mục logs nếu nó không tồn tại
			os.MkdirAll(storagePath, 0755)
		}

		// Thêm một file handler cho việc ghi log liên tục
		// Lưu ý: maxSize được đặt thành 10 bytes cho mục đích demo
		// Trong môi trường production, sử dụng giá trị lớn hơn như 10*1024*1024 (10MB)
		fileHandler, err := handler.NewFileHandler(filepath.Join(storagePath, "app.log"), 10)
		if err == nil {
			manager.AddHandler("file", fileHandler)
		}

		// Đăng ký log manager trong container
		c.Instance("log", manager)         // Dịch vụ logging chung
		c.Instance("log.manager", manager) // Binding đặc biệt cho manager
	}
}

// Boot thực hiện thiết lập sau đăng ký cho dịch vụ logging.
//
// Đối với provider logging, hiện tại đây là no-op vì tất cả thiết lập
// được thực hiện trong quá trình đăng ký.
//
// Tham số:
//   - app: interface{} - instance của ứng dụng
func (p *ServiceProvider) Boot(app interface{}) {
	// Không yêu cầu thiết lập bổ sung sau khi đăng ký
}
