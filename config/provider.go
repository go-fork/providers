package config

import (
	"go.fork.vn/di"
)

// ServiceProvider cung cấp cấu hình cho module config và tích hợp với DI container.
//
// ServiceProvider là một implementation của interface di.ServiceProvider, cho phép tự động
// đăng ký config manager vào DI container của ứng dụng. ServiceProvider thực hiện công việc:
//   - Tạo một config manager mới sử dụng Viper
//   - Đăng ký config manager vào DI container với key "config"
//
// Việc cấu hình cụ thể (đọc file, biến môi trường, v.v.) được thực hiện bởi ứng dụng.
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
// tạo một config manager mới và đăng ký vào DI container của ứng dụng.
//
// Params:
//   - app: interface{} - Đối tượng ứng dụng phải implement interface:
//     Container() *di.Container - Trả về DI container
//
// Luồng thực thi:
//  1. Kiểm tra app có implement Container() không, nếu không thì return
//  2. Lấy container từ app, kiểm tra nếu nil thì panic
//  3. Tạo config manager mới
//  4. Đăng ký config manager vào container với key "config"
//
// Việc cấu hình (đọc từ file, biến môi trường, vv) sẽ được thực hiện bởi ứng dụng,
// cho phép mỗi ứng dụng tùy chỉnh cấu hình theo nhu cầu riêng.
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

	// Tạo một config manager mới và đăng ký vào container
	manager := NewConfig()
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
	// Safety check, though method is a no-op
	if app == nil {
		return
	}
}

// Requires trả về danh sách các provider mà provider này phụ thuộc vào.
//
// ConfigServiceProvider không phụ thuộc vào bất kỳ provider nào khác,
// nên phương thức này trả về một slice rỗng.
//
// Returns:
//   - []string: Một slice rỗng vì không có dependencies
func (p *ServiceProvider) Requires() []string {
	return []string{} // Không có dependencies
}

// Providers trả về danh sách các service mà provider này đăng ký.
//
// ConfigServiceProvider đăng ký config manager vào container với key "config"
//
// Returns:
//   - []string: Mảng chứa tên của service được đăng ký - "config"
func (p *ServiceProvider) Providers() []string {
	return []string{"config"}
}
