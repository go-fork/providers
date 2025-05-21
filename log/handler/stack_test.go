package handler

import (
	"errors"
	"testing"
)

// MockTestHandler giả lập một Handler để kiểm tra StackHandler
type MockTestHandler struct {
	LogCalled   bool
	CloseCalled bool
	ShouldError bool
	LogLevel    Level
	LogMessage  string
	LogArgs     []interface{}
}

func (m *MockTestHandler) Log(level Level, message string, args ...interface{}) error {
	m.LogCalled = true
	m.LogLevel = level
	m.LogMessage = message
	m.LogArgs = args
	if m.ShouldError {
		return errors.New("mock log error")
	}
	return nil
}

func (m *MockTestHandler) Close() error {
	m.CloseCalled = true
	if m.ShouldError {
		return errors.New("mock close error")
	}
	return nil
}

func TestNewStackHandler(t *testing.T) {
	// Test tạo stack handler không có handler con
	stack1 := NewStackHandler()
	if stack1 == nil {
		t.Fatal("NewStackHandler() trả về nil")
	}
	if len(stack1.handlers) != 0 {
		t.Errorf("NewStackHandler() không khởi tạo đúng slice handlers, got len %d",
			len(stack1.handlers))
	}

	// Test tạo stack handler với một handler con
	h1 := &MockTestHandler{}
	stack2 := NewStackHandler(h1)
	if stack2 == nil {
		t.Fatal("NewStackHandler(h1) trả về nil")
	}
	if len(stack2.handlers) != 1 {
		t.Errorf("NewStackHandler(h1) không khởi tạo đúng slice handlers, got len %d",
			len(stack2.handlers))
	}

	// Test tạo stack handler với nhiều handler con
	h2 := &MockTestHandler{}
	h3 := &MockTestHandler{}
	stack3 := NewStackHandler(h1, h2, h3)
	if stack3 == nil {
		t.Fatal("NewStackHandler(h1, h2, h3) trả về nil")
	}
	if len(stack3.handlers) != 3 {
		t.Errorf("NewStackHandler(h1, h2, h3) không khởi tạo đúng slice handlers, got len %d",
			len(stack3.handlers))
	}
}

func TestStackHandlerLogWithError(t *testing.T) {
	// Tạo các handler giả lập, h2 sẽ trả về lỗi
	h1 := &MockTestHandler{}
	h2 := &MockTestHandler{ShouldError: true}
	h3 := &MockTestHandler{}

	// Tạo stack handler với các handler con
	stack := NewStackHandler(h1, h2, h3)

	// Test Log method với một handler con lỗi
	err := stack.Log(ErrorLevel, "test message")
	if err == nil {
		t.Error("Log() không trả về lỗi khi một handler lỗi")
	}

	// Kiểm tra tất cả handler con được gọi dù có lỗi
	if !h1.LogCalled {
		t.Error("Log() không gọi Log() trên handler con thứ nhất khi có lỗi")
	}
	if !h2.LogCalled {
		t.Error("Log() không gọi Log() trên handler con có lỗi")
	}
	if !h3.LogCalled {
		t.Error("Log() không gọi Log() trên handler con thứ ba sau khi gặp lỗi")
	}
}

func TestStackHandlerClose(t *testing.T) {
	// Tạo các handler giả lập
	h1 := &MockTestHandler{}
	h2 := &MockTestHandler{}

	// Tạo stack handler với các handler con
	stack := NewStackHandler(h1, h2)

	// Test Close method
	err := stack.Close()
	if err != nil {
		t.Errorf("Close() trả về lỗi: %v", err)
	}

	// Kiểm tra cả hai handler con được đóng
	if !h1.CloseCalled {
		t.Error("Close() không gọi Close() trên handler con thứ nhất")
	}
	if !h2.CloseCalled {
		t.Error("Close() không gọi Close() trên handler con thứ hai")
	}
}

func TestStackHandlerAddHandler(t *testing.T) {
	// Tạo stack handler trống
	stack := NewStackHandler()

	// Kiểm tra ban đầu không có handler
	if len(stack.handlers) != 0 {
		t.Errorf("Stack handler mới không rỗng, got len %d", len(stack.handlers))
	}

	// Thêm một handler
	h1 := &MockTestHandler{}
	stack.AddHandler(h1)

	// Kiểm tra handler được thêm
	if len(stack.handlers) != 1 {
		t.Errorf("AddHandler không thêm đúng, got len %d, want 1", len(stack.handlers))
	}

	// Thêm một handler nữa
	h2 := &MockTestHandler{}
	stack.AddHandler(h2)

	// Kiểm tra handler thứ hai được thêm
	if len(stack.handlers) != 2 {
		t.Errorf("AddHandler lần thứ hai không thêm đúng, got len %d, want 2",
			len(stack.handlers))
	}

	// Kiểm tra thứ tự handlers
	if stack.handlers[0] != h1 || stack.handlers[1] != h2 {
		t.Error("AddHandler không duy trì thứ tự handlers")
	}
}
