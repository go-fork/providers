package driver

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-fork/providers/cache/config"
	redisManager "github.com/go-fork/providers/redis"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

type RedisDriver interface {
	Driver
	WithSerializer(serializer string) RedisDriver
}

// redisDriver cài đặt cache driver sử dụng Redis.
//
// redisDriver lưu trữ dữ liệu cache trong Redis, một hệ thống lưu trữ key-value
// có tính năng persistence và phân tán. Driver này phù hợp cho các ứng dụng yêu cầu
// khả năng mở rộng, phân tán cache giữa nhiều instance ứng dụng và khả năng phục hồi
// sau khi khởi động lại. Nó cũng tận dụng các tính năng của Redis như key expiration.
type redisDriver struct {
	client       *redis.Client                     // Redis client để giao tiếp với Redis server
	prefix       string                            // Tiền tố cho các key cache để tránh xung đột
	default_ttl  time.Duration                     // Thời gian sống mặc định cho các entry không chỉ định TTL
	serializer   func(interface{}) ([]byte, error) // Hàm serialization để chuyển đổi giá trị thành dạng binary
	deserializer func([]byte, interface{}) error   // Hàm deserialization để chuyển đổi từ binary
	hits         int64                             // Số lần cache hit
	misses       int64                             // Số lần cache miss
}

// NewRedisDriver tạo một Redis driver mới với cấu hình mặc định.
//
// Phương thức này khởi tạo một RedisDriver mới với thông tin kết nối cơ bản.
// Nó sử dụng cấu hình mặc định cho prefix là "cache:" và default_ttl là 5 phút.
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
func NewRedisDriver(config config.DriverRedisConfig, redis_manager redisManager.Manager) (RedisDriver, error) {
	if !config.Enabled {
		return nil, fmt.Errorf("redis driver is not enabled")
	}
	client, err := redis_manager.Client()
	if err != nil {
		return nil, fmt.Errorf("could not create Redis client: %w", err)
	}
	// Khởi tạo driver
	driver := &redisDriver{
		client:       client,
		prefix:       "cache:", // Tiền tố mặc định
		default_ttl:  time.Duration(config.DefaultTTL) * time.Second,
		serializer:   json.Marshal,
		deserializer: json.Unmarshal,
		hits:         0,
		misses:       0,
	}
	switch config.Serializer {

	case "gob":
		driver.serializer = func(v interface{}) ([]byte, error) {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(v); err != nil {
				return nil, fmt.Errorf("could not serialize value: %w", err)
			}
			return buf.Bytes(), nil
		}
		driver.deserializer = func(data []byte, v interface{}) error {
			buf := bytes.NewBuffer(data)
			dec := gob.NewDecoder(buf)
			if err := dec.Decode(v); err != nil {
				return fmt.Errorf("could not deserialize value: %w", err)
			}
			return nil
		}
	case "msgpack":
		driver.serializer = func(v interface{}) ([]byte, error) {
			data, err := msgpack.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("could not serialize value: %w", err)
			}
			return data, nil
		}
		driver.deserializer = func(data []byte, v interface{}) error {
			if err := msgpack.Unmarshal(data, v); err != nil {
				return fmt.Errorf("could not deserialize value: %w", err)
			}
			return nil
		}
	default:
		driver.serializer = json.Marshal
		driver.deserializer = json.Unmarshal
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
func (d *redisDriver) prefixKey(key string) string {
	return d.prefix + key
}

// Get lấy một giá trị từ cache.
func (d *redisDriver) Get(ctx context.Context, key string) (interface{}, bool) {
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
	// Giải mã dữ liệu - cần xử lý khác nhau tùy theo serializer
	var value interface{}

	// For GOB and MSGPACK, we need to decode differently
	if d.deserializer != nil {
		// Create a buffer to hold the decoded value
		var decodedValue interface{}
		if err := d.deserializer(data, &decodedValue); err != nil {
			d.misses++
			return nil, false
		}
		value = decodedValue
	} else {
		// Fallback to JSON
		if err := json.Unmarshal(data, &value); err != nil {
			d.misses++
			return nil, false
		}
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
func (d *redisDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	prefixedKey := d.prefixKey(key)

	// Mã hóa dữ liệu
	data, err := d.serializer(value)
	if err != nil {
		return fmt.Errorf("could not serialize value: %w", err)
	}

	// Xác định thời gian hết hạn
	if ttl == 0 {
		ttl = d.default_ttl
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
func (d *redisDriver) Has(ctx context.Context, key string) bool {
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
func (d *redisDriver) Delete(ctx context.Context, key string) error {
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
func (d *redisDriver) Flush(ctx context.Context) error {
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
func (d *redisDriver) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, []string) {
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

		// Use the same deserialization logic as Get method
		if d.deserializer != nil {
			if err := d.deserializer([]byte(data), &decoded); err != nil {
				missed = append(missed, keys[i])
				continue
			}
		} else {
			if err := json.Unmarshal([]byte(data), &decoded); err != nil {
				missed = append(missed, keys[i])
				continue
			}
		}

		results[keys[i]] = decoded
	}

	return results, missed
}

// SetMultiple đặt nhiều giá trị vào cache
func (d *redisDriver) SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = d.default_ttl
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
func (d *redisDriver) DeleteMultiple(ctx context.Context, keys []string) error {
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
func (d *redisDriver) Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
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
func (d *redisDriver) Stats(ctx context.Context) map[string]interface{} {
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
func (d *redisDriver) Close() error {
	return d.client.Close()
}

// WithSerializer thiết lập serializer theo tên
func (d *redisDriver) WithSerializer(serializerName string) RedisDriver {
	newDriver := &redisDriver{
		client:      d.client,
		prefix:      d.prefix,
		default_ttl: d.default_ttl,
		hits:        d.hits,
		misses:      d.misses,
	}

	switch serializerName {
	case "gob":
		newDriver.serializer = func(v interface{}) ([]byte, error) {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(v); err != nil {
				return nil, fmt.Errorf("could not serialize value: %w", err)
			}
			return buf.Bytes(), nil
		}
		newDriver.deserializer = func(data []byte, v interface{}) error {
			buf := bytes.NewBuffer(data)
			dec := gob.NewDecoder(buf)
			if err := dec.Decode(v); err != nil {
				return fmt.Errorf("could not deserialize value: %w", err)
			}
			return nil
		}
	case "msgpack":
		newDriver.serializer = func(v interface{}) ([]byte, error) {
			data, err := msgpack.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("could not serialize value: %w", err)
			}
			return data, nil
		}
		newDriver.deserializer = func(data []byte, v interface{}) error {
			if err := msgpack.Unmarshal(data, v); err != nil {
				return fmt.Errorf("could not deserialize value: %w", err)
			}
			return nil
		}
	default: // json
		newDriver.serializer = json.Marshal
		newDriver.deserializer = json.Unmarshal
	}

	return newDriver
}
