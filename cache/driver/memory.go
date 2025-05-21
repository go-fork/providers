package driver

import (
	"context"
	"sync"
	"time"
)

// Item đại diện cho một mục trong memory cache.
//
// Cấu trúc này lưu trữ giá trị được cache và thời điểm hết hạn của nó.
// Nó được sử dụng nội bộ bởi MemoryDriver để quản lý các cache entry.
type Item struct {
	Value      interface{} // Giá trị được lưu trong cache
	Expiration int64       // Thời điểm hết hạn (UnixNano), 0 nếu không hết hạn
}

// Expired kiểm tra xem item đã hết hạn hay chưa.
//
// Phương thức này so sánh thời điểm hết hạn của item với thời gian hiện tại
// để xác định xem item đã hết hạn hay chưa.
//
// Returns:
//   - bool: true nếu item đã hết hạn, false nếu chưa hết hạn hoặc không có thời hạn
func (i *Item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}

// MemoryDriver cài đặt cache driver sử dụng memory (in-memory).
//
// MemoryDriver lưu trữ dữ liệu cache trong bộ nhớ chính của ứng dụng.
// Driver này cung cấp khả năng truy xuất dữ liệu nhanh nhất nhưng dữ liệu sẽ
// bị mất khi ứng dụng khởi động lại. Nó hỗ trợ TTL và tự động dọn dẹp
// các entry đã hết hạn.
type MemoryDriver struct {
	items             map[string]Item // Map lưu trữ các cache item
	mu                sync.RWMutex    // Mutex cho các thao tác thread-safe
	janitorInterval   time.Duration   // Khoảng thời gian giữa các lần dọn dẹp
	stopJanitor       chan bool       // Channel để dừng goroutine dọn dẹp
	janitorRunning    bool            // Flag đánh dấu goroutine dọn dẹp đang chạy
	defaultExpiration time.Duration   // Thời gian sống mặc định cho các entry không chỉ định TTL
	hits              int64           // Số lần cache hit
	misses            int64           // Số lần cache miss
}

// NewMemoryDriver tạo một memory driver mới với các tùy chọn mặc định.
//
// Phương thức này khởi tạo một MemoryDriver mới với các giá trị mặc định cho
// defaultExpiration (5 phút) và cleanupInterval (10 phút).
//
// Returns:
//   - *MemoryDriver: Driver đã được khởi tạo
func NewMemoryDriver() *MemoryDriver {
	return NewMemoryDriverWithOptions(5*time.Minute, 10*time.Minute)
}

// NewMemoryDriverWithOptions tạo một memory driver mới với các tùy chọn chi tiết.
//
// Phương thức này khởi tạo một MemoryDriver mới với các tùy chọn thời gian hết hạn mặc định
// và khoảng thời gian dọn dẹp được chỉ định.
//
// Params:
//   - defaultExpiration: Thời gian sống mặc định cho các cache entry không chỉ định TTL
//   - cleanupInterval: Khoảng thời gian giữa các lần dọn dẹp tự động
//
// Returns:
//   - *MemoryDriver: Driver đã được khởi tạo
func NewMemoryDriverWithOptions(defaultExpiration, cleanupInterval time.Duration) *MemoryDriver {
	driver := &MemoryDriver{
		items:             make(map[string]Item),
		janitorInterval:   cleanupInterval,
		defaultExpiration: defaultExpiration,
		stopJanitor:       make(chan bool),
	}

	// Chỉ chạy janitor nếu có khoảng thời gian dọn dẹp > 0
	if cleanupInterval > 0 {
		go driver.startJanitor()
		driver.janitorRunning = true
	}

	return driver
}

