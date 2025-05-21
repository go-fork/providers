package handler

import (
	"testing"
)

func TestNewConsoleHandler(t *testing.T) {
	// Test với colored = true
	h1 := NewConsoleHandler(true)
	if h1 == nil {
		t.Fatal("NewConsoleHandler(true) trả về nil")
	}
	if !h1.colored {
		t.Error("NewConsoleHandler(true) không đặt colored đúng")
	}

	// Test với colored = false
	h2 := NewConsoleHandler(false)
	if h2 == nil {
		t.Fatal("NewConsoleHandler(false) trả về nil")
	}
	if h2.colored {
		t.Error("NewConsoleHandler(false) không đặt colored đúng")
	}
}

func TestConsoleHandlerLog(t *testing.T) {
	// Test non-colored console handler
	h := &ConsoleHandler{colored: false}

	// Ghi log ở mỗi level
	err1 := h.Log(DebugLevel, "debug message")
	err2 := h.Log(InfoLevel, "info message")
	err3 := h.Log(WarningLevel, "warning message")
	err4 := h.Log(ErrorLevel, "error message")
	err5 := h.Log(FatalLevel, "fatal message")

	// Kiểm tra không có lỗi
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		t.Errorf("ConsoleHandler.Log() trả về lỗi: %v, %v, %v, %v, %v",
			err1, err2, err3, err4, err5)
	}

	// Test với tham số định dạng
	err := h.Log(InfoLevel, "formatted %s %d", "message", 123)
	if err != nil {
		t.Errorf("ConsoleHandler.Log() với định dạng trả về lỗi: %v", err)
	}
}

func TestColoredConsoleHandler(t *testing.T) {
	// Test colored console handler
	h := &ConsoleHandler{colored: true}

	// Ghi log ở mỗi level
	err1 := h.Log(DebugLevel, "debug message")
	err2 := h.Log(InfoLevel, "info message")
	err3 := h.Log(WarningLevel, "warning message")
	err4 := h.Log(ErrorLevel, "error message")
	err5 := h.Log(FatalLevel, "fatal message")

	// Kiểm tra không có lỗi
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		t.Errorf("ColoredConsoleHandler.Log() trả về lỗi: %v, %v, %v, %v, %v",
			err1, err2, err3, err4, err5)
	}
}

func TestConsoleHandlerClose(t *testing.T) {
	h := NewConsoleHandler(true)

	// Close không nên trả về lỗi
	err := h.Close()
	if err != nil {
		t.Errorf("ConsoleHandler.Close() error = %v", err)
	}
}
