// Package mocks provides mock implementations for cache drivers.
package mocks

import (
	"context"
	"sync"
	"time"
)

// MockDriver là một implementation giả của interface driver.Driver dùng cho testing.
//
// MockDriver cung cấp một triển khai giả lập đầy đủ của interface driver.Driver với
// các hàm có thể tùy chỉnh được. Nó được thiết kế để sử dụng trong kiểm thử đơn vị
// để mô phỏng các tình huống khác nhau mà không cần phụ thuộc vào cơ sở hạ tầng thực.
type MockDriver struct {
	data               map[string]interface{} // Map lưu trữ dữ liệu cache nội bộ
	mutex              sync.RWMutex           // Mutex để đảm bảo thread-safety
	GetFunc            func(ctx context.Context, key string) (interface{}, bool)
	SetFunc            func(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	HasFunc            func(ctx context.Context, key string) bool
	DeleteFunc         func(ctx context.Context, key string) error
	FlushFunc          func(ctx context.Context) error
	GetMultipleFunc    func(ctx context.Context, keys []string) (map[string]interface{}, []string)
	SetMultipleFunc    func(ctx context.Context, values map[string]interface{}, ttl time.Duration) error
	DeleteMultipleFunc func(ctx context.Context, keys []string) error
	RememberFunc       func(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error)
	StatsFunc          func(ctx context.Context) map[string]interface{}
	CloseFunc          func() error
}

// NewMockDriver tạo một instance mới của MockDriver với các triển khai mặc định.
//
// Phương thức này khởi tạo một MockDriver mới với các triển khai mặc định cho mọi phương thức,
// sử dụng map nội bộ để lưu trữ dữ liệu. Mỗi phương thức mặc định có thể được ghi đè
// để điều chỉnh hành vi trong các test cases cụ thể.
//
// Returns:
//   - *MockDriver: Đối tượng MockDriver mới đã được cấu hình với các hành vi mặc định
func NewMockDriver() *MockDriver {
	driver := &MockDriver{
		data: make(map[string]interface{}),
	}

	// Default implementation for Get
	driver.GetFunc = func(ctx context.Context, key string) (interface{}, bool) {
		driver.mutex.RLock()
		defer driver.mutex.RUnlock()
		val, ok := driver.data[key]
		return val, ok
	}

	// Default implementation for Set
	driver.SetFunc = func(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
		driver.mutex.Lock()
		defer driver.mutex.Unlock()
		driver.data[key] = value
		return nil
	}

	// Default implementation for Has
	driver.HasFunc = func(ctx context.Context, key string) bool {
		driver.mutex.RLock()
		defer driver.mutex.RUnlock()
		_, ok := driver.data[key]
		return ok
	}

	// Default implementation for Delete
	driver.DeleteFunc = func(ctx context.Context, key string) error {
		driver.mutex.Lock()
		defer driver.mutex.Unlock()
		delete(driver.data, key)
		return nil
	}

	// Default implementation for Flush
	driver.FlushFunc = func(ctx context.Context) error {
		driver.mutex.Lock()
		defer driver.mutex.Unlock()
		driver.data = make(map[string]interface{})
		return nil
	}

	// Default implementation for GetMultiple
	driver.GetMultipleFunc = func(ctx context.Context, keys []string) (map[string]interface{}, []string) {
		driver.mutex.RLock()
		defer driver.mutex.RUnlock()
		result := make(map[string]interface{})
		var missingKeys []string
		for _, key := range keys {
			if val, ok := driver.data[key]; ok {
				result[key] = val
			} else {
				missingKeys = append(missingKeys, key)
			}
		}
		return result, missingKeys
	}

	// Default implementation for SetMultiple
	driver.SetMultipleFunc = func(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
		driver.mutex.Lock()
		defer driver.mutex.Unlock()
		for k, v := range values {
			driver.data[k] = v
		}
		return nil
	}

	// Default implementation for DeleteMultiple
	driver.DeleteMultipleFunc = func(ctx context.Context, keys []string) error {
		driver.mutex.Lock()
		defer driver.mutex.Unlock()
		for _, key := range keys {
			delete(driver.data, key)
		}
		return nil
	}

	// Default implementation for Remember
	driver.RememberFunc = func(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
		driver.mutex.RLock()
		val, ok := driver.data[key]
		driver.mutex.RUnlock()
		if ok {
			return val, nil
		}

		// Call callback to generate value
		value, err := callback()
		if err != nil {
			return nil, err
		}

		// Store the value
		driver.mutex.Lock()
		driver.data[key] = value
		driver.mutex.Unlock()
		return value, nil
	}

	// Default implementation for Stats
	driver.StatsFunc = func(ctx context.Context) map[string]interface{} {
		driver.mutex.RLock()
		defer driver.mutex.RUnlock()
		return map[string]interface{}{
			"items": len(driver.data),
		}
	}

	// Default implementation for Close
	driver.CloseFunc = func() error {
		return nil
	}

	return driver
}

// Get lấy một giá trị từ mock cache.
//
// Phương thức này gọi GetFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó sẽ trả về giá trị được lưu trữ trong map data nội bộ nếu key tồn tại.
//
// Params:
//   - ctx: Context thực thi
//   - key: Cache key cần lấy
//
// Returns:
//   - interface{}: Giá trị được lưu trong cache
//   - bool: true nếu key tồn tại, false nếu ngược lại
func (d *MockDriver) Get(ctx context.Context, key string) (interface{}, bool) {
	return d.GetFunc(ctx, key)
}

// Set lưu trữ một giá trị vào mock cache.
//
// Phương thức này gọi SetFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó lưu trữ giá trị vào map data nội bộ. TTL được bỏ qua trong triển khai mặc định.
//
// Params:
//   - ctx: Context thực thi
//   - key: Cache key để lưu
//   - value: Giá trị cần lưu trữ
//   - ttl: Thời gian sống của key (được bỏ qua trong triển khai mặc định)
//
// Returns:
//   - error: Lỗi nếu có trong quá trình lưu trữ
func (d *MockDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return d.SetFunc(ctx, key, value, ttl)
}

// Has kiểm tra xem một key có tồn tại trong mock cache không.
//
// Phương thức này gọi HasFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó kiểm tra sự tồn tại của key trong map data nội bộ.
//
// Params:
//   - ctx: Context thực thi
//   - key: Cache key cần kiểm tra
//
// Returns:
//   - bool: true nếu key tồn tại, false nếu ngược lại
func (d *MockDriver) Has(ctx context.Context, key string) bool {
	return d.HasFunc(ctx, key)
}

// Delete xóa một giá trị khỏi mock cache.
//
// Phương thức này gọi DeleteFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó xóa key khỏi map data nội bộ.
//
// Params:
//   - ctx: Context thực thi
//   - key: Cache key cần xóa
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa
func (d *MockDriver) Delete(ctx context.Context, key string) error {
	return d.DeleteFunc(ctx, key)
}

// Flush xóa tất cả các giá trị khỏi mock cache.
//
// Phương thức này gọi FlushFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó tạo một map data mới, loại bỏ tất cả các giá trị hiện có.
//
// Params:
//   - ctx: Context thực thi
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa
func (d *MockDriver) Flush(ctx context.Context) error {
	return d.FlushFunc(ctx)
}

// GetMultiple lấy nhiều giá trị từ mock cache.
//
// Phương thức này gọi GetMultipleFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó lặp qua danh sách keys, thu thập các giá trị tìm thấy vào map kết quả, và các key không tìm thấy
// vào slice missingKeys.
//
// Params:
//   - ctx: Context thực thi
//   - keys: Danh sách các cache key cần lấy
//
// Returns:
//   - map[string]interface{}: Map chứa các key tìm thấy và giá trị tương ứng
//   - []string: Danh sách các key không tìm thấy
func (d *MockDriver) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, []string) {
	return d.GetMultipleFunc(ctx, keys)
}

