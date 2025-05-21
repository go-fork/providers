// Package mocks provides mock implementations for cache drivers.
//
// Các mock objects trong package này được thiết kế để sử dụng trong testing,
// cung cấp cách triển khai giả của các drivers cho cache system.
//
// Ví dụ sử dụng MockDriver:
//
//	mockDriver := mocks.NewMockDriver()
//
//	// Tùy chỉnh hành vi cho Get
//	mockDriver.GetFunc = func(ctx context.Context, key string) (interface{}, bool) {
//		if key == "special_key" {
//			return "special_value", true
//		}
//		return nil, false
//	}
//
// MockDriver triển khai đầy đủ interface driver.Driver và có thể được sử dụng
// trong bất kỳ ngữ cảnh nào yêu cầu một Driver implementation.
package mocks
