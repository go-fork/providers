package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.fork.vn/providers/cache/driver"
)

// Manager là interface chính để quản lý các driver cache.
//
// Manager cung cấp một lớp trừu tượng để làm việc với nhiều cache driver khác nhau.
// Nó cho phép đăng ký nhiều driver dựa trên tên và đặt một driver mặc định
// để sử dụng cho các thao tác cache. Nó cũng cung cấp các phương thức tiện ích
// cho tất cả các thao tác cơ bản trên cache mà không cần trực tiếp tương tác với driver.
type Manager interface {
	// Get lấy một giá trị từ cache.
	//
	// Phương thức này tìm kiếm và trả về giá trị từ cache mặc định dựa trên key được cung cấp.
	//
	// Params:
	//   - key: Cache key cần tìm
	//
	// Returns:
	//   - interface{}: Giá trị được lưu trong cache (nil nếu không tìm thấy)
	//   - bool: true nếu tìm thấy key và chưa hết hạn, false nếu ngược lại
	Get(key string) (interface{}, bool)

	// Set đặt một giá trị vào cache với TTL tùy chọn.
	//
	// Phương thức này lưu trữ một cặp key-value vào cache mặc định với thời gian sống
	// được chỉ định. Nếu key đã tồn tại, giá trị sẽ bị ghi đè.
	//
	// Params:
	//   - key: Cache key để lưu giá trị
	//   - value: Giá trị cần lưu trữ
	//   - ttl: Thời gian sống của giá trị (0 để sử dụng mặc định của driver, -1 để không hết hạn)
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình lưu trữ hoặc driver mặc định không được cấu hình
	Set(key string, value interface{}, ttl time.Duration) error

	// Has kiểm tra xem một key có tồn tại trong cache không.
	//
	// Phương thức này xác định liệu một key có tồn tại trong cache mặc định và chưa hết hạn hay không.
	//
	// Params:
	//   - key: Cache key cần kiểm tra
	//
	// Returns:
	//   - bool: true nếu key tồn tại và chưa hết hạn, false nếu ngược lại
	Has(key string) bool

	// Delete xóa một key khỏi cache.
	//
	// Phương thức này xóa key và giá trị tương ứng khỏi cache mặc định nếu tồn tại.
	//
	// Params:
	//   - key: Cache key cần xóa
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình xóa hoặc driver mặc định không được cấu hình
	Delete(key string) error

	// Flush xóa tất cả các key khỏi cache.
	//
	// Phương thức này xóa tất cả dữ liệu trong cache mặc định, làm trống hoàn toàn bộ nhớ cache.
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình xóa hoặc driver mặc định không được cấu hình
	Flush() error

	// GetMultiple lấy nhiều giá trị từ cache.
	//
	// Phương thức này lấy các giá trị tương ứng với nhiều key từ cache mặc định trong một lần gọi.
	//
	// Params:
	//   - keys: Danh sách các khóa cần lấy
	//
	// Returns:
	//   - map[string]interface{}: Map chứa các key tìm thấy và giá trị tương ứng
	//   - []string: Danh sách các key không tìm thấy hoặc đã hết hạn
	GetMultiple(keys []string) (map[string]interface{}, []string)

	// SetMultiple đặt nhiều giá trị vào cache.
	//
	// Phương thức này lưu trữ nhiều cặp key-value vào cache mặc định trong một lần gọi
	// với cùng một thời gian sống.
	//
	// Params:
	//   - values: Map chứa các key và giá trị tương ứng cần lưu trữ
	//   - ttl: Thời gian sống chung cho tất cả các giá trị
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình lưu trữ hoặc driver mặc định không được cấu hình
	SetMultiple(values map[string]interface{}, ttl time.Duration) error

	// DeleteMultiple xóa nhiều key khỏi cache.
	//
	// Phương thức này xóa nhiều key và giá trị tương ứng khỏi cache mặc định trong một lần gọi.
	//
	// Params:
	//   - keys: Danh sách các khóa cần xóa
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình xóa hoặc driver mặc định không được cấu hình
	DeleteMultiple(keys []string) error

	// Remember lấy một giá trị từ cache hoặc thực thi callback nếu không tìm thấy.
	//
	// Phương thức này kiểm tra xem một key có tồn tại trong cache mặc định không, nếu có thì
	// trả về giá trị tương ứng. Nếu key không tồn tại hoặc đã hết hạn, phương thức
	// sẽ gọi hàm callback để lấy dữ liệu, lưu kết quả vào cache và trả về giá trị đó.
	//
	// Params:
	//   - key: Cache key cần tìm hoặc lưu vào cache
	//   - ttl: Thời gian sống của giá trị nếu phải lấy từ callback
	//   - callback: Hàm được gọi để lấy dữ liệu khi key không có trong cache
	//
	// Returns:
	//   - interface{}: Giá trị từ cache hoặc từ callback
	//   - error: Lỗi nếu có trong quá trình thực hiện, từ callback, hoặc driver mặc định không được cấu hình
	Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error)

	// AddDriver thêm một driver vào manager.
	//
	// Phương thức này đăng ký một driver mới với manager theo tên xác định.
	// Nếu chưa có driver mặc định được đặt, driver đầu tiên được thêm vào sẽ trở thành mặc định.
	//
	// Params:
	//   - name: Tên định danh cho driver
	//   - driver: Đối tượng driver cần thêm vào
	AddDriver(name string, driver driver.Driver)

	// SetDefaultDriver đặt driver mặc định.
	//
	// Phương thức này thiết lập driver mặc định được sử dụng cho các thao tác cache.
	// Driver phải đã được đăng ký với manager trước đó.
	//
	// Params:
	//   - name: Tên của driver cần đặt làm mặc định
	SetDefaultDriver(name string)

	// Driver trả về một driver cụ thể theo tên.
	//
	// Phương thức này lấy một driver đã đăng ký dựa theo tên của nó.
	//
	// Params:
	//   - name: Tên của driver cần lấy
	//
	// Returns:
	//   - driver.Driver: Đối tượng driver được yêu cầu
	//   - error: Lỗi nếu driver không tồn tại
	Driver(name string) (driver.Driver, error)

	// Stats trả về thông tin thống kê về tất cả các driver.
	//
	// Phương thức này thu thập và trả về các thông tin thống kê về trạng thái hiện tại
	// của tất cả các driver đã đăng ký.
	//
	// Returns:
	//   - map[string]map[string]interface{}: Map chứa thông tin thống kê của từng driver, với key là tên driver
	Stats() map[string]map[string]interface{}

	// Close đóng tất cả các driver.
	//
	// Phương thức này giải phóng tài nguyên của tất cả các driver đã đăng ký.
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình đóng bất kỳ driver nào
	Close() error
}