// Get lấy một giá trị từ cache.
//
// Phương thức này tìm kiếm và trả về giá trị từ cache dựa trên key được cung cấp.
// Nếu key không tồn tại hoặc đã hết hạn, phương thức trả về false ở giá trị thứ hai
// và cập nhật bộ đếm miss. Nếu tìm thấy và còn hạn, phương thức trả về giá trị và
// cập nhật bộ đếm hit.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần tìm
//
// Returns:
//   - interface{}: Giá trị được lưu trong cache (nil nếu không tìm thấy)
//   - bool: true nếu tìm thấy key và chưa hết hạn, false nếu ngược lại
func (d *MemoryDriver) Get(ctx context.Context, key string) (interface{}, bool) {
	d.mu.RLock()
	item, found := d.items[key]
	d.mu.RUnlock()

	if !found {
		d.mu.Lock()
		d.misses++
		d.mu.Unlock()
		return nil, false
	}

	if item.Expired() {
		d.mu.Lock()
		d.misses++
		delete(d.items, key)
		d.mu.Unlock()
		return nil, false
	}

	d.mu.Lock()
	d.hits++
	d.mu.Unlock()
	return item.Value, true
}

// Set đặt một giá trị vào cache với TTL tùy chọn.
//
// Phương thức này lưu trữ một cặp key-value vào cache với thời gian sống
// được chỉ định. Nếu key đã tồn tại, giá trị sẽ bị ghi đè.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key để lưu giá trị
//   - value: Giá trị cần lưu trữ
//   - ttl: Thời gian sống của giá trị (0 để sử dụng mặc định, -1 để không hết hạn)
//
// Returns:
//   - error: Luôn trả về nil trong memory driver
func (d *MemoryDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var exp int64

	if ttl == 0 {
		if d.defaultExpiration > 0 {
			exp = time.Now().Add(d.defaultExpiration).UnixNano()
		}
	} else if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.items[key] = Item{
		Value:      value,
		Expiration: exp,
	}
	return nil
}

// Has kiểm tra xem một key có tồn tại trong cache không.
//
// Phương thức này xác định liệu một key có tồn tại trong cache và chưa hết hạn hay không.
// Nó sử dụng phương thức Get để thực hiện kiểm tra và do đó cũng cập nhật bộ đếm hit/miss.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần kiểm tra
//
// Returns:
//   - bool: true nếu key tồn tại và chưa hết hạn, false nếu ngược lại
func (d *MemoryDriver) Has(ctx context.Context, key string) bool {
	_, exists := d.Get(ctx, key)
	return exists
}

// Delete xóa một key khỏi cache.
//
// Phương thức này xóa key và giá trị tương ứng khỏi cache nếu tồn tại.
// Nếu key không tồn tại, thao tác này không có tác dụng và không trả về lỗi.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần xóa
//
// Returns:
//   - error: Luôn trả về nil trong memory driver
func (d *MemoryDriver) Delete(ctx context.Context, key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.items, key)
	return nil
}

// Flush xóa tất cả các key khỏi cache.
//
// Phương thức này xóa tất cả dữ liệu trong cache, làm trống hoàn toàn bộ nhớ cache.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//
// Returns:
//   - error: Luôn trả về nil trong memory driver
func (d *MemoryDriver) Flush(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.items = make(map[string]Item)
	return nil
}

// GetMultiple lấy nhiều giá trị từ cache.
//
// Phương thức này lấy các giá trị tương ứng với nhiều key trong một lần gọi.
// Nó thực hiện gọi Get cho từng key và tổng hợp kết quả.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - keys: Danh sách các khóa cần lấy
//
// Returns:
//   - map[string]interface{}: Map chứa các key tìm thấy và giá trị tương ứng
//   - []string: Danh sách các key không tìm thấy hoặc đã hết hạn
func (d *MemoryDriver) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, []string) {
	results := make(map[string]interface{})
	missed := make([]string, 0)

	for _, key := range keys {
		value, found := d.Get(ctx, key)
		if found {
			results[key] = value
		} else {
			missed = append(missed, key)
		}
	}

	return results, missed
}

