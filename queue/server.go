package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-fork/providers/queue/adapter"
	"github.com/redis/go-redis/v9"
)

// HandlerFunc là một hàm xử lý tác vụ.
type HandlerFunc func(ctx context.Context, task *Task) error

// Server là interface cho việc xử lý tác vụ từ hàng đợi.
type Server interface {
	// RegisterHandler đăng ký một handler cho một loại tác vụ.
	RegisterHandler(taskName string, handler HandlerFunc)

	// RegisterHandlers đăng ký nhiều handler cùng một lúc.
	RegisterHandlers(handlers map[string]HandlerFunc)

	// Start bắt đầu xử lý tác vụ (worker).
	Start() error

	// Stop dừng xử lý tác vụ.
	Stop() error
}

// ServerImpl triển khai interface Server.
type ServerImpl struct {
	queue           adapter.IQueue
	handlers        sync.Map
	started         bool
	stopCh          chan struct{}
	workerDoneCh    chan struct{}
	schedulerDoneCh chan struct{}
	mu              sync.Mutex
	options         ServerOptions
	queues          []string
}

// ServerOptions chứa các tùy chọn cấu hình cho server.
type ServerOptions struct {
	// Concurrency xác định số lượng worker xử lý tác vụ song song.
	Concurrency int

	// Queues xác định các hàng đợi và mức ưu tiên của chúng.
	Queues map[string]int

	// StrictPriority xác định liệu có nên ưu tiên nghiêm ngặt giữa các hàng đợi.
	StrictPriority bool

	// ShutdownTimeout xác định thời gian chờ để các worker hoàn tất tác vụ khi dừng server.
	ShutdownTimeout time.Duration

	// LogLevel xác định mức log.
	LogLevel int

	// RetryLimit xác định số lần thử lại tối đa cho tác vụ bị lỗi.
	RetryLimit int
}

// NewServer tạo một Server mới.
func NewServer(redisClient redis.UniversalClient, opts ServerOptions) Server {
	// Khởi tạo với client thông thường
	redisStdClient, ok := redisClient.(*redis.Client)
	if !ok {
		// Fallback cho các trường hợp client không phải là *redis.Client
		redisStdClient = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
	}

	queue := adapter.NewRedisQueue(redisStdClient, "queue:")

	// Tạo danh sách các queue từ cấu hình
	queues := make([]string, 0, len(opts.Queues))
	for q := range opts.Queues {
		queues = append(queues, q)
	}

	// Nếu không có queue nào được chỉ định, sử dụng queue mặc định
	if len(queues) == 0 {
		queues = append(queues, "default")
	}

	return &ServerImpl{
		queue:           queue,
		handlers:        sync.Map{},
		started:         false,
		stopCh:          make(chan struct{}),
		workerDoneCh:    make(chan struct{}),
		schedulerDoneCh: make(chan struct{}),
		options:         opts,
		queues:          queues,
	}
}

// NewMemoryServer tạo một Server mới với bộ nhớ trong.
func NewMemoryServer(opts ServerOptions) Server {
	queue := adapter.NewMemoryQueue("queue:")

	// Tạo danh sách các queue từ cấu hình
	queues := make([]string, 0, len(opts.Queues))
	for q := range opts.Queues {
		queues = append(queues, q)
	}

	// Nếu không có queue nào được chỉ định, sử dụng queue mặc định
	if len(queues) == 0 {
		queues = append(queues, "default")
	}

	return &ServerImpl{
		queue:           queue,
		handlers:        sync.Map{},
		started:         false,
		stopCh:          make(chan struct{}),
		workerDoneCh:    make(chan struct{}),
		schedulerDoneCh: make(chan struct{}),
		options:         opts,
		queues:          queues,
	}
}

// RegisterHandler đăng ký một handler cho một loại tác vụ.
func (s *ServerImpl) RegisterHandler(taskName string, handler HandlerFunc) {
	s.handlers.Store(taskName, handler)
}

