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
	// Tạo mock handlers
	handler1 := &MockTestHandler{}
	handler2 := &MockTestHandler{}

	// Test với không có handler
	stack1 := NewStackHandler()
	if stack1 == nil {
		t.Fatal("NewStackHandler() không handler trả về nil")
	}
	if len(stack1.handlers) != 0 {
		t.Errorf("NewStackHandler() không handler nên có slice handlers trống, got length = %d",
			len(stack1.handlers))
	}

	// Test với một handler
	stack2 := NewStackHandler(handler1)
	if stack2 == nil {
		t.Fatal("NewStackHandler() với một handler trả về nil")
	}
	if len(stack2.handlers) != 1 {
		t.Errorf("NewStackHandler() với một handler nên có len(handlers) = 1, got = %d",
			len(stack2.handlers))
	}

	// Test với nhiều handlers
	stack3 := NewStackHandler(handler1, handler2)
	if stack3 == nil {
		t.Fatal("NewStackHandler() với nhiều handlers trả về nil")
	}
	if len(stack3.handlers) != 2 {
		t.Errorf("NewStackHandler() với hai handlers nên có len(handlers) = 2, got = %d",
			len(stack3.handlers))
	}
}

func TestStackHandlerLog(t *testing.T) {
	// Test với không handler
	stack1 := NewStackHandler()
	err := stack1.Log(InfoLevel, "test message")
	if err != nil {
		t.Errorf("StackHandler.Log() với không handler nên không lỗi, got = %v", err)
	}

	// Test với một handler thành công
	handler1 := &MockTestHandler{}
	stack2 := NewStackHandler(handler1)
	err = stack2.Log(ErrorLevel, "error message")
	if err != nil {
		t.Errorf("StackHandler.Log() với handler thành công nên không lỗi, got = %v", err)
	}
	if !handler1.LogCalled {
		t.Error("StackHandler.Log() không gọi handler.Log()")
	}
	if handler1.LogLevel != ErrorLevel {
		t.Errorf("StackHandler.Log() không truyền đúng level, got = %v, want = %v",
			handler1.LogLevel, ErrorLevel)
	}
	if handler1.LogMessage != "error message" {
		t.Errorf("StackHandler.Log() không truyền đúng message, got = %v, want = %v",
			handler1.LogMessage, "error message")
	}

	// Test với một handler lỗi
	handler2 := &MockTestHandler{ShouldError: true}
	stack3 := NewStackHandler(handler2)
	err = stack3.Log(InfoLevel, "test message")
	if err == nil {
		t.Error("StackHandler.Log() với handler lỗi nên trả về lỗi")
	}

	// Test với nhiều handlers, một lỗi một thành công
	handler3 := &MockTestHandler{}
	handler4 := &MockTestHandler{ShouldError: true}
	stack4 := NewStackHandler(handler3, handler4)
	err = stack4.Log(WarningLevel, "warning message")
	if err == nil {
		t.Error("StackHandler.Log() với một handler lỗi nên trả về lỗi")
	}
	if !handler3.LogCalled || !handler4.LogCalled {
		t.Error("StackHandler.Log() không gọi tất cả các handlers khi có lỗi")
	}
}

func TestStackHandlerClose(t *testing.T) {
	// Test với không handler
	stack1 := NewStackHandler()
	err := stack1.Close()
	if err != nil {
		t.Errorf("StackHandler.Close() với không handler nên không lỗi, got = %v", err)
	}

	// Test với một handler thành công
	handler1 := &MockTestHandler{}
	stack2 := NewStackHandler(handler1)
	err = stack2.Close()
	if err != nil {
		t.Errorf("StackHandler.Close() với handler thành công nên không lỗi, got = %v", err)
	}
	if !handler1.CloseCalled {
		t.Error("StackHandler.Close() không gọi handler.Close()")
	}

	// Test với một handler lỗi
	handler2 := &MockTestHandler{ShouldError: true}
	stack3 := NewStackHandler(handler2)
	err = stack3.Close()
	if err == nil {
		t.Error("StackHandler.Close() với handler lỗi nên trả về lỗi")
	}

	// Test với nhiều handlers, một lỗi một thành công
	handler3 := &MockTestHandler{}
	handler4 := &MockTestHandler{ShouldError: true}
	stack4 := NewStackHandler(handler3, handler4)
	err = stack4.Close()
	if err == nil {
		t.Error("StackHandler.Close() với một handler lỗi nên trả về lỗi")
	}
	if !handler3.CloseCalled || !handler4.CloseCalled {
		t.Error("StackHandler.Close() không gọi tất cả các handlers khi có lỗi")
	}
}

func TestStackHandlerAddHandler(t *testing.T) {
	// Tạo stack handler ban đầu
	stack := NewStackHandler()

	// Thêm một handler
	handler1 := &MockTestHandler{}
	stack.AddHandler(handler1)
	if len(stack.handlers) != 1 {
		t.Errorf("StackHandler.AddHandler() sau thêm 1 handler nên có len(handlers) = 1, got = %d",
			len(stack.handlers))
	}

	// Thêm handler thứ hai
	handler2 := &MockTestHandler{}
	stack.AddHandler(handler2)
	if len(stack.handlers) != 2 {
		t.Errorf("StackHandler.AddHandler() sau thêm 2 handlers nên có len(handlers) = 2, got = %d",
			len(stack.handlers))
	}

	// Kiểm tra cả hai handlers có hoạt động không
	err := stack.Log(InfoLevel, "test after add")
	if err != nil {
		t.Errorf("StackHandler.Log() sau khi thêm handlers nên không lỗi, got = %v", err)
	}
	if !handler1.LogCalled || !handler2.LogCalled {
		t.Error("StackHandler.Log() không gọi tất cả các handlers sau khi thêm")
	}
}
