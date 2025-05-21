package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisDriver cài đặt cache driver sử dụng Redis.
//
// RedisDriver lưu trữ dữ liệu cache trong Redis, một hệ thống lưu trữ key-value
// có tính năng persistence và phân tán. Driver này phù hợp cho các ứng dụng yêu cầu
// khả năng mở rộng, phân tán cache giữa nhiều instance ứng dụng và khả năng phục hồi
// sau khi khởi động lại. Nó cũng tận dụng các tính năng của Redis như key expiration.
type RedisDriver struct {
	client            *redis.Client                     // Redis client để giao tiếp với Redis server
	prefix            string                            // Tiền tố cho các key cache để tránh xung đột
	defaultExpiration time.Duration                     // Thời gian sống mặc định cho các entry không chỉ định TTL
	serializer        func(interface{}) ([]byte, error) // Hàm serialization để chuyển đổi giá trị thành dạng binary
	deserializer      func([]byte, interface{}) error   // Hàm deserialization để chuyển đổi từ binary
	hits              int64                             // Số lần cache hit
	misses            int64                             // Số lần cache miss
}

// RedisConfig cấu hình cho Redis driver.
//
// Cấu trúc này cung cấp các tùy chọn cấu hình chi tiết cho Redis driver
// như thông tin kết nối, prefix cho key và thời gian sống mặc định.
type RedisConfig struct {
	Host              string        // Hostname hoặc IP của Redis server
	Port              int           // Port của Redis server
	Password          string        // Mật khẩu xác thực (nếu có)
	DB                int           // Database index để sử dụng
	Prefix            string        // Tiền tố cho các key cache
	DefaultExpiration time.Duration // Thời gian sống mặc định
}

// NewRedisDriver tạo một Redis driver mới với cấu hình mặc định.
//
// Phương thức này khởi tạo một RedisDriver mới với thông tin kết nối cơ bản.
// Nó sử dụng cấu hình mặc định cho prefix là "cache:" và defaultExpiration là 5 phút.
//
// Params:
//   - host: Hostname hoặc IP của Redis server
//   - port: Port của Redis server
//   - password: Mật khẩu xác thực (empty string nếu không cần)
//   - db: Database index để sử dụng
//
// Returns:
//   - *RedisDriver: Driver đã được khởi tạo
//   - error: Lỗi nếu không thể kết nối tới Redis server
func NewRedisDriver(host string, port int, password string, db int) (*RedisDriver, error) {
	config := RedisConfig{
		Host:              host,
		Port:              port,
		Password:          password,
		DB:                db,
		Prefix:            "cache:",
		DefaultExpiration: 5 * time.Minute,
	}
	return NewRedisDriverWithConfig(config)
}

// NewRedisDriverWithConfig tạo một Redis driver mới với cấu hình chi tiết.
//
// Phương thức này khởi tạo một RedisDriver mới với cấu hình được cung cấp đầy đủ.
//
// Params:
//   - config: Cấu trúc chứa toàn bộ thông tin cấu hình cho driver
//
// Returns:
//   - *RedisDriver: Driver đã được khởi tạo
//   - error: Lỗi nếu không thể kết nối tới Redis server
func NewRedisDriverWithConfig(config RedisConfig) (*RedisDriver, error) {
	// Tạo Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Kiểm tra kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	// Khởi tạo driver
	driver := &RedisDriver{
		client:            client,
		prefix:            config.Prefix,
		defaultExpiration: config.DefaultExpiration,
		serializer:        json.Marshal,
		deserializer:      json.Unmarshal,
	}

	return driver, nil
}

// prefixKey thêm prefix vào key.
//
// Phương thức này thêm tiền tố đã cấu hình vào cache key để tạo thành Redis key hoàn chỉnh.
// Việc này giúp tránh xung đột với các key khác trong cùng một Redis database.
//
// Params:
//   - key: Cache key cần thêm tiền tố
//
// Returns:
//   - string: Redis key đã được thêm tiền tố
func (d *RedisDriver) prefixKey(key string) string {
	return d.prefix + key
}

// Get lấy một giá trị từ cache.
//
// Phương thức này lấy và giải mã giá trị từ Redis dựa trên key được cung cấp.
// Nếu key không tồn tại hoặc đã hết hạn, phương thức trả về false ở giá trị thứ hai
// và cập nhật bộ đếm miss. Nếu tìm thấy, phương thức trả về giá trị đã giải mã
// và cập nhật bộ đếm hit.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần lấy
//
// Returns:
//   - interface{}: Giá trị được lưu trong cache (nil nếu không tìm thấy)
//   - bool: true nếu tìm thấy key, false nếu ngược lại
func (d *RedisDriver) Get(ctx context.Context, key string) (interface{}, bool) {
	prefixedKey := d.prefixKey(key)

	// Lấy giá trị từ Redis
	data, err := d.client.Get(ctx, prefixedKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			// Key không tồn tại
			d.misses++
			return nil, false
		}
		// Lỗi khác
		return nil, false
	}

	// Giải mã dữ liệu
	var value interface{}
	if err := d.deserializer(data, &value); err != nil {
		return nil, false
	}

	d.hits++
	return value, true
}