// manager là implementation mặc định của Manager.
//
// manager quản lý nhiều driver cache thông qua một map driver và cung cấp
// cơ chế để thực hiện các thao tác cache qua driver mặc định. Nó đảm bảo thread-safety
// thông qua RWMutex và cung cấp các phương thức tiện ích để tương tác với nhiều driver.
type manager struct {
	drivers       map[string]driver.Driver // Map chứa tất cả các driver đã đăng ký
	defaultDriver string                   // Tên của driver mặc định
	mu            sync.RWMutex             // Mutex cho các thao tác thread-safe
}

// NewManager tạo một manager mới.
//
// Phương thức này khởi tạo một DefaultManager mới, sẵn sàng để đăng ký các driver.
//
// Returns:
//   - Manager: Đối tượng Manager mới được khởi tạo
func NewManager() Manager {
	return &manager{
		drivers: make(map[string]driver.Driver),
	}
}

// Get lấy một giá trị từ cache.
//
// Phương thức này tìm kiếm và trả về giá trị từ cache mặc định dựa trên key được cung cấp.
//
// Params:
//   - key: Cache key cần tìm
//
// Returns:
//   - interface{}: Giá trị được lưu trong cache (nil nếu không tìm thấy)
//   - bool: true nếu tìm thấy key và chưa hết hạn, false nếu ngược lại
func (m *manager) Get(key string) (interface{}, bool) {
	driver, err := m.DefaultDriver()
	if err != nil {
		return nil, false
	}
	return driver.Get(context.Background(), key)
}

