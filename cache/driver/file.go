package driver

import (
	"context"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileDriver cài đặt cache driver sử dụng file system.
//
// FileDriver lưu trữ dữ liệu cache dưới dạng các file trên hệ thống file,
// với mỗi cache entry tương ứng với một file riêng biệt. Driver này thích hợp
// cho các ứng dụng cần persistence và có thể phục hồi dữ liệu cache sau khi khởi động lại.
// Nó cũng hỗ trợ TTL (Time To Live) và tự động dọn dẹp các entry đã hết hạn.
type FileDriver struct {
	directory         string        // Đường dẫn thư mục lưu trữ cache
	defaultExpiration time.Duration // Thời gian sống mặc định cho các entry không chỉ định TTL
	mu                sync.RWMutex  // Mutex cho các thao tác thread-safe
	janitorInterval   time.Duration // Khoảng thời gian giữa các lần dọn dẹp
	stopJanitor       chan bool     // Channel để dừng goroutine dọn dẹp
	janitorRunning    bool          // Flag đánh dấu goroutine dọn dẹp đang chạy
	hits              int64         // Số lần cache hit
	misses            int64         // Số lần cache miss
}

// FileCache là cấu trúc lưu trữ dữ liệu trong file.
//
// Cấu trúc này được dùng để serialize và deserialize dữ liệu cache
// khi lưu trữ và đọc từ file. Nó chứa giá trị cần cache và thời gian hết hạn.
type FileCache struct {
	Value      interface{} // Giá trị được lưu trong cache
	Expiration int64       // Thời điểm hết hạn (UnixNano), 0 nếu không hết hạn
}

// NewFileDriver tạo một file driver mới với các tùy chọn mặc định.
//
// Phương thức này khởi tạo một FileDriver mới với thư mục lưu trữ được chỉ định
// và các giá trị mặc định cho defaultExpiration (5 phút) và cleanupInterval (10 phút).
//
// Params:
//   - directory: Đường dẫn thư mục để lưu trữ các file cache
//
// Returns:
//   - *FileDriver: Driver đã được khởi tạo
//   - error: Lỗi nếu không thể tạo thư mục cache
func NewFileDriver(directory string) (*FileDriver, error) {
	return NewFileDriverWithOptions(directory, 5*time.Minute, 10*time.Minute)
}

// NewFileDriverWithOptions tạo một file driver mới với các tùy chọn chi tiết.
//
// Phương thức này khởi tạo một FileDriver mới với thư mục lưu trữ và các tùy chọn
// thời gian hết hạn mặc định và khoảng thời gian dọn dẹp được chỉ định.
//
// Params:
//   - directory: Đường dẫn thư mục để lưu trữ các file cache
//   - defaultExpiration: Thời gian sống mặc định cho các cache entry không chỉ định TTL
//   - cleanupInterval: Khoảng thời gian giữa các lần dọn dẹp tự động
//
// Returns:
//   - *FileDriver: Driver đã được khởi tạo
//   - error: Lỗi nếu không thể tạo thư mục cache
func NewFileDriverWithOptions(directory string, defaultExpiration, cleanupInterval time.Duration) (*FileDriver, error) {
	// Tạo thư mục nếu không tồn tại
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("unable to create cache directory: %w", err)
	}

	driver := &FileDriver{
		directory:         directory,
		defaultExpiration: defaultExpiration,
		janitorInterval:   cleanupInterval,
		stopJanitor:       make(chan bool),
	}

	// Chỉ chạy janitor nếu có khoảng thời gian dọn dẹp > 0
	if cleanupInterval > 0 {
		go driver.startJanitor()
		driver.janitorRunning = true
	}

	return driver, nil
}