// Set đặt một giá trị vào cache với TTL tùy chọn.
//
// Phương thức này mã hóa và lưu trữ một cặp key-value vào Redis với thời gian sống
// được chỉ định. Nếu key đã tồn tại, giá trị sẽ bị ghi đè.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key để lưu giá trị
//   - value: Giá trị cần lưu trữ
//   - ttl: Thời gian sống của giá trị (0 để sử dụng mặc định, -1 để không hết hạn)
//
// Returns:
//   - error: Lỗi nếu có trong quá trình mã hóa hoặc lưu trữ
func (d *RedisDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	prefixedKey := d.prefixKey(key)

	// Mã hóa dữ liệu
	data, err := d.serializer(value)
	if err != nil {
		return fmt.Errorf("could not serialize value: %w", err)
	}

	// Xác định thời gian hết hạn
	if ttl == 0 {
		ttl = d.defaultExpiration
	}

	// Lưu vào Redis
	return d.client.Set(ctx, prefixedKey, data, ttl).Err()
}

// Has kiểm tra xem một key có tồn tại trong cache không.
//
// Phương thức này kiểm tra sự tồn tại của key trong Redis bằng cách
// sử dụng lệnh EXISTS của Redis.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần kiểm tra
//
// Returns:
//   - bool: true nếu key tồn tại, false nếu ngược lại
func (d *RedisDriver) Has(ctx context.Context, key string) bool {
	prefixedKey := d.prefixKey(key)
	exists, _ := d.client.Exists(ctx, prefixedKey).Result()
	return exists > 0
}

// Delete xóa một key khỏi cache.
//
// Phương thức này xóa key và giá trị tương ứng khỏi Redis nếu tồn tại.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần xóa
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa
func (d *RedisDriver) Delete(ctx context.Context, key string) error {
	prefixedKey := d.prefixKey(key)
	return d.client.Del(ctx, prefixedKey).Err()
}

// Flush xóa tất cả các key khỏi cache có prefix đã định.
//
// Phương thức này quét và xóa tất cả các key có tiền tố đã cấu hình
// trong Redis database được sử dụng. Phương pháp này an toàn hơn so với
// FLUSHDB vì nó chỉ xóa các key thuộc về cache này.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//
// Returns:
//   - error: Lỗi nếu có trong quá trình quét hoặc xóa
func (d *RedisDriver) Flush(ctx context.Context) error {
	// Tìm tất cả các key có prefix
	pattern := d.prefix + "*"
	iter := d.client.Scan(ctx, 0, pattern, 0).Iterator()

	// Xóa từng key
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())

		// Xóa theo batch để tối ưu hiệu suất
		if len(keys) >= 100 {
			if err := d.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
			keys = []string{}
		}
	}

	// Xóa batch cuối cùng
	if len(keys) > 0 {
		return d.client.Del(ctx, keys...).Err()
	}

	return iter.Err()
}

// GetMultiple lấy nhiều giá trị từ cache
func (d *RedisDriver) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, []string) {
	results := make(map[string]interface{})
	missed := make([]string, 0)

	// Chuẩn bị các key với prefix
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = d.prefixKey(key)
	}

	// Sử dụng MGET để lấy nhiều giá trị cùng lúc
	values, err := d.client.MGet(ctx, prefixedKeys...).Result()
	if err != nil {
		// Lỗi, trả về tất cả keys là missed
		return results, keys
	}

	// Xử lý kết quả
	for i, value := range values {
		if value == nil {
			missed = append(missed, keys[i])
			continue
		}

		// Giải mã dữ liệu
		var decoded interface{}
		data, ok := value.(string)
		if !ok {
			missed = append(missed, keys[i])
			continue
		}

		if err := d.deserializer([]byte(data), &decoded); err != nil {
			missed = append(missed, keys[i])
			continue
		}

		results[keys[i]] = decoded
	}

	return results, missed
}

// SetMultiple đặt nhiều giá trị vào cache
func (d *RedisDriver) SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = d.defaultExpiration
	}

	// Sử dụng pipe để tối ưu hiệu suất
	pipe := d.client.Pipeline()

	for key, value := range values {
		// Mã hóa dữ liệu
		data, err := d.serializer(value)
		if err != nil {
			return fmt.Errorf("could not serialize value for key '%s': %w", key, err)
		}

		prefixedKey := d.prefixKey(key)
		pipe.Set(ctx, prefixedKey, data, ttl)
	}

	// Thực thi pipeline
	_, err := pipe.Exec(ctx)
	return err
}

// DeleteMultiple xóa nhiều key khỏi cache
func (d *RedisDriver) DeleteMultiple(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	// Chuẩn bị các key với prefix
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = d.prefixKey(key)
	}

	// Xóa tất cả các key cùng lúc
	return d.client.Del(ctx, prefixedKeys...).Err()
}

// Remember lấy một giá trị từ cache hoặc thực thi callback nếu không tìm thấy
func (d *RedisDriver) Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
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

// Stats trả về thông tin thống kê về cache
func (d *RedisDriver) Stats(ctx context.Context) map[string]interface{} {
	// Đếm số lượng key với prefix
	pattern := d.prefix + "*"
	count, err := d.client.Keys(ctx, pattern).Result()
	countVal := len(count)
	if err != nil {
		countVal = -1
	}

	// Lấy thông tin từ INFO command
	info, err := d.client.Info(ctx).Result()
	if err != nil {
		info = ""
	}

	return map[string]interface{}{
		"count":  countVal,
		"hits":   d.hits,
		"misses": d.misses,
		"type":   "redis",
		"prefix": d.prefix,
		"info":   info,
	}
}

// Close giải phóng tài nguyên của driver
func (d *RedisDriver) Close() error {
	return d.client.Close()
}

// WithSerializer thiết lập hàm serializer tùy chỉnh
func (d *RedisDriver) WithSerializer(serializer func(interface{}) ([]byte, error), deserializer func([]byte, interface{}) error) *RedisDriver {
	if serializer != nil && deserializer != nil {
		d.serializer = serializer
		d.deserializer = deserializer
	}
	return d
}
