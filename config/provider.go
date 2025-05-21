package config

import (
	"fmt"
	"os"

	"github.com/go-fork/di"
	"github.com/go-fork/providers/config/formatter"
)

// ServiceProvider cung cấp cấu hình cho module config và tích hợp với DI container.
//
// ServiceProvider đăng ký config manager vào container, tự động nạp cấu hình từ biến môi trường và thư mục configs.
type ServiceProvider struct{}

// NewServiceProvider trả về một ServiceProvider mới cho module config.
//
// Hàm này khởi tạo và trả về một đối tượng ServiceProvider để sử dụng với DI container.
func NewServiceProvider() di.ServiceProvider {
	return &ServiceProvider{}
}

// Register đăng ký các binding cấu hình vào DI container.
//
// Tham số app phải implement interface có Container() *di.Container và BasePath().
// Hàm này sẽ tạo config manager, nạp cấu hình từ biến môi trường và thư mục configs, sau đó đăng ký vào container.
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
// Hàm này là một lifecycle hook, mặc định không thực hiện gì.
func (p *ServiceProvider) Boot(app interface{}) {
	// Không cần thực hiện gì trong Boot
}
