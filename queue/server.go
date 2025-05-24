// filepath: /Users/cluster/dev/go/github.com/go-fork/providers/queue/server.go
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

// ServerOptions chứa các tùy chọn cấu hình cho server.
type ServerOptions struct {
	// Concurrency xác định số lượng worker xử lý tác vụ song song.
	Concurrency int

	// PollingInterval xác định thời gian chờ giữa các lần kiểm tra tác vụ (tính bằng mili giây).
	PollingInterval int

	// DefaultQueue xác định tên queue mặc định nếu không có queue nào được chỉ định.
	DefaultQueue string

	// StrictPriority xác định liệu có nên ưu tiên nghiêm ngặt giữa các hàng đợi.
	StrictPriority bool

	// Queues xác định danh sách các queue cần lắng nghe theo thứ tự ưu tiên.
	Queues []string

	// ShutdownTimeout xác định thời gian chờ để các worker hoàn tất tác vụ khi dừng server.
	ShutdownTimeout time.Duration

	// LogLevel xác định mức log.
	LogLevel int

	// RetryLimit xác định số lần thử lại tối đa cho tác vụ bị lỗi.
	RetryLimit int
}

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
	queue           adapter.QueueAdapter
	handlers        sync.Map
	started         bool
	stopCh          chan struct{}
	workerDoneCh    chan struct{}
	schedulerDoneCh chan struct{}
	mu              sync.Mutex
	options         ServerOptions
	queues          []string
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

	// Sử dụng danh sách các queue từ cấu hình
	queues := opts.Queues

	// Nếu không có queue nào được chỉ định, sử dụng queue mặc định
	if len(queues) == 0 {
		defaultQueue := "default"
		if opts.DefaultQueue != "" {
			defaultQueue = opts.DefaultQueue
		}
		queues = append(queues, defaultQueue)
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

// NewServerWithAdapter tạo một Server mới với adapter QueueAdapter được cung cấp.
func NewServerWithAdapter(adapter adapter.QueueAdapter, opts ServerOptions) Server {
	// Sử dụng danh sách các queue từ cấu hình
	queues := opts.Queues

	// Nếu không có queue nào được chỉ định, sử dụng queue mặc định
	if len(queues) == 0 {
		defaultQueue := "default"
		if opts.DefaultQueue != "" {
			defaultQueue = opts.DefaultQueue
		}
		queues = append(queues, defaultQueue)
	}

	return &ServerImpl{
		queue:           adapter,
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
		return fmt.Errorf("server already started")
	}

	s.started = true
	log.Println("Starting queue worker server...")

	// TODO: Implement worker start logic

	return nil
}

// Stop dừng xử lý tác vụ.
func (s *ServerImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return fmt.Errorf("server not started")
	}

	log.Println("Stopping queue worker server...")
	close(s.stopCh)

	// TODO: Implement worker stop logic

	s.started = false
	return nil
}
