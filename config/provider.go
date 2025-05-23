// Package config cung cấp giải pháp quản lý cấu hình linh hoạt và mở rộng cho ứng dụng Go.
package config

import (
	"github.com/go-fork/di"
	"github.com/go-fork/providers/config/formatter"
)

// ServiceProvider cung cấp cấu hình cho module config và tích hợp với DI container.
//
// ServiceProvider là một implementation của interface di.ServiceProvider, cho phép tự động
// đăng ký config manager vào DI container của ứng dụng. ServiceProvider thực hiện các công việc:
//   - Tạo một config manager mới
//   - Nạp cấu hình từ biến môi trường với prefix "APP"
//   - Nạp cấu hình từ thư mục configs của ứng dụng (các file YAML)
//   - Đăng ký config manager vào DI container với key "config"
//
// Để sử dụng ServiceProvider, ứng dụng cần:
//   - Implement interface Container() *di.Container để cung cấp DI container
//   - Implement interface BasePath(path ...string) string để cung cấp đường dẫn đến thư mục configs
type ServiceProvider struct{}

// NewServiceProvider trả về một ServiceProvider mới cho module config.
//
// Hàm này khởi tạo và trả về một đối tượng ServiceProvider để sử dụng với DI container.
// ServiceProvider cho phép tự động đăng ký và cấu hình config manager cho ứng dụng.
//
// Returns:
//   - di.ServiceProvider: Interface di.ServiceProvider đã được implement bởi ServiceProvider
//
// Example:
//
//	app.Register(config.NewServiceProvider())
func NewServiceProvider() di.ServiceProvider {
	return &ServiceProvider{}
}

// Register đăng ký các binding cấu hình vào DI container.
//
// Register được gọi khi đăng ký ServiceProvider vào ứng dụng. Phương thức này
// tạo một config manager mới, nạp cấu hình từ các nguồn khác nhau,
// và đăng ký manager này vào DI container của ứng dụng.
//
// Params:
//   - app: interface{} - Đối tượng ứng dụng phải implement các interface:
//   - Container() *di.Container - Trả về DI container
//   - BasePath(path ...string) string - Trả về đường dẫn gốc của ứng dụng
//
// Các formatters được đăng ký theo thứ tự ưu tiên từ cao xuống thấp:
//  1. EnvFormatter: Đọc biến môi trường với prefix "APP_"
//  2. JSONFormatter: Đọc file JSON từ thư mục configs
//  3. YAMLFormatter: Đọc file YAML từ thư mục configs
//
// Luồng thực thi:
//  1. Kiểm tra app có implement Container() không, nếu không thì return
//  2. Lấy container từ app, kiểm tra nếu nil thì return
//  3. Tạo config manager mới
//  4. Nếu app implement BasePath():
//     a. Nạp cấu hình từ biến môi trường với prefix "APP"
//     b. Nạp cấu hình từ thư mục configs của ứng dụng
//  5. Đăng ký config manager vào container với key "config"
func (p *ServiceProvider) Register(app interface{}) {
	// Lấy container từ app nếu có
	appWithContainer, ok := app.(interface{ Container() *di.Container })
	if !ok {
		return // Không cần xử lý nếu app không implement Container()
	}

	container := appWithContainer.Container()
	if container == nil {
		return // Không xử lý nếu container nil
	}

	manager := NewConfig()

	// Kiểm tra app có implement BasePath không
	appWithBasePath, ok := app.(interface{ BasePath(path ...string) string })
	if ok {
		// Nạp cấu hình từ biến môi trường với prefix "APP_"
		envFormatter := formatter.NewEnvFormatter("APP_")
		if err := manager.Load(envFormatter); err != nil {
			// Tiếp tục ngay cả khi có lỗi từ env formatter
		}

		// Lấy đường dẫn tới thư mục configs
		configPath := appWithBasePath.BasePath("configs")

		// Nạp cấu hình từ file YAML (ưu tiên thấp nhất)
		yamlFormatter := formatter.NewYAMLFormatter(configPath + "/config.yaml")
		if err := manager.Load(yamlFormatter); err != nil {
			// Tiếp tục ngay cả khi có lỗi từ yaml formatter
		}

		// Nạp cấu hình từ file JSON (ưu tiên cao hơn YAML)
		jsonFormatter := formatter.NewJSONFormatter(configPath + "/config.json")
		if err := manager.Load(jsonFormatter); err != nil {
			// Tiếp tục ngay cả khi có lỗi từ json formatter
		}
	}

	// Đăng ký config manager vào container
	container.Instance("config", manager)
}

// Boot được gọi sau khi tất cả các service provider đã được đăng ký.
//
// Boot là một lifecycle hook của di.ServiceProvider mà thực hiện sau khi tất cả
// các service provider đã được đăng ký xong. Trong trường hợp của ConfigServiceProvider,
// không cần thực hiện thêm hành động nào ở giai đoạn boot.
//
// Params:
//   - app: interface{} - Đối tượng ứng dụng, không sử dụng trong phương thức này
func (p *ServiceProvider) Boot(app interface{}) {
	// Không cần thực hiện gì trong Boot
}
