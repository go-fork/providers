// Package cache cung cấp một hệ thống cache linh hoạt và có thể mở rộng cho ứng dụng.
//
// Package này sử dụng thiết kế dựa trên driver, cho phép lưu trữ và truy xuất dữ liệu cached
// từ nhiều nguồn khác nhau (memory, file, redis, mongodb, v.v.) với khả năng tùy chỉnh cao.
// Hỗ trợ các thao tác cơ bản như get, set, has, delete và cả auto-expiration.
//
// Các tính năng chính:
//   - Lưu trữ dữ liệu với TTL (Time To Live)
//   - Hỗ trợ nhiều driver cùng một lúc
//   - Khả năng mở rộng với driver tùy chỉnh
//   - Quản lý trung tâm thông qua Manager
//   - Caching trong bộ nhớ, file, Redis, MongoDB
//   - Hỗ trợ remember pattern để tính toán lười biếng
//   - Hỗ trợ thao tác hàng loạt (batch operations)
//
// Cấu trúc package:
//   - driver/: Định nghĩa interface và các driver thực thi
//   - driver.go: Interface chính và các định nghĩa cần thiết
//   - memory.go: Driver lưu cache trong bộ nhớ
//   - file.go: Driver lưu cache trong hệ thống file
//   - redis.go: Driver sử dụng Redis (v9+)
//   - mongodb.go: Driver sử dụng MongoDB
//   - manager.go: Manager để quản lý các driver cache
//   - provider.go: ServiceProvider để tích hợp với container DI
//
// Ví dụ sử dụng cơ bản:
//
//	// Khởi tạo manager
//	manager := cache.NewManager()
//
//	// Thêm driver
//	memoryDriver := driver.NewMemoryDriver()
//	manager.AddDriver("memory", memoryDriver)
//
//	// Lưu trữ và truy xuất dữ liệu
//	manager.Set("user:1", userData, 3600*time.Second) // TTL 1 giờ
//	userData, exists := manager.Get("user:1")
//	manager.Delete("user:1")
//
// Ví dụ sử dụng remember pattern:
//
//	userData, err := manager.Remember("user:1", 1*time.Hour, func() (interface{}, error) {
//		// Hàm này chỉ được gọi khi key không tồn tại trong cache
//		return fetchUserDataFromDatabase(1)
//	})
//
// Package cache giúp tối ưu hóa hiệu suất ứng dụng bằng cách giảm thiểu
// tải trọng truy vấn và cải thiện thời gian phản hồi.
package cache
