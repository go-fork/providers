// Package scheduler cung cấp giải pháp lên lịch và chạy các task định kỳ cho ứng dụng Go,
// dựa trên thư viện gocron.
//
// Tính năng nổi bật:
//   - Wrap toàn bộ tính năng của thư viện gocron - một thư viện lập lịch và chạy task hiệu quả
//   - Hỗ trợ nhiều loại lịch trình: theo khoảng thời gian, theo thời điểm cụ thể, biểu thức cron
//   - Hỗ trợ chế độ singleton để tránh chạy song song cùng một task
//   - Hỗ trợ distributed locking với Redis cho môi trường phân tán
//   - Hỗ trợ tag để nhóm và quản lý các task
//   - Tích hợp với DI container thông qua ServiceProvider
//   - API fluent cho trải nghiệm lập trình dễ dàng
//
// Kiến trúc và cách hoạt động:
//   - Sử dụng mô hình embedding để triển khai interface Manager trực tiếp nhúng gocron.Scheduler
//   - Cung cấp fluent interface để cấu hình task một cách dễ dàng và rõ ràng
//   - ServiceProvider giúp tích hợp dễ dàng vào ứng dụng thông qua DI container
//   - Tự động khởi động scheduler khi ứng dụng boot
//   - Hỗ trợ tự động gia hạn khóa cho distributed locking trong môi trường phân tán
//
// Ví dụ sử dụng cơ bản:
//
//	// Đăng ký service provider
//	app := di.New()
//	app.Register(scheduler.NewServiceProvider())
//
//	// Lấy scheduler từ container
//	container := app.Container()
//	sched := container.Get("scheduler").(scheduler.Manager)
//
//	// Đăng ký task chạy mỗi 5 phút
//	sched.Every(5).Minutes().Do(func() {
//		fmt.Println("Task runs every 5 minutes")
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
// Sử dụng với Distributed Locker (cho môi trường phân tán):
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
//	// Hoặc với tùy chọn tùy chỉnh
//	customLocker, err := scheduler.NewRedisLocker(redisClient, scheduler.RedisLockerOptions{
//		KeyPrefix:    "myapp_scheduler:",
//		LockDuration: 60 * time.Second,
//		MaxRetries:   5,
//		RetryDelay:   200 * time.Millisecond,
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
// đồng thời tích hợp dễ dàng với kiến trúc ứng dụng thông qua DI container.
package scheduler
