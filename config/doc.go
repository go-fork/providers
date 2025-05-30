// Package config là thư viện quản lý cấu hình (configuration) hiện đại, linh hoạt
// và mở rộng cho ứng dụng Go, dựa trên nền tảng thư viện Viper.
//
// # Giới thiệu
//
// Package config cung cấp wrapper cho thư viện Viper nổi tiếng, đồng thời mở rộng tính năng
// và chuẩn hóa API để dễ dàng tích hợp vào các ứng dụng thông qua Dependency Injection.
// Thư viện này thiết kế nhằm tối ưu quy trình quản lý cấu hình, đảm bảo tính nhất quán
// và linh hoạt cho các ứng dụng Go hiện đại.
//
// # Đối tượng chính
//
//   - Manager: Interface chính để tương tác với hệ thống cấu hình, cung cấp các phương thức
//     để truy xuất và quản lý cấu hình, đồng thời bổ sung các tiện ích so với Viper gốc.
//
//   - ServiceProvider: Implementation của di.ServiceProvider để tích hợp với DI container,
//     cho phép đăng ký và cấu hình tự động cho Manager. Tuân thủ đầy đủ interface di.ServiceProvider
//     với các phương thức Register(), Boot(), Requires(), và Providers().
//
// # Tính năng nổi bật
//
// - Hỗ trợ đa dạng định dạng cấu hình: YAML, JSON, TOML, HCL, INI, properties, dotenv
// - Tự động đọc từ biến môi trường (Environment Variables)
// - Hot reload: tự động cập nhật khi file cấu hình thay đổi
// - Phân cấp cấu hình với hỗ trợ đầy đủ cho dot notation (ví dụ: "database.host")
// - Kiểu dữ liệu phong phú với hỗ trợ unmarshalling vào struct
// - Cấu hình mặc định và override thông qua nhiều nguồn
// - API đồng nhất, safe type với pattern return value, ok cho mọi kiểu dữ liệu
// - Tích hợp với DI container thông qua ServiceProvider
//
// # Kiến trúc và cách hoạt động
//
//   - Sử dụng mô hình composition thông qua embedded struct để kế thừa tính năng của Viper
//   - Interface Manager định nghĩa API chuẩn cho việc truy xuất và quản lý cấu hình
//   - manager struct cài đặt interface Manager bằng cách nhúng *viper.Viper
//   - ServiceProvider cài đặt đầy đủ di.ServiceProvider interface, cung cấp các phương thức:
//   - Register(): Đăng ký Manager vào DI container với key "config"
//   - Boot(): Khởi tạo cấu hình sau khi đăng ký
//   - Requires(): Trả về danh sách các provider mà config phụ thuộc (hiện tại không có)
//   - Providers(): Trả về danh sách các service mà config cung cấp (key "config")
//   - Các phương thức của Manager đều được bổ sung kiểm tra tồn tại (IsSet) trước khi truy xuất,
//     trả về giá trị mặc định phù hợp và boolean cho biết key có tồn tại không
//
// # Ví dụ sử dụng cơ bản
//
//	// Đăng ký với DI container
//	app := di.New()
//	app.Register(config.NewServiceProvider())
//
//	// Lấy config từ container
//	container := app.Container()
//	cfg := container.Get("config").(config.Manager)
//
//	// Cấu hình và đọc file
//	cfg.SetConfigFile("config.yaml")
//	err := cfg.ReadInConfig()
//	if err != nil {
//	    log.Fatalf("Không thể đọc file cấu hình: %v", err)
//	}
//
//	// Sử dụng cấu hình với kiểm tra tồn tại
//	if name, ok := cfg.GetString("app.name"); ok {
//	    fmt.Printf("Tên ứng dụng: %s\n", name)
//	}
//
//	if port, ok := cfg.GetInt("app.port"); ok {
//	    fmt.Printf("Cổng: %d\n", port)
//	}
//
//	// Unmarshalling vào struct
//	type DatabaseConfig struct {
//	    Host     string
//	    Port     int
//	    Username string
//	    Password string
//	}
//
//	var dbConfig DatabaseConfig
//	err = cfg.UnmarshalKey("database", &dbConfig)
//	if err != nil {
//	    log.Fatalf("Lỗi unmarshalling cấu hình database: %v", err)
//	}
//
// Package này giúp nhất quán hóa quản lý cấu hình trong ứng dụng Go, tận dụng
// sức mạnh của Viper và bổ sung các tính năng hữu ích cho ứng dụng hiện đại.
//
// Module này tương thích đầy đủ với go.fork.vn/di từ phiên bản v0.0.5 trở lên,
// cài đặt đầy đủ interface ServiceProvider với các phương thức Register, Boot, Requires và Providers.
package config
