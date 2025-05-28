package mailer

import (
	"time"

	"github.com/go-fork/providers/queue"
)

// Manager định nghĩa interface cho việc quản lý các thành phần mailer.
type Manager interface {
	// Mailer trả về instance của Mailer.
	Mailer() Mailer

	// Config trả về cấu hình hiện tại của mailer.
	Config() *Config

	// QueueSettings trả về cấu hình queue cho mailer.
	QueueSettings() *QueueSettings

	// QueueEnabled kiểm tra xem chức năng queue có được bật không.
	QueueEnabled() bool

	// QueueManager trả về queue Manager instance nếu queue được bật.
	QueueManager() queue.Manager

	// QueueClient trả về queue Client instance nếu queue được bật.
	QueueClient() queue.Client

	// EnqueueMessage đưa một message vào queue để gửi không đồng bộ.
	EnqueueMessage(message *Message) error

	// ProcessMessage xử lý một message từ queue và gửi mail.
	ProcessMessage(messageData []byte) error
}

// manager quản lý các thành phần của mailer
type manager struct {
	config        *Config
	queueSettings *QueueSettings
	mailer        Mailer
	queueManager  queue.Manager
}

// NewManager tạo một manager mới với cấu hình.
//
// Tham số:
//   - cfg: *Config - cấu hình cho mailer
//
// Trả về:
//   - Manager: một manager instance
//   - error: lỗi nếu có trong quá trình khởi tạo
func NewManager(cfg *Config) (Manager, error) {
	// Đảm bảo luôn có cấu hình
	if cfg == nil {
		cfg = NewConfig()
	}

	return NewManagerWithConfig(cfg)
}

// NewManagerWithConfig tạo một manager mới với cấu hình tùy chỉnh
func NewManagerWithConfig(cfg *Config) (Manager, error) {
	queueSettings := queueSettingsFromConfig(cfg.Queue)

	m := &manager{
		config:        cfg,
		queueSettings: queueSettings,
	}

	// Khởi tạo mailer
	m.mailer = NewSMTPMailer(cfg)

	// Khởi tạo queue manager nếu cần
	if queueSettings.QueueEnabled {
		queueCfg := queue.DefaultConfig()
		// Cấu hình queue dựa vào queueSettings
		queueCfg.Adapter.Default = queueSettings.QueueAdapter
		queueCfg.Server.DefaultQueue = queueSettings.QueueName
		queueCfg.Server.Concurrency = queueSettings.WorkerConcurrency
		queueCfg.Server.PollingInterval = int(queueSettings.PollingInterval / time.Millisecond)
		queueCfg.Server.RetryLimit = queueSettings.MaxRetries

		if queueSettings.QueueAdapter == "redis" {
			// Cấu hình Redis cho API mới
			queueCfg.Adapter.Redis.Prefix = queueSettings.QueuePrefix
			queueCfg.Adapter.Redis.ProviderKey = "redis"
		}

		m.queueManager = queue.NewManager(queueCfg)
	}

	return m, nil
}

// Mailer trả về instance của Mailer
func (m *manager) Mailer() Mailer {
	return m.mailer
}

// Config trả về cấu hình hiện tại của mailer
func (m *manager) Config() *Config {
	return m.config
}

// QueueSettings trả về cấu hình queue cho mailer
func (m *manager) QueueSettings() *QueueSettings {
	return m.queueSettings
}

// QueueEnabled kiểm tra xem chức năng queue có được bật không
func (m *manager) QueueEnabled() bool {
	return m.queueSettings != nil && m.queueSettings.QueueEnabled
}

// QueueManager trả về queue Manager instance nếu queue được bật
func (m *manager) QueueManager() queue.Manager {
	return m.queueManager
}

// QueueClient trả về queue Client instance nếu queue được bật
func (m *manager) QueueClient() queue.Client {
	if !m.QueueEnabled() || m.queueManager == nil {
		return nil
	}
	return m.queueManager.Client()
}

// EnqueueMessage đưa một message vào queue để gửi không đồng bộ
func (m *manager) EnqueueMessage(message *Message) error {
	if !m.QueueEnabled() {
		// Nếu queue không được bật, gửi trực tiếp
		return m.mailer.Send(message)
	}

	// Chuyển đổi message thành dữ liệu để lưu vào queue
	data, err := message.MarshalJSON()
	if err != nil {
		return err
	}

	// Lấy client và đưa task vào queue
	client := m.QueueClient()
	if client == nil {
		return m.mailer.Send(message)
	}

	// Sử dụng queue Options thay vì TaskOptions
	queueOpts := []queue.Option{
		queue.WithQueue(m.queueSettings.QueueName),
		queue.WithMaxRetry(m.queueSettings.MaxRetries),
		queue.WithTimeout(m.queueSettings.ProcessTimeout),
	}

	_, err = client.Enqueue("mailer:send", data, queueOpts...)
	return err
}

// ProcessMessage xử lý một message từ queue và gửi mail
func (m *manager) ProcessMessage(messageData []byte) error {
	// Tạo message từ dữ liệu
	message := NewMessage()
	err := message.UnmarshalJSON(messageData)
	if err != nil {
		return err
	}

	// Gửi mail
	return m.mailer.Send(message)
}