// RegisterHandlers đăng ký nhiều handler cùng một lúc.
func (s *ServerImpl) RegisterHandlers(handlers map[string]HandlerFunc) {
	for taskName, handler := range handlers {
		s.RegisterHandler(taskName, handler)
	}
}

// Start bắt đầu xử lý tác vụ (worker).
func (s *ServerImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return nil
	}

	concurrency := s.options.Concurrency
	if concurrency <= 0 {
		concurrency = 10 // Giá trị mặc định
	}

	// Khởi động scheduler để chuyển các tác vụ đã lập lịch vào queue xử lý
	go s.runScheduler()

	// Khởi động các worker để xử lý tác vụ
	for i := 0; i < concurrency; i++ {
		go s.runWorker(i)
	}

	s.started = true
	return nil
}

// runScheduler kiểm tra các tác vụ đã lập lịch và chuyển chúng vào queue xử lý
// khi đến thời điểm thực hiện.
func (s *ServerImpl) runScheduler() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer close(s.schedulerDoneCh)

	for {
		select {
		case <-ticker.C:
			s.processScheduledTasks()
		case <-s.stopCh:
			return
		}
	}
}

// processScheduledTasks di chuyển các tác vụ đã lên lịch từ scheduled queue vào pending queue
// khi đến thời điểm xử lý.
func (s *ServerImpl) processScheduledTasks() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()

	for _, queueName := range s.queues {
		scheduledQueueName := fmt.Sprintf("%s:scheduled", queueName)
		pendingQueueName := fmt.Sprintf("%s:pending", queueName)

		// Kiểm tra nếu không có tác vụ đã lập lịch nào
		isEmpty, err := s.queue.IsEmpty(ctx, scheduledQueueName)
		if err != nil || isEmpty {
			continue
		}

		// Kiểm tra tất cả các tác vụ đã lập lịch
		var maxCount int64 = 100 // Giới hạn số lượng tác vụ kiểm tra mỗi lần
		for i := int64(0); i < maxCount; i++ {
			var scheduledTask struct {
				TaskID    string    `json:"task_id"`
				ProcessAt time.Time `json:"process_at"`
			}

			// Lấy tác vụ từ scheduled queue mà không xóa
			// Chúng ta chỉ cần kiểm tra và không xóa nó vội
			err := s.queue.Dequeue(ctx, scheduledQueueName, &scheduledTask)
			if err != nil {
				// Không còn tác vụ nào trong scheduled queue
				break
			}

			// Nếu chưa đến thời gian xử lý
			if scheduledTask.ProcessAt.After(now) {
				// Đưa lại tác vụ vào queue để kiểm tra sau
				_ = s.queue.Enqueue(ctx, scheduledQueueName, scheduledTask)
				continue
			} // Lấy tác vụ từ queue và chuyển vào pending để xử lý
			var task Task
			// Chúng ta cần duyệt qua tất cả các tác vụ trong queue để tìm đúng tác vụ theo ID
			// Trong triển khai thực tế, chúng ta có thể dùng Redis ZPOP để lấy theo thời gian
			_ = s.queue.Enqueue(ctx, pendingQueueName, &task)
		}
	}
}

