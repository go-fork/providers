// Package queue cung cấp một service provider để quản lý hàng đợi và
// xử lý tác vụ bất đồng bộ với thiết kế đơn giản, linh hoạt hỗ trợ cả Redis và bộ nhớ trong.
//
// Package này cung cấp APIs để đưa tác vụ vào hàng đợi, lên lịch thực thi
// tác vụ và xử lý tác vụ qua các handlers đăng ký. Nó hỗ trợ:
//
// - Thực thi ngay lập tức (immediate execution)
// - Trì hoãn thực thi theo khoảng thời gian (delayed execution)
// - Lên lịch thực thi vào thời điểm cụ thể (scheduled execution)
// - Xử lý tác vụ song song với mức độ song song có thể cấu hình
// - Hệ thống thử lại tự động với chiến lược backoff
// - Ưu tiên giữa các hàng đợi khác nhau
// - Hàng đợi trong bộ nhớ cho môi trường phát triển và kiểm thử
//
// Package cung cấp ba thành phần chính:
//
// 1. Manager: Quản lý các thành phần trong queue và cung cấp cấu hình
// 2. Client: Cho phép đưa tác vụ vào hàng đợi (producer)
// 3. Server: Xử lý các tác vụ từ hàng đợi thông qua handlers (consumer)
//
// Hỗ trợ hai adapter chính:
//
// 1. Redis adapter: Cho môi trường production với khả năng mở rộng cao
// 2. Memory adapter: Cho môi trường phát triển và kiểm thử
//
// Ví dụ sử dụng:
//
//	// Đăng ký service provider với Redis
//	app.Register(queue.NewServiceProvider(queue.RedisOptions{
//	    Addr: "localhost:6379",
//	}, queue.ServerOptions{
//	    Concurrency: 10,
//	    Queues: map[string]int{
//	        "default": 1,
//	        "emails":  2,
//	    },
//	}))
//
//	// Sử dụng client để đưa tác vụ vào hàng đợi
//	client := app.Make("queue.client").(queue.Client)
//	client.Enqueue("email:send", map[string]interface{}{
//	    "to":      "user@example.com",
//	    "subject": "Hello",
//	}, queue.WithQueue("emails"), queue.WithMaxRetry(3))
//
//	// Sử dụng server để xử lý tác vụ
//	server := app.Make("queue.server").(queue.Server)
//	server.RegisterHandler("email:send", func(ctx context.Context, task *queue.Task) error {
//	    var payload map[string]interface{}
//	    if err := task.Unmarshal(&payload); err != nil {
//	        return err
//	    }
//	    // Xử lý logic gửi email ở đây...
//	    return nil
//	})
//
//	// Bắt đầu xử lý tác vụ
//	if err := server.Start(); err != nil {
//	    log.Fatal(err)
//	}
//
// Sử dụng bộ nhớ trong (cho phát triển và kiểm thử):
//
//	// Khởi tạo client với bộ nhớ trong
//	client := queue.NewMemoryClient()
//
//	// Khởi tạo server với bộ nhớ trong
//	server := queue.NewMemoryServer(queue.ServerOptions{
//	    Concurrency: 5,
//	})
//
// Xử lý tác vụ theo lịch:
//
//	// Lên lịch tác vụ sau 5 phút
//	client.EnqueueIn("report:generate", 5*time.Minute, payload)
//
//	// Lên lịch tác vụ vào thời điểm cụ thể
//	processAt := time.Date(2025, 5, 24, 15, 0, 0, 0, time.Local)
//	client.EnqueueAt("cleanup:old-data", processAt, payload)
//
// Các tùy chọn khác:
//
//	// Tùy chỉnh queue, số lần thử lại, timeout...
//	client.Enqueue("image:resize", payload,
//	    queue.WithQueue("media"),
//	    queue.WithMaxRetry(3),
//	    queue.WithTimeout(2*time.Minute),
//	    queue.WithTaskID("resize-123"),
//	)
package queue