// keyToFilename chuyển đổi key thành tên file an toàn.
//
// Phương thức này chuyển đổi một cache key thành đường dẫn file đầy đủ trong thư mục cache.
// Trong thực tế, nên sử dụng một hàm hash như md5 hoặc sha1 để tránh các vấn đề với ký tự đặc biệt
// và đảm bảo tên file hợp lệ trên hệ thống file.
//
// Params:
//   - key: Cache key cần chuyển đổi
//
// Returns:
//   - string: Đường dẫn đầy đủ đến file cache
func (d *FileDriver) keyToFilename(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("invalid key: key is empty")
	}
	for _, c := range key {
		if c == '/' || c == '\\' || c == 0 {
			return "", fmt.Errorf("invalid key: contains forbidden character")
		}
	}
	h := sha1.New()
	_, err := h.Write([]byte(key))
	if err != nil {
		return "", fmt.Errorf("invalid key: %w", err)
	}
	hash := hex.EncodeToString(h.Sum(nil))
	return filepath.Join(d.directory, hash), nil
}

// Get lấy một giá trị từ cache.
//
// Phương thức này đọc và giải mã dữ liệu từ file cache tương ứng với key được chỉ định.
// Nếu file không tồn tại hoặc đã hết hạn, phương thức sẽ trả về false và cập nhật
// bộ đếm miss. Nếu tìm thấy và còn hạn, phương thức trả về giá trị và cập nhật bộ đếm hit.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần lấy
//
// Returns:
//   - interface{}: Giá trị được lưu trong cache (nil nếu không tìm thấy)
//   - bool: true nếu tìm thấy key và chưa hết hạn, false nếu ngược lại
func (d *FileDriver) Get(ctx context.Context, key string) (interface{}, bool) {
	filename, err := d.keyToFilename(key)
	if err != nil {
		d.mu.Lock()
		d.misses++
		d.mu.Unlock()
		return nil, false
	}

	// Kiểm tra xem file có tồn tại không
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		d.mu.Lock()
		d.misses++
		d.mu.Unlock()
		return nil, false
	}

	// Mở file
	file, err := os.Open(filename)
	if err != nil {
		d.mu.Lock()
		d.misses++
		d.mu.Unlock()
		return nil, false
	}
	defer file.Close()

	// Giải mã dữ liệu
	var cache FileCache
	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(&cache); err != nil {
		d.mu.Lock()
		d.misses++
		d.mu.Unlock()
		return nil, false
	}

	// Kiểm tra xem đã hết hạn chưa
	if cache.Expiration > 0 && time.Now().UnixNano() > cache.Expiration {
		d.mu.Lock()
		d.misses++
		d.mu.Unlock()
		os.Remove(filename) // Xóa file đã hết hạn
		return nil, false
	}

	d.mu.Lock()
	d.hits++
	d.mu.Unlock()
	return cache.Value, true
}

// Set đặt một giá trị vào cache với TTL tùy chọn.
//
// Phương thức này mã hóa và lưu trữ một cặp key-value vào một file trong thư mục cache.
// Thời gian sống (TTL) có thể được chỉ định, hoặc sử dụng giá trị mặc định nếu ttl = 0,
// hoặc không có thời hạn nếu ttl = -1.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key để lưu giá trị
//   - value: Giá trị cần lưu trữ
//   - ttl: Thời gian sống của giá trị (0 để sử dụng mặc định, -1 để không hết hạn)
//
// Returns:
//   - error: Lỗi nếu có trong quá trình tạo, mã hóa hoặc ghi file
func (d *FileDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	filename, err := d.keyToFilename(key)
	if err != nil {
		return err
	}
	var exp int64

	if ttl == 0 {
		if d.defaultExpiration > 0 {
			exp = time.Now().Add(d.defaultExpiration).UnixNano()
		}
	} else if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	// Tạo cấu trúc cache
	cache := FileCache{
		Value:      value,
		Expiration: exp,
	}

	// Mở file để ghi
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create cache file: %w", err)
	}
	defer file.Close()

	// Mã hóa và ghi vào file
	encoder := gob.NewEncoder(file)
	return encoder.Encode(cache)
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
func (d *FileDriver) Has(ctx context.Context, key string) bool {
	_, exists := d.Get(ctx, key)
	return exists
}

