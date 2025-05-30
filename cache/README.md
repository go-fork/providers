# Cache Provider

## Giới thiệu

Cache Provider là một package cung cấp hệ thống quản lý cache hiện đại, linh hoạt và có khả năng mở rộng cao cho framework dependency injection go-fork. Provider này cung cấp tích hợp cache với nhiều backend khác nhau như Memory, File, Redis và MongoDB trong ứng dụng Go. Package này được thiết kế để giúp đơn giản hóa việc tích hợp cache vào ứng dụng Go của bạn, đồng thời hỗ trợ các tính năng nâng cao như Remember pattern, batch operations và TTL management.

## Tổng quan

Cache Provider hỗ trợ:
- Tích hợp dễ dàng với framework dependency injection go-fork
- Đa dạng driver storage backend (Memory, File, Redis, MongoDB)
- TTL (Time To Live) tự động cho cache entries
- Remember pattern để lazy computation và caching
- Batch operations để tối ưu hiệu suất
- Thread-safe operations với concurrent access
- Giao diện đơn giản cho các thao tác cache phổ biến
- Monitoring và statistics cho performance tracking

## Cài đặt

```bash
go get go.fork.vn/providers/cache
```

## Cấu hình

Sao chép file cấu hình mẫu và chỉnh sửa theo nhu cầu:

```bash
cp configs/app.sample.yaml configs/app.yaml
```

### Ví dụ cấu hình

```yaml
cache:
  # Driver mặc định sẽ được sử dụng
  default_driver: "memory"
  
  # Cấu hình các drivers
  drivers:
    # Memory driver - cache trong RAM
    memory:
      type: "memory"
      default_ttl: 3600         # TTL mặc định (giây)
      cleanup_interval: 600     # Interval dọn dẹp expired entries (giây)
      max_items: 1000          # Số lượng item tối đa
    
    # File driver - cache trong file system
    file:
      type: "file"
      path: "/tmp/cache"        # Thư mục lưu cache
      default_ttl: 1800        # TTL mặc định (giây)
      file_permissions: "0644"  # Quyền file
    
    # Redis driver - cache trong Redis
    redis:
      type: "redis"
      addr: "localhost:6379"    # Địa chỉ Redis server
      password: ""              # Password (nếu có)
      db: 0                     # Database number
      prefix: "myapp:"          # Prefix cho tất cả keys
      default_ttl: 3600        # TTL mặc định (giây)
    
    # MongoDB driver - cache trong MongoDB
    mongodb:
      type: "mongodb"
      uri: "mongodb://localhost:27017"  # MongoDB connection URI
      database: "cache_db"              # Database name
      collection: "cache_collection"    # Collection name
      default_ttl: 3600                # TTL mặc định (giây)
```

## Sử dụng

### Thiết lập cơ bản

```go
package main

import (
    "context"
    "log"
    "time"
    
    "go.fork.vn/di"
    "go.fork.vn/providers/config"
    "go.fork.vn/providers/cache"
    "go.fork.vn/providers/cache/driver"
)

func main() {
    // Tạo DI container
    container := di.New()
    
    // Đăng ký provider config (nếu sử dụng service config)
    configProvider := config.NewServiceProvider()
    container.Register(configProvider)
    
    // Đăng ký Cache provider
    cacheProvider := cache.NewServiceProvider()
    container.Register(cacheProvider)
    
    // Boot các service providers
    container.Boot()
    
    // Lấy Cache manager sử dụng MustMake
    // MustMake sẽ panic nếu service không tồn tại hoặc không thể tạo được
    cacheManager := container.MustMake("cache").(cache.Manager)
    
    // Tạo context với timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Lưu dữ liệu vào cache
    err := cacheManager.Set("user:1", map[string]interface{}{
        "name":      "Nguyễn Văn A",
        "email":     "nguyen@example.com",
        "createdAt": time.Now(),
    }, 1*time.Hour)
    if err != nil {
        log.Fatal("Không thể lưu vào cache:", err)
    }
    
    // Lấy dữ liệu từ cache
    userData, exists := cacheManager.Get("user:1")
    if exists {
        log.Printf("Dữ liệu user từ cache: %+v\n", userData)
    }
    
    log.Println("Cache hoạt động thành công!")
}
```

### Sử dụng các phương thức của Manager

