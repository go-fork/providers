package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

// Manager là interface chính cho việc quản lý lịch trình công việc, wrapping gocron.
type Manager interface {
	// WithDistributedLocker thiết lập distributed locker (như Redis) cho scheduler.
	// Hữu ích khi chạy scheduler trên nhiều máy chủ trong môi trường phân tán.
	WithDistributedLocker(locker gocron.Locker) Manager

	// Every tạo một công việc mới với khoảng thời gian được chỉ định.
	// Trả về Manager để hỗ trợ fluent interface.
	Every(interval interface{}) Manager

	// Second chỉ định đơn vị thời gian là giây (đơn lẻ).
	// Trả về Manager để hỗ trợ fluent interface.
	Second() Manager

	// Seconds chỉ định đơn vị thời gian là giây.
	// Trả về Manager để hỗ trợ fluent interface.
	Seconds() Manager

	// Minutes chỉ định đơn vị thời gian là phút.
	// Trả về Manager để hỗ trợ fluent interface.
	Minutes() Manager

	// Hours chỉ định đơn vị thời gian là giờ.
	// Trả về Manager để hỗ trợ fluent interface.
	Hours() Manager

	// Days chỉ định đơn vị thời gian là ngày.
	// Trả về Manager để hỗ trợ fluent interface.
	Days() Manager

	// Weeks chỉ định đơn vị thời gian là tuần.
	// Trả về Manager để hỗ trợ fluent interface.
	Weeks() Manager

	// At chỉ định thời điểm trong ngày để chạy công việc.
	// Định dạng: "HH:MM" hoặc "HH:MM:SS".
	// Trả về Manager để hỗ trợ fluent interface.
	At(time string) Manager

	// StartAt chỉ định thời điểm bắt đầu cho công việc.
	// Trả về Manager để hỗ trợ fluent interface.
	StartAt(time time.Time) Manager

	// Cron thiết lập biểu thức cron cho công việc.
	// Trả về Manager để hỗ trợ fluent interface.
	Cron(cronExpression string) Manager

	// CronWithSeconds thiết lập biểu thức cron có hỗ trợ giây.
	// Trả về Manager để hỗ trợ fluent interface.
	CronWithSeconds(cronExpression string) Manager

	// Tag đánh dấu công việc với các tag được chỉ định.
	// Trả về Manager để hỗ trợ fluent interface.
	Tag(tags ...string) Manager

	// SingletonMode đặt công việc ở chế độ singleton (không chạy đồng thời).
	// Trả về Manager để hỗ trợ fluent interface.
	SingletonMode() Manager

	// Do đặt hàm để thực thi cho công việc với các tham số tùy chọn.
	// Trả về Job và error nếu có.
	Do(jobFun interface{}, params ...interface{}) (*gocron.Job, error)

	// Name đặt tên cho công việc đang được cấu hình.
	// Trả về Manager để hỗ trợ fluent interface.
	Name(name string) Manager

	// RemoveByTag xóa các công việc theo tag.
	RemoveByTag(tag string) error

	// RemoveByTags xóa các công việc khớp với TẤT CẢ tags đã chỉ định.
	RemoveByTags(tags ...string) error

	// FindJobsByTag tìm công việc theo tag.
	FindJobsByTag(tags ...string) ([]*gocron.Job, error)

	// StartAsync bắt đầu scheduler trong một goroutine riêng.
	StartAsync()

	// StartBlocking bắt đầu scheduler và chặn luồng hiện tại.
	StartBlocking()

	// Stop dừng scheduler.
	Stop()

	// IsRunning kiểm tra xem scheduler có đang chạy không.
	IsRunning() bool

	// Clear xóa tất cả các công việc đã đăng ký.
	Clear()

	// GetScheduler trả về đối tượng scheduler gốc của gocron.
	GetScheduler() *gocron.Scheduler

	// RegisterEventListeners đăng ký các listener cho các sự kiện.
	RegisterEventListeners(eventListeners ...gocron.EventListener)
}

// manager triển khai interface Manager bằng cách nhúng gocron.Scheduler.
type manager struct {
	*gocron.Scheduler
}

// NewScheduler tạo một đối tượng Manager mới sử dụng gocron làm backend.
// Nhận tham số config để cấu hình scheduler.
func NewScheduler(cfg ...Config) Manager {
	scheduler := gocron.NewScheduler(time.Local)
	return &manager{
		Scheduler: scheduler,
	}
}

// NewSchedulerWithConfig tạo một đối tượng Manager mới với cấu hình cụ thể.
func NewSchedulerWithConfig(cfg Config) Manager {
	scheduler := gocron.NewScheduler(time.Local)
	return &manager{
		Scheduler: scheduler,
	}
}

