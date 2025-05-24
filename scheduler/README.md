# Scheduler Provider

Scheduler Provider là giải pháp lên lịch và chạy các task định kỳ cho ứng dụng Go, được xây dựng dựa trên thư viện [go-co-op/gocron](https://github.com/go-co-op/gocron).

## Tính năng nổi bật

- Tích hợp toàn bộ tính năng của thư viện gocron vào DI container của ứng dụng
- Hỗ trợ nhiều loại lịch trình: theo khoảng thời gian, theo thời điểm cụ thể, biểu thức cron
- Hỗ trợ chế độ singleton để tránh chạy song song cùng một task
- Hỗ trợ distributed locking với Redis cho môi trường phân tán
- Hỗ trợ tag để nhóm và quản lý các task
- API fluent cho trải nghiệm lập trình dễ dàng

## Cài đặt

Để cài đặt Scheduler Provider, bạn có thể sử dụng lệnh go get:

```bash
go get github.com/go-fork/providers/scheduler
```

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

### 3. Sử dụng Distributed Locker với Redis

Để đảm bảo task chỉ chạy một lần trong môi trường phân tán (nhiều instance của ứng dụng), bạn có thể sử dụng Redis Distributed Locker:

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

// Tạo Redis Locker với tùy chọn mặc định
locker, err := scheduler.NewRedisLocker(redisClient)
if err != nil {
    log.Fatal(err)
}

// Hoặc với tùy chọn tùy chỉnh
customLocker, err := scheduler.NewRedisLocker(redisClient, scheduler.RedisLockerOptions{
    KeyPrefix:    "myapp_scheduler:",
    LockDuration: 60 * time.Second,
    MaxRetries:   5,
    RetryDelay:   200 * time.Millisecond,
})
if err != nil {
    log.Fatal(err)
}

// Thiết lập Redis Locker cho scheduler
sched := container.Get("scheduler").(scheduler.Manager)
sched.WithDistributedLocker(locker)

// Từ bây giờ, tất cả các jobs sẽ sử dụng distributed locking với Redis
// để đảm bảo chỉ chạy một lần trong môi trường phân tán
```

Các tùy chọn cấu hình của Redis Locker:

| Tùy chọn | Mô tả | Giá trị mặc định |
|----------|-------|------------------|
| KeyPrefix | Tiền tố được thêm vào trước mỗi khóa trong Redis | `scheduler_lock:` |
| LockDuration | Thời gian một khóa sẽ tồn tại trước khi tự động hết hạn | `30 * time.Second` |
| MaxRetries | Số lần thử tối đa khi gặp lỗi khi tương tác với Redis | `3` |
| RetryDelay | Thời gian chờ giữa các lần thử | `100 * time.Millisecond` |

### 4. Quản lý các task

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

### 5. Tùy chỉnh Scheduler

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
