package handler

import (
	"os"
	"strings"
	"testing"
)

// Tạo một buffer để bắt đầu ra của console
type captureOutput struct {
	oldStdout *os.File
	oldStderr *os.File
	readPipe  *os.File
	writePipe *os.File
}

func newCaptureOutput() (*captureOutput, error) {
	// Lưu trữ os.Stdout và os.Stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	// Tạo pipe để capture đầu ra
	rOut, wOut, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	// Đặt pipe làm đầu ra mới
	os.Stdout = wOut
	os.Stderr = wOut

	return &captureOutput{
		oldStdout: oldStdout,
		oldStderr: oldStderr,
		readPipe:  rOut,
		writePipe: wOut,
	}, nil
}

func (c *captureOutput) read() (string, error) {
	// Đóng write pipe để flush đầu ra
	c.writePipe.Close()

	// Đọc tất cả đầu ra từ pipe
	buf := make([]byte, 1024)
	n, err := c.readPipe.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}

	// Khôi phục os.Stdout và os.Stderr
	os.Stdout = c.oldStdout
	os.Stderr = c.oldStderr

	// Đóng read pipe
	c.readPipe.Close()

	return string(buf[:n]), nil
}

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
	// Tạo capture để bắt đầu ra
	capture, err := newCaptureOutput()
	if err != nil {
		t.Fatalf("Không thể tạo capture: %v", err)
	}

	// Test non-colored console handler
	h := &ConsoleHandler{colored: false}

	// Ghi log ở mỗi level
	h.Log(DebugLevel, "debug message")
	h.Log(InfoLevel, "info message")
	h.Log(WarningLevel, "warning message")
	h.Log(ErrorLevel, "error message")
	h.Log(FatalLevel, "fatal message")

	// Đọc đầu ra
	output, err := capture.read()
	if err != nil {
		t.Fatalf("Không thể đọc đầu ra: %v", err)
	}

	// Kiểm tra đầu ra
	if !strings.Contains(output, "DEBUG") || !strings.Contains(output, "debug message") {
		t.Error("ConsoleHandler.Log(DebugLevel) không ghi đúng thông điệp")
	}
	if !strings.Contains(output, "INFO") || !strings.Contains(output, "info message") {
		t.Error("ConsoleHandler.Log(InfoLevel) không ghi đúng thông điệp")
	}
	if !strings.Contains(output, "WARNING") || !strings.Contains(output, "warning message") {
		t.Error("ConsoleHandler.Log(WarningLevel) không ghi đúng thông điệp")
	}
}

func TestColoredConsoleHandler(t *testing.T) {
	// Tạo capture để bắt đầu ra
	capture, err := newCaptureOutput()
	if err != nil {
		t.Fatalf("Không thể tạo capture: %v", err)
	}

	// Test colored console handler
	h := &ConsoleHandler{colored: true}

	// Ghi log ở mỗi level
	h.Log(DebugLevel, "debug message")
	h.Log(InfoLevel, "info message")
	h.Log(WarningLevel, "warning message")
	h.Log(ErrorLevel, "error message")
	h.Log(FatalLevel, "fatal message")

	// Đọc đầu ra
	output, err := capture.read()
	if err != nil {
		t.Fatalf("Không thể đọc đầu ra: %v", err)
	}

	// Kiểm tra đầu ra có mã ANSI
	if !strings.Contains(output, "\033[") {
		t.Error("ConsoleHandler với colored=true không sử dụng mã ANSI")
	}
}

func TestConsoleHandlerWithArgs(t *testing.T) {
	// Tạo capture để bắt đầu ra
	capture, err := newCaptureOutput()
	if err != nil {
		t.Fatalf("Không thể tạo capture: %v", err)
	}

	// Test handler
	h := NewConsoleHandler(false)

	// Ghi log với tham số định dạng
	// Lưu ý: phương thức Log không dùng args để định dạng message
	h.Log(InfoLevel, "formatted message 123")

	// Đọc đầu ra
	output, err := capture.read()
	if err != nil {
		t.Fatalf("Không thể đọc đầu ra: %v", err)
	}

	// Kiểm tra đầu ra có định dạng đúng
	if !strings.Contains(output, "formatted message 123") {
		t.Error("ConsoleHandler.Log() không ghi đúng thông điệp")
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

func TestConsoleHandlerColorize(t *testing.T) {
	h := &ConsoleHandler{colored: true}
	message := "test message"

	// Kiểm tra colorize với các level khác nhau
	tests := []struct {
		name      string
		level     Level
		colorCode string
	}{
		{"Debug", DebugLevel, "\033[36m"},     // Cyan
		{"Info", InfoLevel, "\033[32m"},       // Green
		{"Warning", WarningLevel, "\033[33m"}, // Yellow
		{"Error", ErrorLevel, "\033[31m"},     // Red
		{"Fatal", FatalLevel, "\033[35m"},     // Magenta
		{"Unknown", Level(99), "\033[0m"},     // Default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colored := h.colorize(tt.level, message)
			if !strings.Contains(colored, tt.colorCode) {
				t.Errorf("colorize() không sử dụng mã màu đúng %s, got: %s", tt.colorCode, colored)
			}
			if !strings.Contains(colored, message) {
				t.Errorf("colorize() không bao gồm thông điệp gốc")
			}
			if !strings.Contains(colored, "\033[0m") {
				t.Errorf("colorize() không bao gồm mã reset ở cuối")
			}
		})
	}
}