```go
// Khởi tạo context với timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Kiểm tra key có tồn tại không
exists := cacheManager.Has("user:1")
log.Printf("Key user:1 tồn tại: %v\n", exists)

// Remember pattern - lazy computation
userData, err := cacheManager.Remember("expensive_user:1", 1*time.Hour, func() (interface{}, error) {
    // Function này chỉ được gọi khi cache miss
    return fetchExpensiveUserData(1)
})
if err != nil {
    log.Fatal("Remember pattern thất bại:", err)
}

// Batch operations - lấy nhiều keys cùng lúc
userKeys := []string{"user:1", "user:2", "user:3"}
foundUsers, missingKeys := cacheManager.GetMultiple(userKeys)
log.Printf("Tìm thấy users: %v, thiếu keys: %v\n", len(foundUsers), missingKeys)

// Batch set - lưu nhiều values cùng lúc  
userDataMap := map[string]interface{}{
    "user:4": map[string]string{"name": "User 4"},
    "user:5": map[string]string{"name": "User 5"},
}
err = cacheManager.SetMultiple(userDataMap, 30*time.Minute)
if err != nil {
    log.Fatal("Batch set thất bại:", err)
}

// Xóa key
err = cacheManager.Delete("user:1")
if err != nil {
    log.Fatal("Không thể xóa key:", err)
}

// Xóa nhiều keys cùng lúc
err = cacheManager.DeleteMultiple([]string{"user:2", "user:3"})
if err != nil {
    log.Fatal("Batch delete thất bại:", err)
}

// Xóa toàn bộ cache
err = cacheManager.Clear()
if err != nil {
    log.Fatal("Không thể xóa toàn bộ cache:", err)
}
```

### Các Services đăng ký

Provider đăng ký các services sau trong DI container:

- `cache` - Instance Cache Manager
- `cache.manager` - Alias cho Cache Manager

Ví dụ truy xuất các services này với MustMake:

```go
// Lấy Cache Manager
cacheManager := container.MustMake("cache").(cache.Manager)

// Lấy Cache Manager thông qua alias
manager := container.MustMake("cache.manager").(cache.Manager)
```

## Danh sách phương thức

### Các phương thức cơ bản

| Phương thức | Mô tả |
|------------|-------|
| `Get(key string) (interface{}, bool)` | Lấy một giá trị từ cache theo key |
| `Set(key string, value interface{}, ttl time.Duration) error` | Đặt một giá trị vào cache với TTL |
| `Has(key string) bool` | Kiểm tra xem một key có tồn tại trong cache không |
| `Delete(key string) error` | Xóa một key khỏi cache |
| `Flush() error` | Xóa tất cả các key khỏi cache |

### Các phương thức batch operations

| Phương thức | Mô tả |
|------------|-------|
| `GetMultiple(keys []string) (map[string]interface{}, []string)` | Lấy nhiều giá trị từ cache |
| `SetMultiple(values map[string]interface{}, ttl time.Duration) error` | Đặt nhiều giá trị vào cache |
| `DeleteMultiple(keys []string) error` | Xóa nhiều key khỏi cache |

### Các phương thức nâng cao

| Phương thức | Mô tả |
|------------|-------|
| `Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error)` | Lấy giá trị từ cache hoặc thực thi callback nếu không tìm thấy |

### Các phương thức quản lý driver

| Phương thức | Mô tả |
|------------|-------|
| `AddDriver(name string, driver driver.Driver)` | Thêm một driver vào manager |
| `SetDefaultDriver(name string)` | Đặt driver mặc định |
| `Driver(name string) (driver.Driver, error)` | Trả về một driver cụ thể theo tên |

### Các phương thức tiện ích

| Phương thức | Mô tả |
|------------|-------|
| `Stats() map[string]map[string]interface{}` | Trả về thông tin thống kê về tất cả các driver |
| `Close() error` | Đóng tất cả các driver |

## Lưu ý

1. **TTL Management**: Mỗi driver có thể có cách xử lý TTL khác nhau. Memory driver có automatic cleanup, trong khi File driver kiểm tra TTL khi truy cập.

2. **Thread Safety**: Tất cả các phương thức của Manager đều thread-safe và có thể được gọi đồng thời từ nhiều goroutine.

3. **Error Handling**: Luôn kiểm tra error return từ các phương thức Set, Delete, và các batch operations.

4. **Driver Selection**: Nếu không có driver mặc định được thiết lập, các phương thức cache sẽ trả về error.

5. **Resource Management**: Luôn gọi `Close()` khi kết thúc để giải phóng tài nguyên của tất cả drivers.

6. **Configuration**: Mỗi driver có thể yêu cầu cấu hình riêng. Tham khảo documentation của từng driver để biết chi tiết.

