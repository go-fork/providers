// Package scheduler cung cấp giải pháp lên lịch và chạy các task định kỳ cho ứng dụng Go,
// dựa trên thư viện gocron.
//
// Tính năng nổi bật:
//   - Configuration-driven: Hỗ trợ cấu hình qua file config với struct Config và RedisLockerOptions
//   - Auto-start: Tự động khởi động scheduler khi ứng dụng boot (có thể tắt qua config)
//   - Wrap toàn bộ tính năng của thư viện gocron - một thư viện lập lịch và chạy task hiệu quả
//   - Hỗ trợ nhiều loại lịch trình: theo khoảng thời gian, theo thời điểm cụ thể, biểu thức cron
//   - Hỗ trợ chế độ singleton để tránh chạy song song cùng một task
//   - Hỗ trợ distributed locking với Redis cho môi trường phân tán (tự động cấu hình qua config)
//   - Hỗ trợ tag để nhóm và quản lý các task
//   - Tích hợp với DI container thông qua ServiceProvider
//   - API fluent cho trải nghiệm lập trình dễ dàng
//
// Kiến trúc và cách hoạt động:
//   - Sử dụng mô hình embedding để triển khai interface Manager trực tiếp nhúng gocron.Scheduler
//   - Cung cấp fluent interface để cấu hình task một cách dễ dàng và rõ ràng
//   - ServiceProvider giúp tích hợp dễ dàng vào ứng dụng thông qua DI container
//   - Hỗ trợ dual config system: RedisLockerOptions (int) cho file config và RedisLockerOptionsTime (time.Duration) cho internal use
//   - Tự động khởi động scheduler khi ứng dụng boot (có thể tắt thông qua config auto_start: false)
//   - Tự động thiết lập distributed locking khi được enable trong config và có redis provider
//   - Hỗ trợ tự động gia hạn khóa cho distributed locking trong môi trường phân tán
//
// Ví dụ sử dụng với configuration-driven approach:
//
//	// config/app.yaml
//	scheduler:
//	  auto_start: true
//	  distributed_lock:
//	    enabled: true
//	  options:
//	    key_prefix: "myapp_scheduler:"
//	    lock_duration: 60    # seconds
//	    max_retries: 5
//	    retry_delay: 200     # milliseconds
//
//	redis:
//	  default:
//	    addr: "localhost:6379"
//	    password: ""
//	    db: 0
//
//	// Đăng ký service providers
//	app := di.New()
//	app.Register(config.NewServiceProvider())
//	app.Register(redis.NewServiceProvider())  // Required cho distributed locking
//	app.Register(scheduler.NewServiceProvider())
//
//	// Boot ứng dụng - scheduler tự động cấu hình
//	app.Boot()
//
//	// Lấy scheduler từ container
//	container := app.Container()
//	sched := container.Get("scheduler").(scheduler.Manager)
//
//	// Đăng ký task chạy mỗi 5 phút với distributed locking tự động
//	sched.Every(5).Minutes().Do(func() {
//		fmt.Println("Task runs every 5 minutes with distributed locking")
//	})
//
//	// Đăng ký task với cron expression
//	sched.Cron("0 0 * * *").Do(func() {
//		fmt.Println("Task runs at midnight every day")
//	})
//
//	// Đăng ký task với tag để dễ quản lý
//	sched.Every(1).Hour().Tag("maintenance").Do(func() {
//		fmt.Println("Maintenance task runs hourly")
//	})
//
// Sử dụng Manual Redis Locker (tùy chọn thay vì config):
//
//	import (
//		"github.com/redis/go-redis/v9"
//	)
//
//	// Khởi tạo Redis client
//	redisClient := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "",
//		DB:       0,
//	})
//
//	// Tạo Redis Locker với tùy chọn mặc định
//	locker, err := scheduler.NewRedisLocker(redisClient)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Hoặc với tùy chọn tùy chỉnh (sử dụng int values cho config)
//	customLocker, err := scheduler.NewRedisLocker(redisClient, scheduler.RedisLockerOptions{
//		KeyPrefix:    "myapp_scheduler:",
//		LockDuration: 60,   // seconds (int value)
//		MaxRetries:   5,
//		RetryDelay:   200,  // milliseconds (int value)
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Lấy scheduler từ container
//	sched := container.Get("scheduler").(scheduler.Manager)
//
//	// Thiết lập Redis Locker cho scheduler
//	sched.WithDistributedLocker(locker)
//
//	// Từ bây giờ, tất cả các jobs sẽ sử dụng distributed locking với Redis
//	// để đảm bảo chỉ chạy một lần trong môi trường phân tán
//
// Gói này giúp đơn giản hóa việc lên lịch và chạy các task định kỳ trong ứng dụng Go,
// đồng thời tích hợp dễ dàng với kiến trúc ứng dụng thông qua DI container và hỗ trợ
// configuration-driven approach cho việc thiết lập distributed locking.
package scheduler
