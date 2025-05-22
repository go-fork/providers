// Package config cung cấp giải pháp quản lý cấu hình linh hoạt và mở rộng cho ứng dụng Go.
package config

import (
	"fmt"
	"os"

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

	// Tạo config manager
	manager := NewManager()

	// Nạp cấu hình từ env và thư mục configs nếu app hỗ trợ BasePath
	if appWithPath, ok := app.(interface{ BasePath(path ...string) string }); ok {
		// Nạp từ biến môi trường
		envProvider := formatter.NewEnvFormatter("APP")
		_ = manager.Load(envProvider) // Bỏ qua lỗi khi nạp từ env

		// Nạp từ file YAML trong thư mục configs
		configPath := appWithPath.BasePath("configs")
		if configPath != "" {
			// Kiểm tra thư mục configs tồn tại
			if fileInfo, err := os.Stat(configPath); err == nil && fileInfo.IsDir() {
				// Phải cẩn thận xử lý lỗi từ LoadFromDirectory vì nó bây giờ trả về lỗi cho file YAML không hợp lệ
				values, err := formatter.LoadFromDirectory(configPath)
				if err == nil {
					for k, v := range values {
						if k != "" {
							manager.Set(k, v)
						}
					}
				} else {
					// Ghi log lỗi nhưng không dừng quá trình đăng ký
					fmt.Printf("Warning: Failed to load config from directory %s: %v\n", configPath, err)
				}
			}
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