// runWorker xử lý các tác vụ từ queue.
func (s *ServerImpl) runWorker(id int) {
	defer func() {
		if id == 0 { // Chỉ worker đầu tiên sẽ đóng channel
			close(s.workerDoneCh)
		}
	}()

	for {
		select {
		case <-s.stopCh:
			return
		default:
			// Kiểm tra và xử lý tác vụ từ tất cả các queue theo thứ tự ưu tiên
			s.processNextTask(id)

			// Tránh CPU spike
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// processNextTask xử lý tác vụ tiếp theo từ queue.
func (s *ServerImpl) processNextTask(workerID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Nếu strict priority, xử lý theo thứ tự ưu tiên
	var queueToProcess []string
	if s.options.StrictPriority {
		// Sắp xếp các queue theo ưu tiên (không cần thiết nếu đã sắp xếp khi khởi tạo)
		queueToProcess = s.queues
	} else {
		// Chọn ngẫu nhiên một queue để xử lý
		if len(s.queues) > 0 {
			queueToProcess = []string{s.queues[workerID%len(s.queues)]}
		}
	}

	for _, queueName := range queueToProcess {
		pendingQueue := fmt.Sprintf("%s:pending", queueName)

		// Kiểm tra nếu queue rỗng
		isEmpty, err := s.queue.IsEmpty(ctx, pendingQueue)
		if err != nil || isEmpty {
			continue
		}

		var task Task
		err = s.queue.Dequeue(ctx, pendingQueue, &task)
		if err != nil {
			log.Printf("Worker %d: Error dequeueing task: %v", workerID, err)
			continue
		}

		// Tìm handler phù hợp
		handlerVal, ok := s.handlers.Load(task.Name)
		if !ok {
			log.Printf("Worker %d: No handler found for task: %s", workerID, task.Name)
			continue
		}

		handler, ok := handlerVal.(HandlerFunc)
		if !ok {
			log.Printf("Worker %d: Invalid handler type for task: %s", workerID, task.Name)
			continue
		}

		// Thực thi handler
		taskCtx, taskCancel := context.WithTimeout(context.Background(), 30*time.Second)
		err = handler(taskCtx, &task)
		taskCancel()

		if err != nil {
			log.Printf("Worker %d: Error processing task %s: %v", workerID, task.ID, err)

			// Xử lý retry nếu không vượt quá số lần cho phép
			if task.RetryCount < task.MaxRetry {
				task.RetryCount++

				// Trì hoãn theo số lần retry (backoff strategy)
				delay := time.Duration(task.RetryCount*task.RetryCount) * time.Second
				retryTime := time.Now().Add(delay)

				// Lưu lại tác vụ để retry
				scheduledTask := struct {
					TaskID    string    `json:"task_id"`
					ProcessAt time.Time `json:"process_at"`
				}{
					TaskID:    task.ID,
					ProcessAt: retryTime,
				}

				scheduledQueue := fmt.Sprintf("%s:scheduled", queueName)
				_ = s.queue.Enqueue(ctx, scheduledQueue, scheduledTask)
				_ = s.queue.Enqueue(ctx, pendingQueue, &task)

				log.Printf("Worker %d: Task %s scheduled for retry %d/%d after %v",
					workerID, task.ID, task.RetryCount, task.MaxRetry, delay)
			} else {
				log.Printf("Worker %d: Task %s failed permanently after %d retries",
					workerID, task.ID, task.MaxRetry)

				// Có thể lưu tác vụ thất bại vào một queue riêng để theo dõi
				failedQueue := fmt.Sprintf("%s:failed", queueName)
				_ = s.queue.Enqueue(ctx, failedQueue, &task)
			}
		} else {
			log.Printf("Worker %d: Task %s completed successfully", workerID, task.ID)
		}

		// Đã xử lý được một tác vụ, thoát vòng lặp
		break
	}
}

// Stop dừng xử lý tác vụ.
func (s *ServerImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return nil
	}

	// Gửi tín hiệu dừng tới tất cả các worker
	close(s.stopCh)

	// Thiết lập timeout cho việc shutdown
	timeout := s.options.ShutdownTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second // Mặc định 30s
	}

	// Chờ cho tất cả các worker dừng lại
	select {
	case <-s.workerDoneCh:
		// Workers đã dừng
	case <-time.After(timeout):
		// Hết thời gian timeout
		log.Printf("Shutdown timeout: Some workers did not stop gracefully")
	}

	// Chờ cho scheduler dừng lại
	select {
	case <-s.schedulerDoneCh:
		// Scheduler đã dừng
	case <-time.After(5 * time.Second):
		// Thời gian timeout ngắn hơn cho scheduler
		log.Printf("Shutdown timeout: Scheduler did not stop gracefully")
	}

	// Khởi tạo lại các channel để có thể restart
	s.stopCh = make(chan struct{})
	s.workerDoneCh = make(chan struct{})
	s.schedulerDoneCh = make(chan struct{})

	s.started = false
	return nil
}
