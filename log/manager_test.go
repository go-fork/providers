package log

import (
	"errors"
	"fmt"
	"testing"

	"go.fork.vn/providers/log/handler"
)

// MockHandler triển khai interface handler.Handler để kiểm tra
type MockHandler struct {
	LogCalled   bool
	CloseCalled bool
	ShouldError bool
	LogLevel    handler.Level
	LogMessage  string
	LogArgs     []interface{}
}

func (m *MockHandler) Log(level handler.Level, message string, args ...interface{}) error {
	m.LogCalled = true
	m.LogLevel = level
	m.LogMessage = message
	m.LogArgs = args
	if m.ShouldError {
		return errors.New("mock log error")
	}
	return nil
}

func (m *MockHandler) Close() error {
	m.CloseCalled = true
	if m.ShouldError {
		return errors.New("mock close error")
	}
	return nil
}

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager() trả về nil")
	}

	// Kiểm tra kiểu đúng
	_, ok := m.(*manager)
	if !ok {
		t.Errorf("NewManager() không trả về *manager, got %T", m)
	}

	// Kiểm tra các thuộc tính mặc định
	defaultManager := m.(*manager)
	if defaultManager.minLevel != handler.InfoLevel {
		t.Errorf("Manager mới không đặt minLevel mặc định là InfoLevel, got %v", defaultManager.minLevel)
	}
	if len(defaultManager.handlers) != 0 {
		t.Errorf("Manager mới không có handlers trống, got %d handlers", len(defaultManager.handlers))
	}
}

func TestManagerAddHandler(t *testing.T) {
	m := NewManager().(*manager)
	h := &MockHandler{}

	// Test thêm handler
	m.AddHandler("test", h)

	// Kiểm tra handler được thêm đúng cách
	if _, ok := m.handlers["test"]; !ok {
		t.Errorf("AddHandler không thêm handler vào map")
	}

	// Thêm handler thứ hai
	h2 := &MockHandler{}
	m.AddHandler("test2", h2)

	// Kiểm tra cả hai handler tồn tại
	if _, ok := m.handlers["test"]; !ok {
		t.Errorf("Handler đầu tiên không còn tồn tại sau khi thêm handler thứ hai")
	}
	if _, ok := m.handlers["test2"]; !ok {
		t.Errorf("Handler thứ hai không được thêm vào map")
	}

	// Ghi đè lên handler cũ - kiểm tra handler cũ được đóng
	h3 := &MockHandler{}
	m.AddHandler("test", h3)

	// Xác minh rằng h3 thay thế h trong map
	handler, ok := m.handlers["test"]
	if !ok {
		t.Error("Handler 'test' không tồn tại sau khi ghi đè")
	}
	if handler != h3 {
		t.Error("Handler không được ghi đè đúng cách")
	}
	// Kiểm tra handler cũ đã được đóng
	if !h.CloseCalled {
		t.Error("Handler cũ không được đóng khi bị ghi đè bởi AddHandler")
	}
}

func TestManagerRemoveHandler(t *testing.T) {
	m := NewManager()
	h := &MockHandler{}

	// Thêm handler
	m.AddHandler("test", h)

	// Xóa handler
	m.RemoveHandler("test")

	// Kiểm tra handler được gọi Close
	if !h.CloseCalled {
		t.Error("RemoveHandler không gọi Close() trên handler")
	}

	// Kiểm tra handler đã bị xóa khỏi map
	defaultManager := m.(*manager)
	if _, ok := defaultManager.handlers["test"]; ok {
		t.Error("RemoveHandler không xóa handler khỏi map")
	}

	// Xóa handler không tồn tại không gây lỗi
	m.RemoveHandler("nonexistent")

	// Thêm handler có lỗi khi close
	h2 := &MockHandler{ShouldError: true}
	m.AddHandler("test2", h2)

	// RemoveHandler vẫn nên xóa handler ngay cả khi Close() gây lỗi
	m.RemoveHandler("test2")

	// Kiểm tra handler đã bị xóa khỏi map dù Close() lỗi
	if _, ok := defaultManager.handlers["test2"]; ok {
		t.Error("RemoveHandler không xóa handler khỏi map khi Close() lỗi")
	}
}

