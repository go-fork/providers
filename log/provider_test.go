package log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-fork/di"
	"github.com/go-fork/providers/log/handler"
)

// MockApplication triển khai interface cần thiết cho ServiceProvider
type MockApplication struct {
	container *di.Container
	basePath  string
}

func NewMockApplication() *MockApplication {
	return &MockApplication{
		container: di.New(),
		basePath:  os.TempDir(),
	}
}

func (m *MockApplication) Container() *di.Container {
	return m.container
}

func (m *MockApplication) BasePath(paths ...string) string {
	result := m.basePath
	for _, path := range paths {
		result = filepath.Join(result, path)
	}
	return result
}

func TestNewServiceProvider(t *testing.T) {
	provider := NewServiceProvider()
	if provider == nil {
		t.Fatal("NewServiceProvider() trả về nil")
	}
}

func TestServiceProviderRegister(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider với application
	provider.Register(app)

	// Kiểm tra manager được đăng ký trong container
	container := app.Container()

	// Kiểm tra binding "log"
	managerInstance, exists := container.Make("log")
	if !exists {
		t.Fatal("ServiceProvider không đăng ký binding 'log'")
	}

	manager, ok := managerInstance.(Manager)
	if !ok {
		t.Fatalf("Binding 'log' không phải kiểu Manager, got %T", managerInstance)
	}

	// Kiểm tra binding "log.manager"
	managerInstance2, exists := container.Make("log.manager")
	if !exists {
		t.Fatal("ServiceProvider không đăng ký binding 'log.manager'")
	}

	if managerInstance != managerInstance2 {
		t.Error("'log' và 'log.manager' trỏ đến các instance khác nhau")
	}

	// Kiểm tra handlers được thiết lập
	if consoleHandler, exists := manager.GetHandler("console"); !exists {
		t.Error("Manager không có console handler")
	} else {
		// Kiểm tra đúng kiểu handler
		_, ok := consoleHandler.(*handler.ConsoleHandler)
		if !ok {
			t.Errorf("Console handler không đúng kiểu, got %T", consoleHandler)
		}
	}

	// Kiểm tra thư mục logs tồn tại
	logsPath := app.BasePath("storage", "logs")
	if _, err := os.Stat(logsPath); os.IsNotExist(err) {
		t.Errorf("Thư mục logs không được tạo tại %q", logsPath)
	} else {
		// Dọn dẹp sau khi test
		defer os.RemoveAll(logsPath)
	}

	// Kiểm tra file handler
	if fileHandler, exists := manager.GetHandler("file"); !exists {
		t.Error("Manager không có file handler")
	} else {
		// Kiểm tra đúng kiểu handler
		_, ok := fileHandler.(*handler.FileHandler)
		if !ok {
			t.Errorf("File handler không đúng kiểu, got %T", fileHandler)
		}
	}
}

func TestServiceProviderRegisterWithInvalidApp(t *testing.T) {
	// Tạo một đối tượng không triển khai interface cần thiết
	invalidApp := struct{}{}

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider với invalid application
	// Không nên panic
	provider.Register(invalidApp)
}

func TestServiceProviderBoot(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Tạo service provider
	provider := NewServiceProvider()

	// Boot không làm gì nhưng nên chạy không lỗi
	provider.Boot(app)
}