// Every tạo một công việc mới với khoảng thời gian được chỉ định.
func (m *manager) Every(interval interface{}) Manager {
	m.Scheduler.Every(interval)
	return m
}

// Second chỉ định đơn vị thời gian là giây (đơn lẻ).
func (m *manager) Second() Manager {
	m.Scheduler.Second()
	return m
}

// Seconds chỉ định đơn vị thời gian là giây.
func (m *manager) Seconds() Manager {
	m.Scheduler.Seconds()
	return m
}

// Minutes chỉ định đơn vị thời gian là phút.
func (m *manager) Minutes() Manager {
	m.Scheduler.Minutes()
	return m
}

// Hours chỉ định đơn vị thời gian là giờ.
func (m *manager) Hours() Manager {
	m.Scheduler.Hours()
	return m
}

// Days chỉ định đơn vị thời gian là ngày.
func (m *manager) Days() Manager {
	m.Scheduler.Days()
	return m
}

// Weeks chỉ định đơn vị thời gian là tuần.
func (m *manager) Weeks() Manager {
	m.Scheduler.Weeks()
	return m
}

// At chỉ định thời điểm trong ngày để chạy công việc.
func (m *manager) At(time string) Manager {
	m.Scheduler.At(time)
	return m
}

// StartAt chỉ định thời điểm bắt đầu cho công việc.
func (m *manager) StartAt(startTime time.Time) Manager {
	m.Scheduler.StartAt(startTime)
	return m
}

// Cron thiết lập biểu thức cron cho công việc.
func (m *manager) Cron(cronExpression string) Manager {
	m.Scheduler.Cron(cronExpression)
	return m
}

// CronWithSeconds thiết lập biểu thức cron có hỗ trợ giây.
func (m *manager) CronWithSeconds(cronExpression string) Manager {
	m.Scheduler.CronWithSeconds(cronExpression)
	return m
}

// Tag đánh dấu công việc với các tag được chỉ định.
func (m *manager) Tag(tags ...string) Manager {
	m.Scheduler.Tag(tags...)
	return m
}

// SingletonMode đặt công việc ở chế độ singleton.
func (m *manager) SingletonMode() Manager {
	m.Scheduler.SingletonMode()
	return m
}

// Do đặt hàm để thực thi cho công việc.
func (m *manager) Do(jobFun interface{}, params ...interface{}) (*gocron.Job, error) {
	return m.Scheduler.Do(jobFun, params...)
}

// Name đặt tên cho công việc đang được cấu hình.
func (m *manager) Name(name string) Manager {
	m.Scheduler.Name(name)
	return m
}

// RemoveByTag xóa các công việc theo tag.
func (m *manager) RemoveByTag(tag string) error {
	return m.Scheduler.RemoveByTag(tag)
}

// RemoveByTags xóa các công việc khớp với TẤT CẢ tags đã chỉ định.
func (m *manager) RemoveByTags(tags ...string) error {
	return m.Scheduler.RemoveByTags(tags...)
}

// FindJobsByTag tìm công việc theo tag.
func (m *manager) FindJobsByTag(tags ...string) ([]*gocron.Job, error) {
	return m.Scheduler.FindJobsByTag(tags...)
}

// StartAsync bắt đầu scheduler trong một goroutine riêng.
func (m *manager) StartAsync() {
	m.Scheduler.StartAsync()
}

// StartBlocking bắt đầu scheduler và chặn luồng hiện tại.
func (m *manager) StartBlocking() {
	m.Scheduler.StartBlocking()
}

// Stop dừng scheduler.
func (m *manager) Stop() {
	m.Scheduler.Stop()
}

// Clear xóa tất cả các công việc đã đăng ký.
func (m *manager) Clear() {
	m.Scheduler.Clear()
}

// GetScheduler trả về đối tượng scheduler gốc của gocron.
func (m *manager) GetScheduler() *gocron.Scheduler {
	return m.Scheduler
}

// WithDistributedLocker thiết lập distributed locker cho scheduler.
func (m *manager) WithDistributedLocker(locker gocron.Locker) Manager {
	m.Scheduler.WithDistributedLocker(locker)
	return m
}

// RegisterEventListeners đăng ký các listener cho các sự kiện.
func (m *manager) RegisterEventListeners(eventListeners ...gocron.EventListener) {
	m.Scheduler.RegisterEventListeners(eventListeners...)
}

// IsRunning kiểm tra xem scheduler có đang chạy không.
func (m *manager) IsRunning() bool {
	return m.Scheduler.IsRunning()
}
