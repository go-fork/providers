package handler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "filehandler-test-*")
	if err != nil {
		t.Fatalf("Không thể tạo thư mục tạm: %v", err)
	}
	return dir
}

func TestFileHandlerLog(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn log
	logPath := filepath.Join(dir, "log-test.log")

	// Tạo handler với maxSize nhỏ để test rotation
	h, err := NewFileHandler(logPath, 50)
	if err != nil {
		t.Fatalf("NewFileHandler() error = %v", err)
	}
	defer h.Close()

	// Ghi log
	err = h.Log(InfoLevel, "test message 1")
	if err != nil {
		t.Errorf("Log() error = %v", err)
	}

	// Ghi log với args
	err = h.Log(WarningLevel, "test message %d", 2)
	if err != nil {
		t.Errorf("Log() với args error = %v", err)
	}

	// Đọc nội dung file
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Không thể đọc file log: %v", err)
	}

	// Kiểm tra nội dung
	contentStr := string(content)
	if !contains(contentStr, "[INFO]") || !contains(contentStr, "test message 1") {
		t.Errorf("Log không ghi đúng message 1: %q", contentStr)
	}
	if !contains(contentStr, "[WARNING]") || !contains(contentStr, "test message 2") {
		t.Errorf("Log không ghi đúng message 2: %q", contentStr)
	}

	// Ghi nhiều log để kích hoạt rotation
	for i := 0; i < 10; i++ {
		err = h.Log(ErrorLevel, "rotation test message %d", i)
		if err != nil {
			t.Errorf("Log() trong vòng lặp error = %v", err)
		}
	}

	// Kiểm tra file gốc vẫn tồn tại
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("File log gốc không tồn tại sau rotation")
	}

	// Kiểm tra nếu có ít nhất một file backup được tạo
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Không thể đọc thư mục: %v", err)
	}

	backupFound := false
	for _, file := range files {
		if file.Name() != "log-test.log" && contains(file.Name(), "log-test.log") {
			backupFound = true
			break
		}
	}

	if !backupFound {
		t.Error("Không tìm thấy file backup sau rotation")
	}
}

func TestFileHandlerClose(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn log
	logPath := filepath.Join(dir, "close-test.log")

	// Tạo handler
	h, err := NewFileHandler(logPath, 1024)
	if err != nil {
		t.Fatalf("NewFileHandler() error = %v", err)
	}

	// Ghi log
	err = h.Log(InfoLevel, "test before close")
	if err != nil {
		t.Errorf("Log() error = %v", err)
	}

	// Đóng file
	err = h.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Kiểm tra ghi log sau khi đóng - phải thất bại
	err = h.Log(InfoLevel, "test after close")
	if err == nil {
		t.Error("Log() sau Close() không trả về lỗi như mong đợi")
	}

	// Kiểm tra đóng hai lần
	err = h.Close()
	if err != nil {
		t.Errorf("Close() lần thứ hai error = %v, mong đợi nil", err)
	}
}

func TestNewFileHandlerWithInvalidPath(t *testing.T) {
	// Đường dẫn không hợp lệ cho file log (thư mục không tồn tại và không thể tạo)
	invalidPath := filepath.Join("/ this path doesn't exist /", "log.txt")

	// Cố gắng tạo handler với đường dẫn không hợp lệ
	_, err := NewFileHandler(invalidPath, 1024)
	if err == nil {
		t.Error("NewFileHandler() với đường dẫn không hợp lệ không trả về lỗi")
	}
}

func TestFileHandlerRotate(t *testing.T) {
	// Tạo thư mục tạm thời
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Đường dẫn log
	logPath := filepath.Join(dir, "rotate-test.log")

	// Tạo handler với maxSize rất nhỏ để kích hoạt rotation ngay lập tức
	h, err := NewFileHandler(logPath, 10)
	if err != nil {
		t.Fatalf("NewFileHandler() error = %v", err)
	}
	defer h.Close()

	// Ghi log để vượt quá kích thước tối đa và kích hoạt rotation
	messageBeforeRotation := "test1"
	err = h.Log(InfoLevel, messageBeforeRotation)
	if err != nil {
		t.Errorf("Log() (trước rotation) error = %v", err)
	}

	// Ghi log thứ hai để đảm bảo rotation xảy ra
	messageAfterRotation := "test2"
	err = h.Log(InfoLevel, messageAfterRotation)
	if err != nil {
		t.Errorf("Log() (sau rotation) error = %v", err)
	}

	// Kiểm tra file gốc chỉ chứa message sau rotation
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Không thể đọc file log sau rotation: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, messageAfterRotation) {
		t.Errorf("File log gốc không chứa message sau rotation: %q", contentStr)
	}

	if contains(contentStr, messageBeforeRotation) {
		t.Errorf("File log gốc vẫn chứa message trước rotation: %q", contentStr)
	}

	// Kiểm tra file backup
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Không thể đọc thư mục: %v", err)
	}

	backupFile := ""
	for _, file := range files {
		if file.Name() != filepath.Base(logPath) && strings.HasPrefix(file.Name(), filepath.Base(logPath)) {
			backupFile = filepath.Join(dir, file.Name())
			break
		}
	}

	if backupFile == "" {
		t.Fatal("Không tìm thấy file backup sau rotation")
	}

	// Đọc nội dung file backup
	backupContent, err := os.ReadFile(backupFile)
	if err != nil {
		t.Fatalf("Không thể đọc file backup: %v", err)
	}

	backupContentStr := string(backupContent)
	if !contains(backupContentStr, messageBeforeRotation) {
		t.Errorf("File backup không chứa message trước rotation: %q", backupContentStr)
	}
}

// Hàm trợ giúp kiểm tra chuỗi con
func contains(s, substr string) bool {
	return s != "" && substr != "" && len(s) >= len(substr) && strings.Contains(s, substr)
}
