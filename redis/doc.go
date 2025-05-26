// Package redis cung cấp một trình quản lý kết nối Redis và service provider cho các ứng dụng Go.
//
// # Tổng quan
//
// Package này giúp quản lý các kết nối Redis và cung cấp chúng thông qua dependency injection.
// Nó triển khai cả Redis client tiêu chuẩn và universal client có thể được sử dụng cho
// triển khai Redis Cluster, Sentinel và standalone.
//
// # Thành phần chính
//
//   - Manager: Interface chính để quản lý kết nối Redis, cung cấp các phương thức để tạo và quản lý
//     các kết nối, bao gồm standard client và universal client.
//
//   - ServiceProvider: Tích hợp với DI container để đăng ký các dịch vụ Redis.
//
//   - Config: Cấu hình cho Redis client và universal client, hỗ trợ tất cả các tùy chọn kết nối.
//
// # Sử dụng cơ bản
//
//	// Khởi tạo Redis manager với cấu hình mặc định
//	manager := redis.NewManager()
//
//	// Lấy Redis client
//	client, err := manager.Client()
//	if err != nil {
//	    // Xử lý lỗi
//	}
//
//	// Sử dụng client để thực hiện các thao tác Redis
//	err = client.Set(ctx, "key", "value", 0).Err()
//
// # Tích hợp với DI Container
//
//	// Tạo service provider
//	provider := redis.NewServiceProvider()
//
//	// Đăng ký provider với ứng dụng
//	provider.Register(app)
//
// # Phụ thuộc
//
// Package này phụ thuộc vào:
//   - github.com/go-fork/providers/config cho quản lý cấu hình
//   - github.com/go-fork/di cho dependency injection
//   - github.com/redis/go-redis/v9 cho Redis client
//
// # Khả năng tùy chỉnh
//
// Cung cấp nhiều tùy chọn cấu hình cho Redis client như:
//   - Địa chỉ máy chủ, cổng kết nối
//   - Xác thực (mật khẩu)
//   - Timeout và các tham số kết nối
//   - Cấu hình cluster và sentinel
//   - Tiền tố khóa
package redis