// SetMultiple lưu trữ nhiều giá trị vào mock cache.
//
// Phương thức này gọi SetMultipleFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó lặp qua map values và lưu từng cặp key-value vào map data nội bộ. TTL được bỏ qua trong
// triển khai mặc định.
//
// Params:
//   - ctx: Context thực thi
//   - values: Map chứa các cặp key-value cần lưu trữ
//   - ttl: Thời gian sống chung cho tất cả các giá trị (được bỏ qua trong triển khai mặc định)
//
// Returns:
//   - error: Lỗi nếu có trong quá trình lưu trữ
func (d *MockDriver) SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
	return d.SetMultipleFunc(ctx, values, ttl)
}

// DeleteMultiple xóa nhiều giá trị khỏi mock cache.
//
// Phương thức này gọi DeleteMultipleFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó lặp qua danh sách keys và xóa từng key khỏi map data nội bộ.
//
// Params:
//   - ctx: Context thực thi
//   - keys: Danh sách các cache key cần xóa
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa
func (d *MockDriver) DeleteMultiple(ctx context.Context, keys []string) error {
	return d.DeleteMultipleFunc(ctx, keys)
}

// Remember lấy một giá trị từ cache hoặc thực thi callback để tạo nó.
//
// Phương thức này gọi RememberFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó kiểm tra xem key có tồn tại trong map data nội bộ không. Nếu có, trả về giá trị đó.
// Nếu không, nó gọi callback để tạo giá trị, lưu kết quả vào cache và trả về.
//
// Params:
//   - ctx: Context thực thi
//   - key: Cache key cần lấy hoặc tạo
//   - ttl: Thời gian sống của giá trị nếu phải tạo (được bỏ qua trong triển khai mặc định)
//   - callback: Hàm được gọi để tạo giá trị nếu key không tồn tại
//
// Returns:
//   - interface{}: Giá trị từ cache hoặc callback
//   - error: Lỗi nếu có trong quá trình thực hiện hoặc từ callback
func (d *MockDriver) Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	return d.RememberFunc(ctx, key, ttl, callback)
}

// Stats trả về thông tin thống kê về mock cache.
//
// Phương thức này gọi StatsFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó trả về một map chứa thông tin về số lượng items trong cache.
//
// Params:
//   - ctx: Context thực thi
//
// Returns:
//   - map[string]interface{}: Map chứa thông tin thống kê
func (d *MockDriver) Stats(ctx context.Context) map[string]interface{} {
	return d.StatsFunc(ctx)
}

// Close đóng kết nối mock driver.
//
// Phương thức này gọi CloseFunc được định nghĩa trong MockDriver. Trong triển khai mặc định,
// nó trả về nil để biểu thị thành công, vì mock driver không có tài nguyên thực sự cần giải phóng.
//
// Returns:
//   - error: Lỗi nếu có trong quá trình đóng kết nối
func (d *MockDriver) Close() error {
	return d.CloseFunc()
}
