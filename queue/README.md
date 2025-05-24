# Queue Provider

Queue Provider là giải pháp xử lý hàng đợi và tác vụ nền đơn giản nhưng mạnh mẽ cho ứng dụng Go, không phụ thuộc vào thư viện bên ngoài.

## Tính năng nổi bật

- Triển khai đơn giản, dễ bảo trì và mở rộng
- Hỗ trợ Redis làm message broker và bộ nhớ trong cho môi trường phát triển
- Hỗ trợ đặt lịch công việc (ngay lập tức, sau một khoảng thời gian, vào một thời điểm)
- Tự động thử lại tác vụ thất bại với chiến lược backoff
- Tích hợp dễ dàng với DI container của ứng dụng
- API đơn giản và tiện lợi cho người sử dụng

## Cài đặt

Để cài đặt Queue Provider, bạn có thể sử dụng lệnh go get:

```bash
go get github.com/go-fork/providers/queue
```

## Cách sử dụng

### 1. Đăng ký Service Provider

```go
package main

import (
    "github.com/go-fork/di"
    "github.com/go-fork/providers/queue"
)

func main() {
    app := di.New()
    
    // Đăng ký service provider với cấu hình mặc định
    app.Register(queue.NewServiceProvider())
    
    // Hoặc đăng ký với cấu hình tùy chỉnh
    config := queue.Config{
        DefaultAdapter: "redis",
        RedisConfig: queue.RedisConfig{
            Addr:     "localhost:6379",
            Password: "",
            DB:       0,
        },
    }
    app.Register(queue.NewServiceProviderWithConfig(config))
    
    // Khởi động ứng dụng (sẽ tự động khởi động queue worker)
    app.Boot()
    
    // Giữ ứng dụng chạy để worker có thể xử lý tác vụ
    select {}
}
```

### 2. Thêm tác vụ vào hàng đợi (Producer)

```go
// Lấy queue manager từ container
container := app.Container()
manager := container.MustMake("queue").(*queue.Manager)

// Hoặc lấy trực tiếp client từ container
client := container.MustMake("queue.client").(queue.Client)

// Thêm tác vụ vào hàng đợi để xử lý ngay lập tức
payload := map[string]interface{}{
    "user_id": 123,
    "email":   "user@example.com",
}

taskInfo, err := client.Enqueue("email:welcome", payload)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Đã thêm tác vụ: %s\n", taskInfo.ID)

// Đặt lịch tác vụ chạy sau 5 phút
taskInfo, err = client.EnqueueIn("reminder:task", 5*time.Minute, payload)
if err != nil {
    log.Fatal(err)
}

// Đặt lịch tác vụ chạy vào thời điểm cụ thể
processAt := time.Date(2023, 5, 25, 13, 0, 0, 0, time.Local)
taskInfo, err = client.EnqueueAt("report:generate", processAt, payload)
if err != nil {
    log.Fatal(err)
}
```

### 3. Xử lý tác vụ từ hàng đợi (Consumer)

```go
// Lấy queue server từ container
server := container.MustMake("queue.server").(queue.Server)

// Đăng ký handler cho tác vụ "email:welcome"
server.RegisterHandler("email:welcome", func(ctx context.Context, task *queue.Task) error {
    var payload map[string]interface{}
    if err := task.Unmarshal(&payload); err != nil {
        return err
    }
    
    userID := int(payload["user_id"].(float64))
    email := payload["email"].(string)
    
    fmt.Printf("Gửi email chào mừng đến %s (ID: %d)\n", email, userID)
        
            // Xử lý logic gửi email ở đây...
    return nil
})

// Đăng ký handler cho các loại tác vụ khác
server.RegisterHandler("reminder:task", handleReminder)
server.RegisterHandler("report:generate", handleReportGeneration)

// Hoặc đăng ký nhiều handler cùng một lúc
server.RegisterHandlers(map[string]queue.HandlerFunc{
    "notification:push": handlePushNotification,
    "order:process":     handleOrderProcessing,
})

// Bắt đầu xử lý tác vụ
err := server.Start()
if err != nil {
    log.Fatal(err)
}

// Khi muốn dừng server
// server.Stop()
```

### 4. Tùy chọn cấu hình

```go
// Cấu hình nâng cao cho queue server
app.Register(queue.NewServiceProvider(
    // Redis connection options
    queue.RedisOptions{
        Addr:     "localhost:6379",
        Password: "secret",
        DB:       0,
        PoolSize: 100,
        TLS:      true,
    },
    // Queue server options
    queue.ServerOptions{
        Concurrency: 20,  // Số lượng worker chạy song song
        Queues: map[string]int{
            "critical": 6,  // Ưu tiên cao nhất
            "default":  3,  // Ưu tiên trung bình
            "low":      1,  // Ưu tiên thấp nhất
        },
        StrictPriority:  true,  // Ưu tiên nghiêm ngặt giữa các hàng đợi
        RetryLimit:      5,     // Số lần thử lại tối đa
        ShutdownTimeout: 1 * time.Minute,  // Thời gian chờ khi tắt server
    },
))
```

### 5. Tùy chọn khi thêm tác vụ

```go
// Thêm tác vụ với các tùy chọn
taskInfo, err := client.Enqueue("image:resize", payload,
    queue.WithQueue("media"),     // Chỉ định hàng đợi
    queue.WithMaxRetry(3),        // Số lần thử lại tối đa
    queue.WithTimeout(2*time.Minute), // Thời gian timeout
    queue.WithTaskID("resize-123"),   // Chỉ định ID cho tác vụ
)
```

### 6. Sử dụng bộ nhớ trong (cho môi trường phát triển)

```go
// Khởi tạo client với bộ nhớ trong
client := queue.NewMemoryClient()

// Khởi tạo server với bộ nhớ trong
server := queue.NewMemoryServer(queue.ServerOptions{
    Concurrency: 5,
    Queues: map[string]int{
        "default": 1,
    },
})
```

## Yêu cầu hệ thống

- Go 1.18 trở lên
- Redis 6.0 trở lên (nếu sử dụng Redis adapter)

## Giấy phép

Mã nguồn này được phân phối dưới giấy phép MIT.