// Set đặt một giá trị vào cache với TTL tùy chọn.
//
// Phương thức này lưu trữ một cặp key-value vào cache mặc định với thời gian sống được chỉ định.
//
// Params:
//   - key: Cache key để lưu giá trị
//   - value: Giá trị cần lưu trữ
//   - ttl: Thời gian sống của giá trị
//
// Returns:
//   - error: Lỗi nếu có trong quá trình lưu trữ hoặc driver mặc định không được cấu hình
func (m *manager) Set(key string, value interface{}, ttl time.Duration) error {
	driver, err := m.DefaultDriver()
	if err != nil {
		return err
	}
	return driver.Set(context.Background(), key, value, ttl)
}

// Has kiểm tra xem một key có tồn tại trong cache không.
//
// Phương thức này xác định liệu một key có tồn tại trong cache mặc định và chưa hết hạn hay không.
//
// Params:
//   - key: Cache key cần kiểm tra
//
// Returns:
//   - bool: true nếu key tồn tại và chưa hết hạn, false nếu ngược lại
func (m *manager) Has(key string) bool {
	driver, err := m.DefaultDriver()
	if err != nil {
		return false
	}
	return driver.Has(context.Background(), key)
}

// Delete xóa một key khỏi cache.
//
// Phương thức này xóa key và giá trị tương ứng khỏi cache mặc định nếu tồn tại.
//
// Params:
//   - key: Cache key cần xóa
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa hoặc driver mặc định không được cấu hình
func (m *manager) Delete(key string) error {
	driver, err := m.DefaultDriver()
	if err != nil {
		return err
	}
	return driver.Delete(context.Background(), key)
}

// Flush xóa tất cả các key khỏi cache.
//
// Phương thức này xóa tất cả dữ liệu trong cache mặc định, làm trống hoàn toàn bộ nhớ cache.
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa hoặc driver mặc định không được cấu hình
func (m *manager) Flush() error {
	driver, err := m.DefaultDriver()
	if err != nil {
		return err
	}
	return driver.Flush(context.Background())
}

// GetMultiple lấy nhiều giá trị từ cache.
//
// Phương thức này lấy các giá trị tương ứng với nhiều key từ cache mặc định trong một lần gọi.
//
// Params:
//   - keys: Danh sách các khóa cần lấy
//
// Returns:
//   - map[string]interface{}: Map chứa các key tìm thấy và giá trị tương ứng
//   - []string: Danh sách các key không tìm thấy hoặc đã hết hạn
func (m *manager) GetMultiple(keys []string) (map[string]interface{}, []string) {
	driver, err := m.DefaultDriver()
	if err != nil {
		return make(map[string]interface{}), keys
	}
	return driver.GetMultiple(context.Background(), keys)
}

// SetMultiple đặt nhiều giá trị vào cache.
//
// Phương thức này lưu trữ nhiều cặp key-value vào cache mặc định trong một lần gọi
// với cùng một thời gian sống.
//
// Params:
//   - values: Map chứa các key và giá trị tương ứng cần lưu trữ
//   - ttl: Thời gian sống chung cho tất cả các giá trị
//
// Returns:
//   - error: Lỗi nếu có trong quá trình lưu trữ hoặc driver mặc định không được cấu hình
func (m *manager) SetMultiple(values map[string]interface{}, ttl time.Duration) error {
	driver, err := m.DefaultDriver()
	if err != nil {
		return err
	}
	return driver.SetMultiple(context.Background(), values, ttl)
}

// DeleteMultiple xóa nhiều key khỏi cache.
//
// Phương thức này xóa nhiều key và giá trị tương ứng khỏi cache mặc định trong một lần gọi.
//
// Params:
//   - keys: Danh sách các khóa cần xóa
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa hoặc driver mặc định không được cấu hình
func (m *manager) DeleteMultiple(keys []string) error {
	driver, err := m.DefaultDriver()
	if err != nil {
		return err
	}
	return driver.DeleteMultiple(context.Background(), keys)
}

