// Package mocks cung cấp các implement giả lập (mock) cho các interface trong package config,
// được thiết kế để hỗ trợ viết unit test cho các ứng dụng sử dụng package config.
//
// # Đối tượng chính
//
//   - MockManager: Triển khai giả lập của interface Manager từ package config, cho phép
//     kiểm soát hoàn toàn hành vi và trạng thái của đối tượng quản lý cấu hình trong bối cảnh test.
//
// # Tính năng
//
//   - Giả lập tất cả các phương thức của interface Manager
//   - Khả năng điều khiển kết quả trả về của các phương thức
//   - Lưu trữ cấu hình trong bộ nhớ giúp kiểm tra dễ dàng
//   - Không phụ thuộc vào file cấu hình thực tế, phù hợp với unit test
//
// # Ví dụ sử dụng
//
//	// Tạo mock manager
//	mockCfg := mocks.NewMockManager()
//
//	// Thiết lập dữ liệu giả lập
//	mockCfg.Set("app.name", "TestApp")
//	mockCfg.Set("app.version", "1.0.0")
//	mockCfg.Set("database.port", 5432)
//
//	// Sử dụng trong test
//	name, ok := mockCfg.GetString("app.name")
//	assert.True(t, ok)
//	assert.Equal(t, "TestApp", name)
//
//	// Kiểm tra không tồn tại
//	_, ok = mockCfg.GetString("nonexistent.key")
//	assert.False(t, ok)
//
// Package mocks giúp đơn giản hóa quá trình viết unit test cho các ứng dụng
// sử dụng package config mà không cần phụ thuộc vào cấu hình thực tế.
package mocks