func TestSetMinLevel(t *testing.T) {
	m := NewManager().(*manager)

	// Kiểm tra mức mặc định
	if m.minLevel != handler.InfoLevel {
		t.Errorf("Mức mặc định không phải InfoLevel, got %v", m.minLevel)
	}

	// Đặt mức mới
	m.SetMinLevel(handler.WarningLevel)

	// Kiểm tra mức được đặt đúng
	if m.minLevel != handler.WarningLevel {
		t.Errorf("SetMinLevel không đặt minLevel đúng, got %v, want %v",
			m.minLevel, handler.WarningLevel)
	}

	// Đặt mức thấp nhất
	m.SetMinLevel(handler.DebugLevel)
	if m.minLevel != handler.DebugLevel {
		t.Errorf("SetMinLevel không đặt minLevel thành DebugLevel, got %v", m.minLevel)
	}

	// Đặt mức cao nhất
	m.SetMinLevel(handler.FatalLevel)
	if m.minLevel != handler.FatalLevel {
		t.Errorf("SetMinLevel không đặt minLevel thành FatalLevel, got %v", m.minLevel)
	}
}

func TestLogMethods(t *testing.T) {
	m := NewManager()
	h := &MockHandler{}

	// Set min level to DebugLevel to ensure Debug logs are processed
	m.SetMinLevel(handler.DebugLevel)
	m.AddHandler("test", h)

	tests := []struct {
		name    string
		logFunc func(message string, args ...interface{})
		level   handler.Level
		message string
		args    []interface{}
	}{
		{"Debug", m.Debug, handler.DebugLevel, "debug message", []interface{}{1, 2}},
		{"Info", m.Info, handler.InfoLevel, "info message", []interface{}{3, 4}},
		{"Warning", m.Warning, handler.WarningLevel, "warn message", []interface{}{5, 6}},
		{"Error", m.Error, handler.ErrorLevel, "error message", []interface{}{7, 8}},
		{"Fatal", m.Fatal, handler.FatalLevel, "fatal message", []interface{}{9, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.LogCalled = false
			h.LogMessage = ""

			// Gọi method với message có định dạng
			formattedMsg := fmt.Sprintf("%s %%d %%d", tt.message)
			tt.logFunc(formattedMsg, tt.args...)

			// Kiểm tra handler được gọi với tham số đúng
			if !h.LogCalled {
				t.Errorf("%s không gọi handler", tt.name)
			}
			if h.LogLevel != tt.level {
				t.Errorf("%s truyền sai log level: got %v, want %v", tt.name, h.LogLevel, tt.level)
			}

			// Kiểm tra message được định dạng đúng
			expectedMsg := fmt.Sprintf(formattedMsg, tt.args...)
			if h.LogMessage != expectedMsg {
				t.Errorf("%s không định dạng message đúng: got %q, want %q",
					tt.name, h.LogMessage, expectedMsg)
			}
		})
	}
}

func TestManagerClose(t *testing.T) {
	m := NewManager()
	h1 := &MockHandler{}
	h2 := &MockHandler{}

	m.AddHandler("h1", h1)
	m.AddHandler("h2", h2)

	// Test Close method
	err := m.Close()
	if err != nil {
		t.Errorf("Close trả về lỗi: %v", err)
	}

	// Kiểm tra cả hai handlers được đóng
	if !h1.CloseCalled {
		t.Error("Close không gọi Close() trên handler đầu tiên")
	}
	if !h2.CloseCalled {
		t.Error("Close không gọi Close() trên handler thứ hai")
	}

	// Lưu ý: Theo hiện thực hiện tại, Close() không xóa các handlers khỏi map
	// Nó chỉ đóng các handlers nhưng vẫn giữ chúng trong map
	defaultManager := m.(*manager)
	if len(defaultManager.handlers) == 0 {
		t.Error("Không mong đợi map handlers trống sau khi close, hiện thực chỉ đóng handlers")
	}
}

func TestManagerCloseWithError(t *testing.T) {
	m := NewManager()
	h1 := &MockHandler{}
	h2 := &MockHandler{ShouldError: true}
	h3 := &MockHandler{}

	m.AddHandler("h1", h1)
	m.AddHandler("h2", h2)
	m.AddHandler("h3", h3)

	// Test Close method với handler trả về lỗi
	err := m.Close()
	if err == nil {
		t.Error("Close với handler lỗi không trả về lỗi")
	}

	// Kiểm tra tất cả handlers được đóng dù có lỗi
	if !h1.CloseCalled || !h2.CloseCalled || !h3.CloseCalled {
		t.Error("Close không gọi Close() trên tất cả handlers")
	}
}