7. **Performance**: Batch operations thường hiệu quả hơn multiple single operations, đặc biệt với Redis và MongoDB drivers.

## Phát triển

### Mock cho Testing

Package này cung cấp mock cho việc testing trong thư mục `mocks`. Sử dụng MockManager để test các thành phần phụ thuộc vào cache mà không cần backend thật:

```go
import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "go.fork.vn/providers/cache/mocks"
)

func TestYourFunction(t *testing.T) {
    // Tạo mock manager
    mockManager := mocks.NewMockManager(t)
    
    // Thiết lập expectations cơ bản
    mockManager.On("Get", "user:1").Return(userData, true)
    mockManager.On("Set", "user:1", mock.Anything, mock.Anything).Return(nil)
    mockManager.On("Has", "user:1").Return(true)
    mockManager.On("Delete", "user:1").Return(nil)
    
    // Mock cho Remember pattern
    mockManager.On("Remember", "expensive_key", mock.Anything, mock.AnythingOfType("func() (interface{}, error)")).
        Run(func(args mock.Arguments) {
            // Giả lập thực thi callback function
            callback := args.Get(2).(func() (interface{}, error))
            // Chạy callback
            callback()
        }).
        Return(mockData, nil)
    
    // Mock cho batch operations
    mockManager.On("GetMultiple", []string{"user:1", "user:2"}).
        Return(map[string]interface{}{"user:1": userData}, []string{"user:2"})
    
    mockManager.On("SetMultiple", mock.Anything, mock.Anything).Return(nil)
    mockManager.On("DeleteMultiple", mock.Anything).Return(nil)
    
    // Mock cho Clear và Close
    mockManager.On("Clear").Return(nil)
    mockManager.On("Close").Return(nil)
    
    // Sử dụng mock trong tests
    err := YourFunction(mockManager)
    
    // Kiểm tra kết quả
    assert.NoError(t, err)
    mockManager.AssertExpectations(t)
}

// Ví dụ về hàm sử dụng Cache Manager
func YourFunction(m cache.Manager) error {
    // Set cache
    if err := m.Set("test", "data", 1*time.Hour); err != nil {
        return err
    }
    
    // Get from cache
    data, exists := m.Get("test")
    if !exists {
        return fmt.Errorf("cache miss")
    }
    
    // Do something with data...
    
    return nil
}
```

### Tạo lại Mocks

Mocks được tạo bằng [mockery](https://github.com/vektra/mockery). Để tạo lại mocks, chạy lệnh sau từ thư mục gốc của project:

```bash
mockery
```

Lệnh này sẽ sử dụng cấu hình từ file `.mockery.yaml`.

### Phương pháp cải thiện test coverage

Để cải thiện test coverage của package, hãy chú ý đến các phương pháp sau:

1. **Thiết lập test helper**: Tạo các hàm helper để thiết lập và dọn dẹp môi trường test một cách nhất quán.

2. **Mock external dependencies**: Sử dụng mock cho các dependency bên ngoài như Redis client, MongoDB client để không phụ thuộc vào service thật trong unit tests.

3. **Kiểm tra cả happy path và error path**: Đảm bảo kiểm tra cả trường hợp thành công và thất bại của mỗi hàm.

4. **Sử dụng testify**: Sử dụng các assertion của package testify để làm cho tests dễ đọc hơn.

5. **Docker containers cho integration tests**: Sử dụng Docker để chạy Redis/MongoDB tạm thời cho integration tests.

Ví dụ thiết lập test helper:

```go
// testHelper.go
package cache_test

import (
    "testing"
    "time"
    
    "go.fork.vn/providers/cache"
    "go.fork.vn/providers/cache/driver"
    "go.fork.vn/providers/cache/config"
    "github.com/stretchr/testify/require"
)

func setupTestManager(t *testing.T) (cache.Manager, func()) {
    // Tạo manager cho test
    manager := cache.NewManager()
    
    // Thêm memory driver cho test
    memConfig := config.DriverMemoryConfig{
        DefaultTTL:      60, // 1 minute for fast test
        CleanupInterval: 10, // 10 seconds
        MaxItems:        100,
    }
    memDriver := driver.NewMemoryDriver(memConfig)
    manager.AddDriver("test_memory", memDriver)
    manager.SetDefaultDriver("test_memory")
    
    // Tạo cleanup function
    cleanup := func() {
        err := manager.Clear()
        require.NoError(t, err)
        err = manager.Close()
        require.NoError(t, err)
    }
    
    return manager, cleanup
}
```