// SetMultiple đặt nhiều giá trị vào cache.
//
// Phương thức này lưu trữ nhiều cặp key-value vào cache trong một lần gọi
// với cùng một thời gian sống. Nó thực hiện gọi Set cho từng cặp key-value.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - values: Map chứa các key và giá trị tương ứng cần lưu trữ
//   - ttl: Thời gian sống chung cho tất cả các giá trị
//
// Returns:
//   - error: Lỗi nếu có trong quá trình thực hiện (luôn là nil trong memory driver)
func (d *MemoryDriver) SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
	for key, value := range values {
		if err := d.Set(ctx, key, value, ttl); err != nil {
			return err
		}
	}
	return nil
}

// DeleteMultiple xóa nhiều key khỏi cache.
//
// Phương thức này xóa nhiều key và giá trị tương ứng khỏi cache trong một lần gọi.
// Nó thực hiện gọi Delete cho từng key.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - keys: Danh sách các khóa cần xóa
//
// Returns:
//   - error: Lỗi nếu có trong quá trình thực hiện (luôn là nil trong memory driver)
func (d *MemoryDriver) DeleteMultiple(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := d.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

// Remember lấy một giá trị từ cache hoặc thực thi callback nếu không tìm thấy.
//
// Phương thức này kiểm tra xem một key có tồn tại trong cache không, nếu có thì
// trả về giá trị tương ứng. Nếu key không tồn tại hoặc đã hết hạn, phương thức
// sẽ gọi hàm callback để lấy dữ liệu, lưu kết quả vào cache và trả về giá trị đó.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần tìm hoặc lưu vào cache
//   - ttl: Thời gian sống của giá trị nếu phải lấy từ callback
//   - callback: Hàm được gọi để lấy dữ liệu khi key không có trong cache
//
// Returns:
//   - interface{}: Giá trị từ cache hoặc từ callback
//   - error: Lỗi nếu có trong quá trình thực hiện hoặc từ callback
func (d *MemoryDriver) Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// Kiểm tra cache trước
	value, found := d.Get(ctx, key)
	if found {
		return value, nil
	}

	// Không tìm thấy, gọi callback
	value, err := callback()
	if err != nil {
		return nil, err
	}

	// Lưu kết quả vào cache
	err = d.Set(ctx, key, value, ttl)
	return value, err
}

// Stats trả về thông tin thống kê về cache.
//
// Phương thức này thu thập và trả về các thông tin thống kê về trạng thái
// hiện tại của memory cache như số lượng item, số lần hit/miss, v.v.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//
// Returns:
//   - map[string]interface{}: Map chứa các thông tin thống kê
func (d *MemoryDriver) Stats(ctx context.Context) map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()

	itemCount := len(d.items)
	stats := map[string]interface{}{
		"count":  itemCount,
		"hits":   d.hits,
		"misses": d.misses,
		"type":   "memory",
	}

	return stats
}

// Close giải phóng tài nguyên của driver.
//
// Phương thức này dừng goroutine janitor nếu đang chạy và giải phóng
// các tài nguyên khác được sử dụng bởi driver.
//
// Returns:
//   - error: Lỗi nếu có trong quá trình giải phóng tài nguyên (luôn là nil trong memory driver)
func (d *MemoryDriver) Close() error {
	if d.janitorRunning {
		d.stopJanitor <- true
	}
	return nil
}

// startJanitor bắt đầu một routine định kỳ dọn dẹp các mục đã hết hạn.
//
// Phương thức này chạy một goroutine định kỳ gọi deleteExpired để
// xóa các cache entry đã hết hạn theo khoảng thời gian đã cấu hình.
func (d *MemoryDriver) startJanitor() {
	ticker := time.NewTicker(d.janitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.deleteExpired()
		case <-d.stopJanitor:
			return
		}
	}
}

// deleteExpired xóa tất cả các mục đã hết hạn.
//
// Phương thức này quét qua tất cả các item trong cache map,
// kiểm tra thời gian hết hạn và xóa những item đã quá hạn.
func (d *MemoryDriver) deleteExpired() {
	now := time.Now().UnixNano()

	d.mu.Lock()
	defer d.mu.Unlock()

	for k, v := range d.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(d.items, k)
		}
	}
}
