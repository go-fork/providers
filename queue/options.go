package queue

import "time"

// Option là một hàm để cấu hình tác vụ.
type Option func(*TaskOptions)

// TaskOptions chứa các tùy chọn khi đưa tác vụ vào hàng đợi.
type TaskOptions struct {
	// Queue là tên hàng đợi cho tác vụ
	Queue string

	// MaxRetry là số lần thử lại tối đa nếu tác vụ thất bại
	MaxRetry int

	// Timeout là thời gian tối đa để tác vụ hoàn thành
	Timeout time.Duration

	// Deadline là thời hạn chót để tác vụ hoàn thành
	Deadline time.Time

	// Delay là thời gian trì hoãn trước khi xử lý tác vụ
	Delay time.Duration

	// ProcessAt là thời điểm cụ thể để xử lý tác vụ
	ProcessAt time.Time

	// TaskID là ID tùy chỉnh cho tác vụ
	TaskID string
}

// WithQueue đặt tên hàng đợi cho tác vụ.
func WithQueue(queue string) Option {
	return func(o *TaskOptions) {
		o.Queue = queue
	}
}

// WithMaxRetry đặt số lần thử lại tối đa cho tác vụ.
func WithMaxRetry(n int) Option {
	return func(o *TaskOptions) {
		o.MaxRetry = n
	}
}

// WithTimeout đặt thời gian timeout cho tác vụ.
func WithTimeout(d time.Duration) Option {
	return func(o *TaskOptions) {
		o.Timeout = d
	}
}

// WithDeadline đặt thời hạn thực hiện cho tác vụ.
func WithDeadline(t time.Time) Option {
	return func(o *TaskOptions) {
		o.Deadline = t
	}
}

// WithDelay đặt thời gian trì hoãn cho tác vụ.
func WithDelay(d time.Duration) Option {
	return func(o *TaskOptions) {
		o.Delay = d
	}
}

// WithProcessAt đặt thời điểm xử lý cho tác vụ.
func WithProcessAt(t time.Time) Option {
	return func(o *TaskOptions) {
		o.ProcessAt = t
	}
}

// WithTaskID đặt ID tùy chỉnh cho tác vụ.
func WithTaskID(id string) Option {
	return func(o *TaskOptions) {
		o.TaskID = id
	}
}

// GetDefaultOptions trả về các tùy chọn mặc định.
func GetDefaultOptions() *TaskOptions {
	return &TaskOptions{
		Queue:    "default",
		MaxRetry: 3,
		Timeout:  30 * time.Minute,
	}
}

// ApplyOptions áp dụng các tùy chọn vào TaskOptions.
func ApplyOptions(opts ...Option) *TaskOptions {
	options := GetDefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}
