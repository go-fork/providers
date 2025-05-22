// Package config cung cấp giải pháp quản lý cấu hình (configuration) hiện đại, linh hoạt
// và mở rộng cho ứng dụng Go.
//
// Tính năng nổi bật:
//   - Hỗ trợ đa nguồn cấu hình: Load từ file YAML, JSON, biến môi trường (environment variables),
//     hoặc custom provider.
//   - Truy xuất cấu hình theo dot notation: Cho phép truy vấn key phân cấp dạng "a.b.c",
//     tự động gom các key con khi truy xuất key cha.
//   - Type-safe accessors: Lấy giá trị cấu hình với kiểu dữ liệu mong muốn (string, int, bool, map, slice)
//     và hỗ trợ giá trị mặc định.
//   - Unmarshal cấu hình vào struct: Dễ dàng ánh xạ cấu hình phẳng sang struct Go,
//     kể cả với key cha gom các key con.
//   - Thread-safe: Đảm bảo an toàn khi truy xuất và cập nhật cấu hình đồng thời.
//   - Extensible: Cho phép mở rộng với custom Formatter hoặc ServiceProvider.
//   - Dễ dàng tích hợp với Dependency Injection (DI) container.
//
// Kiến trúc và cách hoạt động:
//   - Tất cả các giá trị cấu hình được lưu trữ phẳng với key dạng dot notation.
//   - Interface Manager định nghĩa các phương thức chính để tương tác với cấu hình.
//   - DefaultManager là implementation mặc định của Manager, thread-safe và hỗ trợ đầy đủ các tính năng.
//   - Khi truy xuất key cha (ví dụ: "database"), package sẽ tự động gom các key con
//     ("database.host", "database.port", ...) thành map hoặc unmarshal vào struct.
//   - ServiceProvider cung cấp tích hợp với DI container, tự động nạp cấu hình từ biến môi trường và thư mục configs.
//   - Module formatter quản lý việc nạp cấu hình từ các nguồn khác nhau thông qua interface Formatter.
//   - Module utils cung cấp các tiện ích xử lý map phẳng và lồng nhau.
//
// Cấu trúc package:
//   - manager.go: Định nghĩa interface Manager và DefaultManager (triển khai mặc định).
//   - provider.go: ServiceProvider tích hợp với DI container.
//   - formatter/: Chứa các Formatter cho từng nguồn cấu hình (YAML, JSON, ENV).
//   - utils/: Các hàm tiện ích xử lý map phẳng và lồng nhau.
//
// Ví dụ sử dụng cơ bản:
//
// / // Khởi tạo manager
//
//	manager := config.NewManager()
//
//	// Nạp cấu hình từ file và biến môi trường
//	manager.Load(formatter.NewYamlFormatter("config.yaml"))
//	manager.Load(formatter.NewEnvFormatter("APP_"))
//
//	// Truy xuất giá trị đơn giản
//	appName := manager.GetString("app.name", "Default App")
//	port := manager.GetInt("app.port", 8080)
//	debug := manager.GetBool("app.debug", false)
//
//	// Truy xuất key cha gom các key con
//	db := manager.GetStringMap("database")
//	fmt.Println(db["host"], db["port"])
//
//	// Unmarshal vào struct
//	type DBConfig struct {
//	    Host string `json:"host"`
//	    Port int    `json:"port"`
//	}
//	var dbConf DBConfig
//	manager.Unmarshal("database", &dbConf)
//
// Tích hợp với DI container:
//
// / // Đăng ký ServiceProvider vào app
//
//	app.Register(config.NewServiceProvider())
//
//	// Lấy config manager từ container
//	manager := app.Make("config").(config.Manager)
//
// Gói này giúp centralize, validate và truy xuất configuration hiệu quả, an toàn, phù hợp cho mọi ứng dụng Go hiện đại.
package config
