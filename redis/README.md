# Redis Provider

Redis provider cung cấp quản lý kết nối Redis và hỗ trợ dependency injection cho ứng dụng Go.

## Tổng quan

Package này cung cấp một cách đơn giản để cấu hình và quản lý các kết nối Redis trong ứng dụng Go. Nó hỗ trợ:

- Redis client tiêu chuẩn
- Universal client cho Redis Cluster, Sentinel và triển khai standalone
- Dependency injection thông qua container `go-fork/di`
- Cấu hình thông qua package `go-fork/providers/config`
- Kiểm thử dễ dàng với hỗ trợ mock (sử dụng mockery)

## Cài đặt

```bash
go get github.com/go-fork/providers/redis
```

## Sử dụng

### Sử dụng cơ bản

```go
import (
    "context"
    "github.com/go-fork/providers/redis"
)

func main() {
    // Khởi tạo Redis manager với cấu hình mặc định
    manager := redis.NewManager()
    
    // Hoặc với cấu hình tùy chỉnh
    config := redis.DefaultConfig()
    config.Client.Host = "redis-server"
    config.Client.Port = 6379
    manager = redis.NewManagerWithConfig(config)
    
    // Lấy Redis client
    client, err := manager.Client()
    if err != nil {
        // Xử lý lỗi
    }
    
    // Sử dụng client để thực hiện các thao tác Redis
    ctx := context.Background()
    err = client.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        // Xử lý lỗi
    }
    
    // Kiểm tra kết nối
    err = manager.Ping(ctx)
    if err != nil {
        // Xử lý lỗi kết nối
    }
}
```

### Sử dụng với Dependency Injection

```go
import (
    "github.com/go-fork/di"
    "github.com/go-fork/providers/redis"
    "github.com/go-fork/providers/config"
)

// Định nghĩa một ứng dụng với container DI
type App struct {
    container *di.Container
}

func (a *App) Container() *di.Container {
    return a.container
}

func main() {
    // Khởi tạo container và app
    container := di.New()
    app := &App{container: container}
    
    // Đăng ký config provider trước
    configProvider := config.NewServiceProvider()
    configProvider.Register(app)
    
    // Đăng ký Redis provider
    redisProvider := redis.NewServiceProvider()
    redisProvider.Register(app)
    
    // Khởi động các provider
    configProvider.Boot(app)
    redisProvider.Boot(app)
    
    // Sau khi provider được đăng ký và khởi động, bạn có thể lấy Redis client:
    redisClient, err := container.Make("redis.client")
    if err != nil {
        // Xử lý lỗi
    }
    
    universalClient, err := container.Make("redis.universal")
    if err != nil {
        // Xử lý lỗi
    }
    
    manager, err := container.Make("redis.manager")
    if err != nil {
        // Xử lý lỗi
    }
}
```

### Cấu hình

Tạo file `redis.yaml` trong thư mục cấu hình của bạn:

```yaml
redis:
  # Cấu hình client tiêu chuẩn
  client:
    host: localhost
    port: 6379
    password: ""
    db: 0
    prefix: "app:"
    dial_timeout: 5
    read_timeout: 3
    write_timeout: 3
    pool_size: 10
    min_idle_conns: 5
  
  # Cấu hình universal client (cho cluster/sentinel)
  universal:
    addresses:
      - localhost:6379
    password: ""
    db: 0
    prefix: "app:"
    dial_timeout: 5
    read_timeout: 3
    write_timeout: 3
    max_retries: 3
    min_retry_backoff: 8
    max_retry_backoff: 512
    pool_size: 10
    min_idle_conns: 5
    cluster_mode: false
    master_name: "" # Chỉ cho chế độ sentinel
    route_randomly: false
```

### Testing với Mock

Package này hỗ trợ mocking thông qua Mockery:

```go
import (
    "testing"
    
    "github.com/go-fork/providers/redis"
    "github.com/go-fork/providers/redis/mocks"
    "github.com/stretchr/testify/assert"
)

func TestWithMockRedisManager(t *testing.T) {
    // Tạo mock cho Redis manager
    mockManager := mocks.NewMockManager(t)
    
    // Thiết lập expectation
    mockManager.EXPECT().Ping(mock.Anything).Return(nil)
    
    // Sử dụng mock trong test
    err := mockManager.Ping(context.Background())
    assert.NoError(t, err)
    
    // Kiểm tra xem tất cả các expectation có được gọi không
    mockManager.AssertExpectations(t)
}
```

## Đóng góp

Nếu bạn phát hiện lỗi hoặc có ý tưởng cải thiện, vui lòng tạo issue hoặc pull request trên GitHub.

## Giấy phép

Package này được cấp phép theo Giấy phép MIT.