func TestLogFilteringByMinLevel(t *testing.T) {
	m := NewManager()
	h := &MockHandler{}

	m.AddHandler("test", h)

	// Đặt mức tối thiểu là ErrorLevel
	m.SetMinLevel(handler.ErrorLevel)

	// Log ở mức thấp hơn không nên gọi handler
	m.Debug("debug message")
	if h.LogCalled {
		t.Error("Debug log không bị lọc khi dưới ngưỡng")
	}

	m.Info("info message")
	if h.LogCalled {
		t.Error("Info log không bị lọc khi dưới ngưỡng")
	}

	m.Warning("warning message")
	if h.LogCalled {
		t.Error("Warning log không bị lọc khi dưới ngưỡng")
	}

	// Reset
	h.LogCalled = false

	// Log ở mức bằng hoặc cao hơn nên gọi handler
	m.Error("error message")
	if !h.LogCalled {
		t.Error("Error log bị lọc khi bằng hoặc trên ngưỡng")
	}

	// Reset
	h.LogCalled = false

	m.Fatal("fatal message")
	if !h.LogCalled {
		t.Error("Fatal log bị lọc khi bằng hoặc trên ngưỡng")
	}
}

// TestLogWithFormatting kiểm tra định dạng thông điệp
func TestLogWithFormatting(t *testing.T) {
	m := NewManager()
	h := &MockHandler{}

	m.AddHandler("test", h)

	// Kiểm tra log với định dạng khác nhau
	tests := []struct {
		name   string
		format string
		args   []interface{}
		want   string
	}{
		{"String", "Hello %s", []interface{}{"world"}, "Hello world"},
		{"Number", "Number: %d", []interface{}{42}, "Number: 42"},
		{"Multiple", "%s: %d, %f", []interface{}{"Test", 123, 45.67}, "Test: 123, 45.670000"},
		{"Empty", "No args", nil, "No args"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.LogCalled = false
			h.LogMessage = ""

			// Gọi phương thức Info với định dạng
			m.Info(tt.format, tt.args...)

			// Kiểm tra định dạng đúng
			if !h.LogCalled {
				t.Error("Log không gọi handler")
			}
			if h.LogMessage != tt.want {
				t.Errorf("Định dạng không đúng: got %q, want %q", h.LogMessage, tt.want)
			}
		})
	}
}

func TestManagerGetHandler(t *testing.T) {
	m := NewManager()
	h1 := &MockHandler{}
	h2 := &MockHandler{}

	// Thêm các handlers
	m.AddHandler("handler1", h1)
	m.AddHandler("handler2", h2)

	// Lấy handler đã đăng ký
	handlerResult := m.GetHandler("handler1")
	if handlerResult == nil {
		t.Error("GetHandler trả về nil cho handler đã đăng ký")
	}

	if handlerResult != h1 {
		t.Errorf("GetHandler không trả về đúng handler, got %v, want %v", handlerResult, h1)
	}

	// Lấy handler thứ hai
	handlerResult = m.GetHandler("handler2")
	if handlerResult != h2 {
		t.Errorf("GetHandler không trả về đúng handler cho key thứ hai, got %v, want %v", handlerResult, h2)
	}

	// Lấy handler không tồn tại
	handlerResult = m.GetHandler("nonexistent")
	if handlerResult != nil {
		t.Errorf("GetHandler không trả về nil cho handler không tồn tại, got %v", handlerResult)
	}

	// Kiểm tra thread-safety (thông qua kiểm tra chức năng cơ bản)
	// Thêm handler mới trong khi đang thực hiện GetHandler
	go func() {
		m.AddHandler("handler3", &MockHandler{})
	}()

	// Xóa handler trong khi đang thực hiện GetHandler
	go func() {
		m.RemoveHandler("handler1")
	}()

	// Lấy handler2 một lần nữa sau các thao tác đồng thời
	handlerResult = m.GetHandler("handler2")
	if handlerResult != h2 {
		t.Errorf("GetHandler không hoạt động đúng sau các thao tác đồng thời, got %v, want %v", handlerResult, h2)
	}
}

func TestLogWithErrorHandler(t *testing.T) {
	m := NewManager()
	h := &MockHandler{ShouldError: true}

	m.AddHandler("test", h)

	// Log với handler trả về lỗi
	// Không cần kiểm tra gì vì lỗi chỉ được in ra stderr
	// Nhưng log call vẫn hoàn thành không bị panic
	m.Info("this will error")

	// Kiểm tra handler vẫn được gọi dù trả về lỗi
	if !h.LogCalled {
		t.Error("Log không gọi handler khi biết handler sẽ lỗi")
	}
}
