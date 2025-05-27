# Scheduler Provider

Scheduler Provider là giải pháp lên lịch và chạy các task định kỳ cho ứng dụng Go, được xây dựng dựa trên thư viện [go-co-op/gocron](https://github.com/go-co-op/gocron).

## Tính năng nổi bật

- **Configuration-driven**: Hỗ trợ cấu hình qua file config với struct Config và RedisLockerOptions
- **Auto-start**: Tự động khởi động scheduler khi ứng dụng boot (có thể tắt qua config)
- **Distributed Locking**: Hỗ trợ distributed locking với Redis cho môi trường phân tán
- **Fluent API**: API fluent cho trải nghiệm lập trình dễ dàng
- **Tích hợp DI**: Tích hợp toàn bộ tính năng vào DI container của ứng dụng
- **Multiple Schedule Types**: Hỗ trợ nhiều loại lịch trình (khoảng thời gian, thời điểm cụ thể, cron expressions)
- **Singleton Mode**: Hỗ trợ chế độ singleton để tránh chạy song song cùng một task
- **Task Management**: Hỗ trợ tag để nhóm và quản lý các task

## Cài đặt

```bash
go get github.com/go-fork/providers/scheduler
```

## Cấu hình

### File cấu hình (config/app.yaml)

```yaml
scheduler:
  # Tự động khởi động scheduler khi ứng dụng boot
  auto_start: true

  # Distributed locking với Redis (tùy chọn)
  distributed_lock:
    enabled: false
  
  # Cài đặt RedisLockerOptions cho distributed locking
  options:
    key_prefix: "scheduler_lock:"
    lock_duration: 30      # seconds
    max_retries: 3
    retry_delay: 100       # milliseconds
```

### Các tùy chọn cấu hình

| Field | Type | Mô tả | Mặc định |
|-------|------|-------|----------|
| `auto_start` | bool | Tự động khởi động scheduler trong Boot() | `true` |
| `distributed_lock.enabled` | bool | Bật distributed locking với Redis | `false` |
| `options.key_prefix` | string | Tiền tố key trong Redis | `"scheduler_lock:"` |
| `options.lock_duration` | int | Thời gian lock (giây) | `30` |
| `options.max_retries` | int | Số lần thử lại | `3` |
| `options.retry_delay` | int | Thời gian chờ giữa các lần thử (ms) | `100` |

## Cách sử dụng

### 1. Đăng ký Service Provider

```go
package main

import (
    "github.com/go-fork/di"
    "github.com/go-fork/providers/scheduler"
)

func main() {
    app := di.New()
    app.Register(scheduler.NewServiceProvider())
    
    // Khởi động ứng dụng
    app.Boot()
    
    // Giữ ứng dụng chạy để scheduler có thể hoạt động
    select {}
}
```

### 2. Lấy scheduler từ container và lên lịch cho task

```go
// Lấy scheduler từ container
container := app.Container()
sched := container.Get("scheduler").(scheduler.Manager)

// Đăng ký task chạy mỗi 5 phút
sched.Every(5).Minutes().Do(func() {
    fmt.Println("Task runs every 5 minutes")
})

// Đăng ký task với cron expression
sched.Cron("0 0 * * *").Do(func() {
    fmt.Println("Task runs at midnight every day")
})

// Đăng ký task với tag để dễ quản lý
sched.Every(1).Hour().Tag("maintenance").Do(func() {
    fmt.Println("Maintenance task runs hourly")
})
```

### 3. Sử dụng Configuration-driven Scheduler

Scheduler provider hỗ trợ cấu hình tự động thông qua file config. Khi có cấu hình distributed locking, provider sẽ tự động tạo và thiết lập Redis locker:

```yaml
# config/app.yaml
scheduler:
  auto_start: true
  distributed_lock:
    enabled: true
  options:
    key_prefix: "myapp_scheduler:"
    lock_duration: 60      # seconds
    max_retries: 5
    retry_delay: 200       # milliseconds

# Cần cả redis provider để kết nối Redis
redis:
  default:
    addr: "localhost:6379"
    password: ""
    db: 0
```

```go
import (
    "github.com/go-fork/di"
    "github.com/go-fork/providers/config"
    "github.com/go-fork/providers/redis"
    "github.com/go-fork/providers/scheduler"
)

func main() {
    app := di.New()
    
    // Đăng ký các providers theo thứ tự phụ thuộc
    app.Register(config.NewServiceProvider())
    app.Register(redis.NewServiceProvider())  // Required cho distributed locking
    app.Register(scheduler.NewServiceProvider())
    
    // Khởi động ứng dụng - scheduler sẽ tự động cấu hình distributed locking
    app.Boot()
    
    // Sử dụng scheduler đã được cấu hình sẵn
    container := app.Container()
    sched := container.Get("scheduler").(scheduler.Manager)
    
    // Tất cả jobs sẽ tự động sử dụng distributed locking nếu được bật
    sched.Every(5).Minutes().Do(func() {
        fmt.Println("This task uses distributed locking automatically")
    })
    
    select {}
}
```

### 4. Sử dụng Manual Redis Locker (Tùy chọn)

Nếu bạn muốn tự thiết lập Redis locker thay vì dùng config:

```go
import (
    "github.com/redis/go-redis/v9"
    "github.com/go-fork/providers/scheduler"
)

// Khởi tạo Redis client
redisClient := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

// Tạo Redis Locker với tùy chọn tùy chỉnh
customLocker, err := scheduler.NewRedisLocker(redisClient, scheduler.RedisLockerOptions{
    KeyPrefix:    "myapp_scheduler:",
    LockDuration: 60,   // seconds (int value)
    MaxRetries:   5,
    RetryDelay:   200,  // milliseconds (int value)
})
if err != nil {
    log.Fatal(err)
}

// Thiết lập Redis Locker cho scheduler
sched := container.Get("scheduler").(scheduler.Manager)
sched.WithDistributedLocker(customLocker)
```

Các tùy chọn cấu hình của Redis Locker:

| Tùy chọn | Mô tả | Giá trị mặc định |
|----------|-------|------------------|
| KeyPrefix | Tiền tố được thêm vào trước mỗi khóa trong Redis | `scheduler_lock:` |
| LockDuration | Thời gian một khóa sẽ tồn tại trước khi tự động hết hạn (giây) | `30` |
| MaxRetries | Số lần thử tối đa khi gặp lỗi khi tương tác với Redis | `3` |
| RetryDelay | Thời gian chờ giữa các lần thử (milliseconds) | `100` |

### 5. Quản lý các task

```go
// Xóa task theo tag
sched.RemoveByTag("maintenance")

// Xóa tất cả task
sched.Clear()

// Tạm dừng scheduler
sched.Stop()

// Khởi động lại scheduler
sched.Start()
```

### 6. Tùy chỉnh Scheduler

```go
// Đặt thời gian múi giờ cho scheduler
sched.Location(time.UTC)

// Đặt singleton mode cho scheduler 
// (không chạy nhiều instance của cùng một job)
sched.SingletonMode().Every(1).Minute().Do(func() {
    // Task sẽ không chạy song song nếu lần chạy trước chưa hoàn thành
    longRunningTask()
})
```

## Tính năng tự động gia hạn khóa Redis

Khi sử dụng Redis Distributed Locker, scheduler triển khai cơ chế tự động gia hạn khóa để đảm bảo job không bị gián đoạn khi chạy thời gian dài:

- Khóa Redis được tự động gia hạn sau khi đã chạy được 2/3 thời gian hết hạn
- Việc gia hạn xảy ra trong một goroutine riêng biệt
- Khi job hoàn thành, khóa sẽ được giải phóng
- Nếu instance gặp sự cố, khóa sẽ tự động hết hạn sau LockDuration

## Các ví dụ nâng cao

### Chạy task với tham số

```go
sched.Every(1).Day().Do(func(name string) {
    fmt.Printf("Hello %s\n", name)
}, "John")
```

### Chạy task vào thời điểm cụ thể

```go
sched.Every(1).Day().At("10:30").Do(func() {
    fmt.Println("Task runs at 10:30 AM daily")
})

// Hoặc sử dụng cron expression
sched.Cron("30 10 * * *").Do(func() {
    fmt.Println("Task runs at 10:30 AM daily")
})
```

### Xử lý lỗi từ task

```go
job, err := sched.Every(1).Minute().Do(func() error {
    // Công việc có thể trả về lỗi
    if somethingWrong {
        return errors.New("something went wrong")
    }
    return nil
})

if err != nil {
    log.Fatal(err)
}

// Đăng ký hàm xử lý khi task trả về lỗi
job.OnError(func(err error) {
    log.Printf("Job failed: %v", err)
})
```

## Yêu cầu hệ thống

- Go 1.18 trở lên
- Redis (tùy chọn, chỉ khi sử dụng distributed locking)

## Giấy phép

Mã nguồn này được phân phối dưới giấy phép MIT.
