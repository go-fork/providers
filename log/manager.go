package log

import (
	"fmt"
	"sync"

	"github.com/go-fork/providers/log/handler"
)

// Manager định nghĩa interface cho hệ thống logging tập trung.
//
// Interface Manager cung cấp các method để ghi log ở nhiều cấp độ nghiêm trọng
// khác nhau và để quản lý các handler.
//
// Các triển khai của interface này cần đảm bảo thread-safe và xử lý
// việc phân phối các log entry đến tất cả các handler đã đăng ký.
type Manager interface {
	// Debug ghi một thông điệp ở cấp độ debug.
	//
	// Tham số:
	//   - message: string - thông điệp log (có thể là chuỗi định dạng)
	//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
	Debug(message string, args ...interface{})

	// Info ghi một thông điệp ở cấp độ info.
	//
	// Tham số:
	//   - message: string - thông điệp log (có thể là chuỗi định dạng)
	//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
	Info(message string, args ...interface{})

	// Warning ghi một thông điệp ở cấp độ warning.
	//
	// Tham số:
	//   - message: string - thông điệp log (có thể là chuỗi định dạng)
	//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
	Warning(message string, args ...interface{})

	// Error ghi một thông điệp ở cấp độ error.
	//
	// Tham số:
	//   - message: string - thông điệp log (có thể là chuỗi định dạng)
	//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
	Error(message string, args ...interface{})

	// Fatal ghi một thông điệp ở cấp độ fatal.
	//
	// Tham số:
	//   - message: string - thông điệp log (có thể là chuỗi định dạng)
	//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
	Fatal(message string, args ...interface{})

	// AddHandler đăng ký một handler mới vào manager.
	//
	// Tham số:
	//   - name: string - định danh duy nhất cho handler
	//   - handler: handler.Handler - instance của handler cần thêm
	AddHandler(name string, handler handler.Handler)

	// RemoveHandler hủy đăng ký và đóng một handler.
	//
	// Tham số:
	//   - name: string - định danh của handler cần xóa
	RemoveHandler(name string)

	// SetMinLevel thiết lập ngưỡng cấp độ log tối thiểu.
	//
	// Tham số:
	//   - level: handler.Level - cấp độ tối thiểu để log
	SetMinLevel(level handler.Level)

	// Close đóng manager và tất cả các handler.
	//
	// Trả về:
	//   - error: một lỗi nếu việc đóng handler thất bại
	Close() error
}

// DefaultManager là triển khai chuẩn của interface Manager.
//
// DefaultManager cung cấp cách quản lý nhiều handler log với bộ lọc
// dựa trên cấp độ thread-safe. Nó được thiết kế cho truy cập đồng thời
// trong môi trường đa goroutine.
//
// Tính năng:
//   - Quản lý handler thread-safe bằng RWMutex
//   - Lọc cấp độ log
//   - Thêm/xóa handler động
//   - Dọn dẹp tài nguyên an toàn khi tắt
type DefaultManager struct {
	handlers map[string]handler.Handler // Map các handler theo tên
	minLevel handler.Level              // Ngưỡng cấp độ log tối thiểu
	mu       sync.RWMutex               // Mutex để đảm bảo thread-safety
}

// NewManager tạo và trả về một instance log manager mới.
//
// Hàm này khởi tạo một DefaultManager không có handler nào và InfoLevel là
// cấp độ log tối thiểu mặc định.
//
// Trả về:
//   - Manager: một instance mới của DefaultManager triển khai interface Manager.
//
// Ví dụ:
//
//	manager := log.NewManager()
//	manager.AddHandler("console", handler.NewConsoleHandler(true))
func NewManager() Manager {
	return &DefaultManager{
		handlers: make(map[string]handler.Handler),
		minLevel: handler.InfoLevel, // Mặc định là InfoLevel
	}
}

// Debug ghi một thông điệp ở cấp độ debug.
//
// Debug logs dành cho thông tin chẩn đoán chi tiết hữu ích trong quá trình
// phát triển hoặc khắc phục sự cố.
//
// Tham số:
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
//
// Ví dụ:
//
//	manager.Debug("Lần thử kết nối %d đến %s", attempt, serverAddress)
func (m *DefaultManager) Debug(message string, args ...interface{}) {
	// Ghi log ở cấp độ DebugLevel
	m.log(handler.DebugLevel, message, args...)
}

// Info ghi một thông điệp ở cấp độ info.
//
// Info logs dành cho thông tin hoạt động chung về hành vi
// bình thường của ứng dụng.
//
// Tham số:
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
//
// Ví dụ:
//
//	manager.Info("Máy chủ đã khởi động trên cổng %d", port)
func (m *DefaultManager) Info(message string, args ...interface{}) {
	// Ghi log ở cấp độ InfoLevel
	m.log(handler.InfoLevel, message, args...)
}

// Warning ghi một thông điệp ở cấp độ warning.
//
// Warning logs chỉ ra các vấn đề tiềm ẩn hoặc điều kiện không mong đợi
// mà không phải lỗi nhưng có thể cần chú ý.
//
// Tham số:
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
//
// Ví dụ:
//
//	manager.Warning("Sử dụng bộ nhớ cao: %d MB", memoryUsage)
func (m *DefaultManager) Warning(message string, args ...interface{}) {
	// Ghi log ở cấp độ WarningLevel
	m.log(handler.WarningLevel, message, args...)
}

