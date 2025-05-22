# Go-Fork Cache Provider

Gói `cache` cung cấp hệ thống quản lý cache hiện đại, linh hoạt và mở rộng cho ứng dụng Go. Được thiết kế theo kiến trúc đa driver, gói này hỗ trợ nhiều loại lưu trữ cache khác nhau, từ bộ nhớ trong RAM đến các giải pháp phân tán như Redis và MongoDB.

## Giới thiệu

Caching là một phần thiết yếu để tối ưu hiệu suất trong các ứng dụng hiện đại. Gói `cache` cung cấp một cách thống nhất để lưu trữ và truy xuất dữ liệu cache từ nhiều nguồn khác nhau. Với thiết kế thread-safe và các tính năng nâng cao như remember pattern, batch operations, gói này là giải pháp hoàn chỉnh cho nhu cầu caching trong ứng dụng Go.

## Tính năng nổi bật

- **Đa dạng driver**: Hỗ trợ bộ nhớ (Memory), tệp tin (File), Redis, MongoDB và khả năng mở rộng driver tùy chỉnh.
- **TTL (Time To Live)**: Tự động quản lý thời gian sống cho các mục trong cache.
- **Remember pattern**: Hỗ trợ tính toán lười biếng và lưu trữ kết quả trong cache.
- **Batch operations**: Thao tác hàng loạt để tối ưu hiệu suất.
- **Serialization**: Tự động chuyển đổi giữa cấu trúc dữ liệu Go và định dạng lưu trữ.
- **Thread-safe**: An toàn khi truy xuất và cập nhật đồng thời.
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container.
- **Extensible**: Dễ dàng mở rộng với driver tùy chỉnh thông qua interface Driver.

## Cấu trúc package

```
cache/
  ├── doc.go                 # Tài liệu tổng quan về package
  ├── manager.go             # Định nghĩa interface Manager và DefaultManager
  ├── provider.go            # ServiceProvider tích hợp với DI container
  └── driver/
      ├── driver.go          # Định nghĩa interface Driver
      ├── memory.go          # Driver lưu cache trong bộ nhớ
      ├── file.go            # Driver lưu cache trong hệ thống file
      ├── redis.go           # Driver sử dụng Redis (v9+)
      └── mongodb.go         # Driver sử dụng MongoDB
```

## Cách hoạt động

### Đăng ký Service Provider

Service Provider cho phép tích hợp dễ dàng gói `cache` vào ứng dụng sử dụng DI container:

```go
// Trong file bootstrap của ứng dụng
import "github.com/go-fork/providers/cache"

func bootstrap(app interface{}) {
    // Đăng ký cache provider
    cacheProvider := cache.NewServiceProvider()
    cacheProvider.Register(app)
    
    // Boot các providers sau khi tất cả đã đăng ký
    cacheProvider.Boot(app)
}
```

ServiceProvider sẽ tự động:
1. Tạo một cache manager mới
2. Cấu hình driver mặc định (thường là memory)
3. Đăng ký manager vào container với key "cache"

### Sử dụng trực tiếp

Bạn có thể tạo và sử dụng cache manager mà không cần thông qua DI container:

```go
// Tạo manager mới
manager := cache.NewManager()

// Thêm memory driver
memoryDriver := driver.NewMemoryDriver()
manager.AddDriver("memory", memoryDriver)
manager.SetDefaultDriver("memory")

// Thêm file driver
fileDriver, err := driver.NewFileDriver("/path/to/cache/dir")
if err == nil {
    manager.AddDriver("file", fileDriver)
}

// Bắt đầu sử dụng cache
manager.Set("user:1", userData, 1*time.Hour)  // Cache với TTL 1 giờ
userData, exists := manager.Get("user:1")
manager.Delete("user:1")
```

### Làm việc với các Driver

#### Memory Driver

Driver này lưu cache trong bộ nhớ RAM, phù hợp cho ứng dụng đơn tiến trình:

```go
// Tạo memory driver
memDriver := driver.NewMemoryDriver()

// Thiết lập thời gian sống mặc định (1 giờ)
memDriver.SetDefaultExpiration(1 * time.Hour)
```

#### File Driver

Driver này lưu cache trong các file trên hệ thống tệp tin:

```go
// Tạo file driver với thư mục cache
fileDriver, err := driver.NewFileDriver("/path/to/cache")
if err != nil {
    // Xử lý lỗi
}

// Thiết lập thời gian sống mặc định
fileDriver.SetDefaultExpiration(30 * time.Minute)
```

#### Redis Driver

Driver này sử dụng Redis làm backend lưu trữ:

```go
// Tạo Redis client
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Tạo Redis driver
redisDriver := driver.NewRedisDriver(redisClient)
redisDriver.SetPrefix("myapp:")  // Tiền tố cho tất cả các key
```

#### MongoDB Driver

Driver này sử dụng MongoDB làm backend lưu trữ:

```go
// Tạo MongoDB client
ctx := context.Background()
clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
mongoClient, err := mongo.Connect(ctx, clientOptions)
if err != nil {
    // Xử lý lỗi
}

// Tạo MongoDB driver
mongoDriver, err := driver.NewMongoDBDriver(mongoClient, "cache_db", "cache_collection")
if err != nil {
    // Xử lý lỗi
}
```

### Sử dụng Remember Pattern

Remember pattern giúp tối ưu việc tính toán giá trị đắt đỏ chỉ khi cần:

```go
userData, err := manager.Remember("user:1", 1*time.Hour, func() (interface{}, error) {
    // Hàm này chỉ được gọi khi "user:1" không tồn tại trong cache
    return fetchUserFromDatabase(1)
})
```

### Batch Operations

Thao tác hàng loạt để giảm số lần giao tiếp với hệ thống cache:

```go
// Lấy nhiều giá trị trong một lần gọi
userIDs := []string{"user:1", "user:2", "user:3"}
foundUsers, missingKeys := manager.GetMultiple(userIDs)

// Lưu nhiều giá trị trong một lần gọi
users := map[string]interface{}{
    "user:4": user4Data,
    "user:5": user5Data,
}
manager.SetMultiple(users, 1*time.Hour)

// Xóa nhiều giá trị trong một lần gọi
manager.DeleteMultiple([]string{"user:1", "user:2"})
```

### Đóng các driver đúng cách

Luôn đóng manager khi không sử dụng để giải phóng tài nguyên:

```go
manager := cache.NewManager()
// ... cấu hình và sử dụng manager ...

// Đóng manager khi kết thúc
defer manager.Close()
```

### Tạo Custom Driver

Bạn có thể triển khai driver của riêng mình bằng cách tuân thủ interface Driver:

```go
type MyCustomDriver struct {
    // Các thành phần nội bộ của driver
}

// Triển khai các phương thức của interface Driver
func (d *MyCustomDriver) Get(ctx context.Context, key string) (interface{}, bool) {
    // Triển khai lấy giá trị
}

func (d *MyCustomDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    // Triển khai lưu trữ giá trị
}

// ... triển khai các phương thức còn lại ...
```

---

Để biết thêm thông tin chi tiết và API reference, vui lòng xem tài liệu trong file `doc.go` hoặc chạy lệnh `go doc github.com/go-fork/providers/cache`.
