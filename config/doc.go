// Package config cung cấp giải pháp quản lý configuration linh hoạt, hiện đại cho ứng dụng Go.
//
// Gói này hỗ trợ load và truy xuất configuration từ nhiều nguồn (sources) khác nhau như file YAML, JSON, environment variables, hoặc custom provider.
// Tất cả các giá trị configuration được quản lý tập trung thông qua Manager, hỗ trợ truy vấn theo dot notation (ví dụ: "database.host", "app.name").
//
// Key Features:
//   - Multi-source configuration: Hỗ trợ load từ file (YAML, JSON), env, hoặc custom provider.
//   - Dot notation access: Truy xuất giá trị cấu hình dạng phân cấp bằng key chuỗi ("a.b.c").
//   - Default value: Hỗ trợ giá trị mặc định khi key không tồn tại.
//   - Type-safe accessors: Lấy giá trị với kiểu dữ liệu mong muốn (string, int, bool, map, slice).
//   - DI integration: Dễ dàng tích hợp với Dependency Injection container.
//   - Extensible: Cho phép mở rộng với custom Formatter hoặc ServiceProvider.
//
// Cấu trúc package:
//   - manager.go: Định nghĩa Manager interface và implementation mặc định.
//   - provider.go: ServiceProvider tích hợp với DI container.
//   - formatter/: Chứa các Formatter cho từng nguồn cấu hình (YAML, JSON, ENV).
//
// Example usage:
//
//	manager := config.NewManager()
//	manager.Load(formatter.NewYamlFormatter("config.yaml"))
//	manager.Load(formatter.NewEnvFormatter("APP_"))
//	appName := manager.GetString("app.name", "Default App")
//	dbHost := manager.GetString("database.host", "localhost")
//
// Gói này là một phần quan trọng trong infrastructure của ứng dụng Go hiện đại, giúp centralize, validate và truy xuất configuration một cách hiệu quả, an toàn.
package config