// Delete xóa một key khỏi cache.
//
// Phương thức này xóa file cache tương ứng với key được chỉ định.
// Nếu file không tồn tại, thao tác này không có tác dụng và không trả về lỗi.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần xóa
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa file
func (d *FileDriver) Delete(ctx context.Context, key string) error {
	filename, err := d.keyToFilename(key)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil // File không tồn tại, không cần xóa
	}
	return os.Remove(filename)
}

// Flush xóa tất cả các key khỏi cache.
//
// Phương thức này xóa tất cả các file trong thư mục cache, làm trống hoàn toàn bộ nhớ cache.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa files
func (d *FileDriver) Flush(ctx context.Context) error {
	dir, err := os.Open(d.directory)
	if err != nil {
		return err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	var errs []error
	for _, name := range names {
		err = os.Remove(filepath.Join(d.directory, name))
		if err != nil {
			errs = append(errs, fmt.Errorf("file '%s': %w", name, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("Flush errors: %v", errs)
	}
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
func (d *FileDriver) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, []string) {
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
//   - error: Lỗi nếu có trong quá trình thực hiện
func (d *FileDriver) SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
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
//   - error: Lỗi nếu có trong quá trình thực hiện
func (d *FileDriver) DeleteMultiple(ctx context.Context, keys []string) error {
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
func (d *FileDriver) Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	value, found := d.Get(ctx, key)
	if found {
		return value, nil
	}

	value, err := callback()
	if err != nil {
		return nil, err
	}

	err = d.Set(ctx, key, value, ttl)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Stats trả về thông tin thống kê về cache.
//
// Phương thức này thu thập và trả về các thông tin thống kê về trạng thái
// hiện tại của file cache như số lượng file, tổng dung lượng, số lần hit/miss, v.v.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//
// Returns:
//   - map[string]interface{}: Map chứa các thông tin thống kê
func (d *FileDriver) Stats(ctx context.Context) map[string]interface{} {
	var itemCount int
	var size int64

	// Đếm số lượng file và kích thước
	filepath.Walk(d.directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			itemCount++
			size += info.Size()
		}
		return nil
	})

	d.mu.RLock()
	defer d.mu.RUnlock()

	return map[string]interface{}{
		"count":  itemCount,
		"size":   size,
		"hits":   d.hits,
		"misses": d.misses,
		"type":   "file",
		"path":   d.directory,
	}
}

// Close giải phóng tài nguyên của driver.
//
// Phương thức này dừng goroutine janitor nếu đang chạy và giải phóng
// các tài nguyên khác được sử dụng bởi driver.
//
// Returns:
//   - error: Lỗi nếu có trong quá trình giải phóng tài nguyên
func (d *FileDriver) Close() error {
	if d.janitorRunning {
		d.stopJanitor <- true
	}
	return nil
}

// startJanitor bắt đầu một routine định kỳ dọn dẹp các file đã hết hạn.
//
// Phương thức này chạy một goroutine định kỳ gọi deleteExpired để
// xóa các file cache đã hết hạn theo khoảng thời gian đã cấu hình.
func (d *FileDriver) startJanitor() {
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

// deleteExpired xóa tất cả các file đã hết hạn.
//
// Phương thức này quét qua tất cả các file trong thư mục cache,
// đọc thông tin thời gian hết hạn và xóa file nếu đã quá hạn.
func (d *FileDriver) deleteExpired() {
	now := time.Now().UnixNano()
	dir, err := os.Open(d.directory)
	if err != nil {
		return
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}

	for _, name := range names {
		filename := filepath.Join(d.directory, name)
		file, err := os.Open(filename)
		if err != nil {
			continue
		}

		var cache FileCache
		decoder := gob.NewDecoder(file)
		if err = decoder.Decode(&cache); err != nil {
			file.Close()
			continue
		}
		file.Close()

		if cache.Expiration > 0 && now > cache.Expiration {
			_ = os.Remove(filename) // Ignore error, continue
		}
	}
}
