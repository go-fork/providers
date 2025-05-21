package log

import (
	"errors"
	"testing"

	"github.com/go-fork/providers/log/handler"
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
}

func TestDefaultManagerAddHandler(t *testing.T) {
	m := NewManager().(*DefaultManager)
	h := &MockHandler{}

	// Test thêm handler
	m.AddHandler("test", h)

	// Kiểm tra handler được thêm đúng cách
	if _, ok := m.handlers["test"]; !ok {
		t.Errorf("AddHandler không thêm handler vào map")
	}
}

func TestDefaultManagerRemoveHandler(t *testing.T) {
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

	// Xóa handler không tồn tại không gây lỗi
	m.RemoveHandler("nonexistent")
}

func TestSetMinLevel(t *testing.T) {
	m := NewManager().(*DefaultManager)

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
}

func TestLogMethods(t *testing.T) {
	m := NewManager()
	h := &MockHandler{}

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
			h.LogLevel = 0
			h.LogMessage = ""
			h.LogArgs = nil

			// Gọi method
			tt.logFunc(tt.message, tt.args...)

			// Kiểm tra handler được gọi với tham số đúng
			if !h.LogCalled {
				t.Errorf("%s không gọi handler", tt.name)
			}
			if h.LogLevel != tt.level {
				t.Errorf("%s truyền sai log level: got %v, want %v", tt.name, h.LogLevel, tt.level)
			}
			if h.LogMessage != tt.message {
				t.Errorf("%s truyền sai message: got %q, want %q", tt.name, h.LogMessage, tt.message)
			}

			// Kiểm tra args
			if len(h.LogArgs) != len(tt.args) {
				t.Errorf("%s truyền sai số lượng args: got %d, want %d",
					tt.name, len(h.LogArgs), len(tt.args))
			}
		})
	}
}

func TestDefaultManagerClose(t *testing.T) {
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
}

func TestDefaultManagerCloseWithError(t *testing.T) {
	m := NewManager()
	h1 := &MockHandler{}
	h2 := &MockHandler{ShouldError: true}

	m.AddHandler("h1", h1)
	m.AddHandler("h2", h2)

	// Test Close method với handler trả về lỗi
	err := m.Close()
	if err == nil {
		t.Error("Close với handler lỗi không trả về lỗi")
	}

	// Kiểm tra cả hai handlers được đóng dù có lỗi
	if !h1.CloseCalled {
		t.Error("Close không gọi Close() trên handler đầu tiên khi handler thứ hai lỗi")
	}
	if !h2.CloseCalled {
		t.Error("Close không gọi Close() trên handler lỗi")
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