// Remember lấy một giá trị từ cache hoặc thực thi callback nếu không tìm thấy.
//
// Phương thức này kiểm tra xem một key có tồn tại trong cache mặc định không, nếu có thì
// trả về giá trị tương ứng. Nếu key không tồn tại hoặc đã hết hạn, phương thức
// sẽ gọi hàm callback để lấy dữ liệu, lưu kết quả vào cache và trả về giá trị đó.
//
// Params:
//   - key: Cache key cần tìm hoặc lưu vào cache
//   - ttl: Thời gian sống của giá trị nếu phải lấy từ callback
//   - callback: Hàm được gọi để lấy dữ liệu khi key không có trong cache
//
// Returns:
//   - interface{}: Giá trị từ cache hoặc từ callback
//   - error: Lỗi nếu có trong quá trình thực hiện, từ callback, hoặc driver mặc định không được cấu hình
func (m *manager) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	driver, err := m.DefaultDriver()
	if err != nil {
		return nil, err
	}
	return driver.Remember(context.Background(), key, ttl, callback)
}

// AddDriver thêm một driver vào manager.
//
// Phương thức này đăng ký một driver mới với manager theo tên xác định.
// Nếu chưa có driver mặc định được đặt, driver đầu tiên được thêm vào sẽ trở thành mặc định.
//
// Params:
//   - name: Tên định danh cho driver
//   - driver: Đối tượng driver cần thêm vào
func (m *manager) AddDriver(name string, driver driver.Driver) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.drivers[name] = driver

	// Đặt driver đầu tiên được thêm làm mặc định nếu chưa có driver mặc định
	if m.defaultDriver == "" {
		m.defaultDriver = name
	}
}

// SetDefaultDriver đặt driver mặc định.
//
// Phương thức này thiết lập driver mặc định được sử dụng cho các thao tác cache.
// Driver phải đã được đăng ký với manager trước đó.
//
// Params:
//   - name: Tên của driver cần đặt làm mặc định
func (m *manager) SetDefaultDriver(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.drivers[name]; ok {
		m.defaultDriver = name
	}
}

// Driver trả về một driver cụ thể theo tên.
//
// Phương thức này lấy một driver đã đăng ký dựa theo tên của nó.
//
// Params:
//   - name: Tên của driver cần lấy
//
// Returns:
//   - driver.Driver: Đối tượng driver được yêu cầu
//   - error: Lỗi nếu driver không tồn tại
func (m *manager) Driver(name string) (driver.Driver, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if driver, ok := m.drivers[name]; ok {
		return driver, nil
	}

	return nil, fmt.Errorf("cache driver '%s' not found", name)
}

// Stats trả về thông tin thống kê về tất cả các driver.
//
// Phương thức này thu thập và trả về các thông tin thống kê về trạng thái hiện tại
// của tất cả các driver đã đăng ký.
//
// Returns:
//   - map[string]map[string]interface{}: Map chứa thông tin thống kê của từng driver, với key là tên driver
func (m *manager) Stats() map[string]map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]map[string]interface{})
	for name, driver := range m.drivers {
		stats[name] = driver.Stats(context.Background())
	}
	return stats
}

// Close đóng tất cả các driver.
//
// Phương thức này giải phóng tài nguyên của tất cả các driver đã đăng ký.
//
// Returns:
//   - error: Lỗi nếu có trong quá trình đóng bất kỳ driver nào
func (m *manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var firstErr error
	for name, driver := range m.drivers {
		if err := driver.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed to close cache driver '%s': %w", name, err)
		}
	}
	return firstErr
}

// DefaultDriver trả về driver mặc định.
//
// Phương thức này lấy driver mặc định hiện tại từ manager, kiểm tra tính hợp lệ
// và trả về lỗi nếu không có driver mặc định hoặc driver đã bị xóa.
//
// Returns:
//   - driver.Driver: Đối tượng driver mặc định
//   - error: Lỗi nếu không có driver mặc định hoặc driver mặc định không tồn tại
func (m *manager) DefaultDriver() (driver.Driver, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.defaultDriver == "" {
		return nil, fmt.Errorf("no default cache driver set")
	}

	if driver, ok := m.drivers[m.defaultDriver]; ok {
		return driver, nil
	}

	return nil, fmt.Errorf("default cache driver '%s' not found", m.defaultDriver)
}