// Error ghi một thông điệp ở cấp độ error.
//
// Error logs chỉ ra các lỗi hoặc thất bại ảnh hưởng đến hoạt động bình thường
// nhưng không yêu cầu phải kết thúc ngay lập tức.
//
// Tham số:
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
//
// Ví dụ:
//
//	manager.Error("Xử lý yêu cầu thất bại: %v", err)
func (m *DefaultManager) Error(message string, args ...interface{}) {
	// Ghi log ở cấp độ ErrorLevel
	m.log(handler.ErrorLevel, message, args...)
}

// Fatal ghi một thông điệp ở cấp độ fatal.
//
// Fatal logs chỉ ra các lỗi nghiêm trọng thường yêu cầu kết thúc ứng dụng
// hoặc cần sự chú ý ngay lập tức của người quản trị.
//
// Tham số:
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - các tham số tùy chọn để định dạng thông điệp
//
// Ví dụ:
//
//	manager.Fatal("Kết nối database thất bại: %v", err)
func (m *DefaultManager) Fatal(message string, args ...interface{}) {
	// Ghi log ở cấp độ FatalLevel
	m.log(handler.FatalLevel, message, args...)
}

// AddHandler thêm một handler log mới vào manager.
//
// Method này đăng ký một handler với tên đã cho. Nếu một handler với cùng tên
// đã tồn tại, nó sẽ bị thay thế mà không đóng handler cũ. Method này là thread-safe.
//
// Tham số:
//   - name: string - tên duy nhất cho handler
//   - handler: handler.Handler - triển khai handler cần thêm
//
// Ví dụ:
//
//	// Thêm một file handler
//	fileHandler, _ := handler.NewFileHandler("app.log", 10*1024*1024)
//	manager.AddHandler("file", fileHandler)
func (m *DefaultManager) AddHandler(name string, handler handler.Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Thêm hoặc thay thế handler theo tên
	m.handlers[name] = handler
}

// RemoveHandler xóa một handler khỏi manager theo tên.
//
// Handler sẽ được đóng đúng cách trước khi xóa để đảm bảo tất cả các tài nguyên
// được giải phóng. Method này là thread-safe.
//
// Tham số:
//   - name: string - tên của handler cần xóa
//
// Nếu handler được chỉ định không tồn tại, thao tác này không làm gì.
//
// Ví dụ:
//
//	manager.RemoveHandler("file") // Xóa và đóng handler "file"
func (m *DefaultManager) RemoveHandler(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Đóng và xóa handler nếu nó tồn tại
	if handler, ok := m.handlers[name]; ok {
		handler.Close()
		delete(m.handlers, name)
	}
}

// SetMinLevel thiết lập cấp độ log tối thiểu cho manager.
//
// Bất kỳ log entry nào có cấp độ dưới ngưỡng này sẽ bị bỏ qua.
// Method này là thread-safe.
//
// Tham số:
//   - level: handler.Level - cấp độ log tối thiểu cần thiết lập
//
// Ví dụ:
//
//	// Chỉ xử lý log Warning, Error và Fatal
//	manager.SetMinLevel(handler.WarningLevel)
func (m *DefaultManager) SetMinLevel(level handler.Level) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.minLevel = level
}

// Close đóng tất cả các handler log đã đăng ký và giải phóng tài nguyên của chúng.
//
// Method này nên được gọi khi ứng dụng đang đóng để đảm bảo
// tất cả các file log được đóng đúng cách và tài nguyên được giải phóng.
//
// Trả về:
//   - error: lỗi đầu tiên gặp phải khi đóng handler, hoặc nil nếu tất cả đều đóng thành công
//
// Ví dụ:
//
//	if err := manager.Close(); err != nil {
//	    fmt.Fprintf(os.Stderr, "Lỗi khi đóng log manager: %v\n", err)
//	}
func (m *DefaultManager) Close() error {
	// Tạo một bản sao của map handlers để giảm thiểu thời gian giữ lock
	m.mu.Lock()
	handlersCopy := make(map[string]handler.Handler, len(m.handlers))
	for k, v := range m.handlers {
		handlersCopy[k] = v
	}
	m.mu.Unlock()

	// Đóng từng handler, theo dõi lỗi đầu tiên
	var firstErr error
	for name, handler := range handlersCopy {
		if err := handler.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed to close handler %s: %w", name, err)
		}
	}
	return firstErr
}

// log là method nội bộ để ghi một log entry đến tất cả các handler.
//
// Method này xử lý lọc cấp độ, định dạng thông điệp và gửi
// log entry đến tất cả các handler đã đăng ký. Nó được thiết kế để giảm thiểu
// thời gian giữ lock để tăng concurrency.
//
// Tham số:
//   - level: handler.Level - cấp độ log của thông điệp
//   - message: string - thông điệp log (có thể là chuỗi định dạng)
//   - args: ...interface{} - tham số tùy chọn để định dạng thông điệp
func (m *DefaultManager) log(level handler.Level, message string, args ...interface{}) {
	// Bỏ qua nếu dưới cấp độ tối thiểu
	if level < m.minLevel {
		return
	}

	// Lấy snapshot của handlers để giảm thiểu thời gian giữ lock
	m.mu.RLock()
	handlersCopy := make(map[string]handler.Handler, len(m.handlers))
	for k, v := range m.handlers {
		handlersCopy[k] = v
	}
	m.mu.RUnlock()

	// Định dạng thông điệp nếu có tham số
	formattedMessage := message
	if len(args) > 0 {
		formattedMessage = fmt.Sprintf(message, args...)
	}

	// Ghi log entry đến tất cả các handler
	for name, handler := range handlersCopy {
		if err := handler.Log(level, formattedMessage); err != nil {
			// Xử lý lỗi logging (ghi ra stderr)
			fmt.Printf("Lỗi khi ghi log đến handler %s: %v\n", name, err)
		}
	}
}
