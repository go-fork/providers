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
	managerInstance, err := container.Make("log")
	if err != nil {
		t.Fatal("ServiceProvider không đăng ký binding 'log':", err)
	}

	manager, ok := managerInstance.(Manager)
	if !ok {
		t.Fatalf("Binding 'log' không phải kiểu Manager, got %T", managerInstance)
	}

	// Kiểm tra binding "log.manager"
	managerInstance2, err := container.Make("log.manager")
	if err != nil {
		t.Fatal("ServiceProvider không đăng ký binding 'log.manager':", err)
	}

	if managerInstance != managerInstance2 {
		t.Error("'log' và 'log.manager' trỏ đến các instance khác nhau")
	}

	// Kiểm tra handlers được thiết lập
	// Vì GetHandler là một phương thức mở rộng, chúng ta cần type assertion trước
	if defaultManager, ok := manager.(*DefaultManager); ok {
		if consoleHandler, exists := defaultManager.GetHandler("console"); !exists {
			t.Error("Manager không có console handler")
		} else {
			// Kiểm tra đúng kiểu handler
			_, ok := consoleHandler.(*handler.ConsoleHandler)
			if !ok {
				t.Errorf("Console handler không đúng kiểu, got %T", consoleHandler)
			}
		}
	} else {
		t.Error("Manager không phải kiểu DefaultManager")
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
	if defaultManager, ok := manager.(*DefaultManager); ok {
		if fileHandler, exists := defaultManager.GetHandler("file"); !exists {
			t.Error("Manager không có file handler")
		} else {
			// Kiểm tra đúng kiểu handler
			_, ok := fileHandler.(*handler.FileHandler)
			if !ok {
				t.Errorf("File handler không đúng kiểu, got %T", fileHandler)
			}
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

// Thêm test case mới kiểm tra chi tiết container.Instance
func TestServiceProviderContainerInstance(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider
	provider.Register(app)

	// Lấy container
	container := app.Container()

	// Kiểm tra xem binding 'log' đã được đăng ký chưa
	if !container.Bound("log") {
		t.Fatal("Binding 'log' chưa được đăng ký trong container")
	}

	// Kiểm tra xem binding 'log.manager' đã được đăng ký chưa
	if !container.Bound("log.manager") {
		t.Fatal("Binding 'log.manager' chưa được đăng ký trong container")
	}
	// Kiểm tra rằng cả hai binding đều trỏ đến cùng một instance
	logManager1, err := container.Make("log")
	if err != nil {
		t.Fatal("Không thể resolve binding 'log':", err)
	}

	logManager2, err := container.Make("log.manager")
	if err != nil {
		t.Fatal("Không thể resolve binding 'log.manager':", err)
	}

	if logManager1 != logManager2 {
		t.Error("container.Instance() không tạo singleton cho 'log' và 'log.manager'")
	}
}

// Test truy cập đến Manager thông qua Container.MustMake
func TestManagerAccessThroughContainer(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider
	provider.Register(app)

	// Lấy container
	container := app.Container()

	// Sử dụng MustMake để lấy manager
	manager := container.MustMake("log").(Manager)

	// Kiểm tra các phương thức của manager
	// Log một tin nhắn test (không cần assert kết quả, chỉ kiểm tra không panic)
	manager.Debug("Test debug message")
	manager.Info("Test info message")
	manager.Warning("Test warning message")
	manager.Error("Test error message")

	// Đóng manager và cleanup
	err := manager.Close()
	if err != nil {
		t.Errorf("Không thể đóng manager: %v", err)
	}
}

// Test ServiceProvider với một container được reset
func TestServiceProviderWithContainerReset(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Reset container trước
	app.Container().Reset()

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider
	provider.Register(app)

	// Kiểm tra binding 'log' tồn tại sau khi reset và register
	if !app.Container().Bound("log") {
		t.Fatal("Binding 'log' không tồn tại sau khi reset container và register lại")
	}

	// Kiểm tra service có hoạt động đúng
	manager, err := app.Container().Make("log")
	if err != nil {
		t.Fatal("Không thể resolve binding 'log':", err)
	}

	// Đảm bảo đúng kiểu
	_, ok := manager.(Manager)
	if !ok {
		t.Fatalf("Binding 'log' không phải Manager, nhận được kiểu %T", manager)
	}
}

// Test cách container giải quyết phụ thuộc qua Call
func TestContainerResolveDependencies(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider
	provider.Register(app)

	// Đăng ký kiểu Manager để container có thể inject
	container := app.Container()
	container.Bind("log.Manager", func(c *di.Container) interface{} {
		manager, _ := c.Make("log")
		return manager
	})

	// Gọi một hàm với tham số là Manager thông qua container.Call
	called := false
	result, err := app.Container().Call(func(logger Manager) {
		// Kiểm tra xem Manager được truyền vào đúng không
		if logger == nil {
			t.Error("Manager được truyền vào là nil")
		} else {
			called = true
			// Thử log một message
			logger.Info("Test message từ container.Call")
		}
	})

	if err != nil {
		t.Fatalf("container.Call thất bại: %v", err)
	}

	if !called {
		t.Error("Hàm callback không được gọi")
	}

	// Kiểm tra kết quả của container.Call
	if result != nil {
		t.Errorf("container.Call trả về không mong đợi: %v", result)
	}
}

// TestGetHandler kiểm tra xem Manager có phương thức GetHandler không
func (m *DefaultManager) GetHandler(name string) (handler.Handler, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	h, exists := m.handlers[name]
	return h, exists
}

// TestBindingTypes kiểm tra các loại binding khác nhau trong container
func TestBindingTypes(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Reset container trước
	app.Container().Reset()

	// Tạo các binding trực tiếp vào container
	container := app.Container()

	// Kiểm tra Bind
	container.Bind("test.bind", func(c *di.Container) interface{} {
		return "bind-value"
	})

	bindValue, err := container.Make("test.bind")
	if err != nil {
		t.Fatal("Không thể resolve binding 'test.bind':", err)
	}
	if bindValue != "bind-value" {
		t.Errorf("Bind không trả về giá trị đúng, expected 'bind-value', got %v", bindValue)
	}

	// Kiểm tra Singleton
	counter := 0
	container.Singleton("test.singleton", func(c *di.Container) interface{} {
		counter++
		return counter
	})

	// Gọi Make nhiều lần, counter chỉ nên tăng một lần nếu là singleton
	val1, err := container.Make("test.singleton")
	if err != nil {
		t.Fatal("Không thể resolve singleton:", err)
	}
	val2, err := container.Make("test.singleton")
	if err != nil {
		t.Fatal("Không thể resolve singleton:", err)
	}

	if val1 != val2 {
		t.Error("Singleton không trả về cùng một instance")
	}
	if val1.(int) != 1 || val2.(int) != 1 {
		t.Errorf("Singleton không tạo đúng giá trị, got %v và %v", val1, val2)
	}

	// Kiểm tra Alias
	container.Alias("test.bind", "test.alias")
	aliasValue, err := container.Make("test.alias")
	if err != nil {
		t.Fatal("Không thể resolve alias:", err)
	}
	if aliasValue != "bind-value" {
		t.Errorf("Alias không trỏ đến binding gốc, expected 'bind-value', got %v", aliasValue)
	}

	// Kiểm tra Instance
	myInstance := &struct{ Name string }{"test-instance"}
	container.Instance("test.instance", myInstance)

	instanceValue, err := container.Make("test.instance")
	if err != nil {
		t.Fatal("Không thể resolve instance:", err)
	}
	if instanceValue != myInstance {
		t.Error("Instance không trả về đúng đối tượng đã đăng ký")
	}
}

// TestServiceProviderWithCustomLogger kiểm tra Provider hoạt động với logger tùy chỉnh
func TestServiceProviderWithCustomLogger(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Tạo một manager tùy chỉnh trước
	customManager := NewManager()

	// Đăng ký manager tùy chỉnh trong container
	app.Container().Instance("custom.log", customManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider (điều này sẽ ghi đè lên custom.log?)
	provider.Register(app)

	// Kiểm tra binding tùy chỉnh và binding tiêu chuẩn
	standardManager, err := app.Container().Make("log")
	if err != nil {
		t.Fatal("Không thể resolve binding 'log':", err)
	}

	// Kiểm tra rằng manager tùy chỉnh vẫn tồn tại và không bị ghi đè
	customValue, err := app.Container().Make("custom.log")
	if err != nil {
		t.Fatal("Không thể resolve binding 'custom.log':", err)
	}

	// Các manager nên khác nhau
	if customValue == standardManager {
		t.Error("Provider ghi đè lên các binding tùy chỉnh")
	}

	// Kiểm tra rằng cả hai đều triển khai interface Manager
	_, ok1 := customValue.(Manager)
	_, ok2 := standardManager.(Manager)
	if !ok1 || !ok2 {
		t.Error("Một trong các binding không triển khai interface Manager")
	}
}

// TestServiceProviderMultipleRegistrations kiểm tra việc đăng ký provider nhiều lần
func TestServiceProviderMultipleRegistrations(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider lần đầu
	provider.Register(app)

	// Lấy manager sau lần đăng ký đầu tiên
	manager1, err := app.Container().Make("log")
	if err != nil {
		t.Fatal("Không thể resolve binding 'log' sau lần đăng ký đầu tiên:", err)
	}

	// Đăng ký provider lần thứ hai
	provider.Register(app)

	// Lấy manager sau lần đăng ký thứ hai
	manager2, err := app.Container().Make("log")
	if err != nil {
		t.Fatal("Không thể resolve binding 'log' sau lần đăng ký thứ hai:", err)
	}

	// Kiểm tra nếu chúng là cùng một instance (tùy thuộc vào cách triển khai)
	// Lưu ý: Nếu Container.Instance() tạo một singleton, các lần đăng ký sau có thể không ghi đè
	if manager1 != manager2 {
		t.Log("Đăng ký lần thứ hai tạo một instance mới - OK nếu được thiết kế như vậy")
	} else {
		t.Log("Đăng ký lần thứ hai giữ nguyên instance - OK nếu được thiết kế như vậy")
	}
}

// TestServiceProviderWithCustomContainer kiểm tra việc sử dụng provider với một container tùy chỉnh
func TestServiceProviderWithCustomContainer(t *testing.T) {
	// Tạo một container tùy chỉnh riêng biệt
	customContainer := di.New()

	// Tạo mock application với container tùy chỉnh
	app := &MockApplication{
		container: customContainer,
		basePath:  os.TempDir(),
	}

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider vào container tùy chỉnh
	provider.Register(app)

	// Kiểm tra xem container có binding 'log' chưa
	if !customContainer.Bound("log") {
		t.Fatal("Container tùy chỉnh không có binding 'log' sau khi register provider")
	}

	// Xác minh binding đúng kiểu và hoạt động
	logService, err := customContainer.Make("log")
	if err != nil {
		t.Fatal("Không thể resolve binding 'log':", err)
	}

	// Kiểm tra kiểu
	logManager, ok := logService.(Manager)
	if !ok {
		t.Fatalf("Binding 'log' không phải Manager, got: %T", logService)
	}

	// Kiểm tra chức năng cơ bản
	// Log một message (không kiểm tra kết quả, chỉ xem nó có panic không)
	logManager.Info("Test message")
}

// TestContainerBindingResolution kiểm tra việc giải quyết các binding phức tạp
func TestContainerBindingResolution(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Reset container
	app.Container().Reset()

	// Tạo một service provider
	provider := NewServiceProvider()

	// Đăng ký provider
	provider.Register(app)

	// Lấy container
	container := app.Container()

	// Thêm một binding phụ thuộc vào log manager
	container.Bind("custom.logger", func(c *di.Container) interface{} {
		// Lấy log manager từ container
		manager, err := c.Make("log")
		if err != nil {
			t.Fatal("Không thể resolve dependency 'log':", err)
		}

		// Trả về một struct sử dụng log manager
		return struct {
			LogManager Manager
			Name       string
		}{
			LogManager: manager.(Manager),
			Name:       "CustomLogger",
		}
	})

	// Giải quyết binding
	customLogger, err := container.Make("custom.logger")
	if err != nil {
		t.Fatal("Không thể resolve binding 'custom.logger':", err)
	}

	// Kiểm tra cấu trúc được trả về
	loggerStruct, ok := customLogger.(struct {
		LogManager Manager
		Name       string
	})

	if !ok {
		t.Fatalf("Binding 'custom.logger' không trả về kiểu đúng, got: %T", customLogger)
	}

	// Kiểm tra các trường
	if loggerStruct.Name != "CustomLogger" {
		t.Errorf("Tên không đúng, expected 'CustomLogger', got %s", loggerStruct.Name)
	}

	if loggerStruct.LogManager == nil {
		t.Error("LogManager là nil")
	}
}

// TestServiceProviderBootSequence kiểm tra trình tự boot và register
func TestServiceProviderBootSequence(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Reset container
	app.Container().Reset()

	// Tạo service provider
	provider := NewServiceProvider()

	// Thứ tự chuẩn: đăng ký trước, sau đó boot
	provider.Register(app)
	provider.Boot(app)

	// Kiểm tra xem các binding đã được thiết lập chưa
	if !app.Container().Bound("log") {
		t.Error("Binding 'log' không được thiết lập sau Register+Boot")
	}

	// Reset lại để thử thứ tự khác
	app.Container().Reset()

	// Thứ tự không chuẩn: boot trước, đăng ký sau
	// (Không nên ảnh hưởng đến kết quả cuối cùng)
	provider.Boot(app)
	provider.Register(app)

	// Kiểm tra xem các binding đã được thiết lập chưa
	if !app.Container().Bound("log") {
		t.Error("Binding 'log' không được thiết lập sau Boot+Register")
	}
}

// TestLoggerConfiguration kiểm tra cấu hình của logger được thiết lập bởi provider
func TestLoggerConfiguration(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider
	provider.Register(app)

	// Lấy log manager
	logManager, err := app.Container().Make("log")
	if err != nil {
		t.Fatal("Không thể resolve binding 'log':", err)
	}

	// Kiểm tra logManager
	defaultManager, ok := logManager.(*DefaultManager)
	if !ok {
		t.Fatalf("LogManager không phải DefaultManager, got %T", logManager)
	}

	// Kiểm tra cấu hình mặc định
	// Vì handlers là private, cần sử dụng GetHandler
	consoleHandler, exists := defaultManager.GetHandler("console")
	if !exists {
		t.Error("ConsoleHandler không tồn tại")
	} else {
		// Kiểm tra kiểu handler
		if _, ok := consoleHandler.(*handler.ConsoleHandler); !ok {
			t.Errorf("ConsoleHandler không đúng kiểu, got %T", consoleHandler)
		}
	}

	// Kiểm tra file handler và đường dẫn
	fileHandler, exists := defaultManager.GetHandler("file")
	if !exists {
		t.Error("FileHandler không tồn tại")
	} else {
		// Kiểm tra kiểu handler
		if _, ok := fileHandler.(*handler.FileHandler); !ok {
			t.Errorf("FileHandler không đúng kiểu, got %T", fileHandler)
		}
	}

	// Đóng manager sau khi kiểm tra
	err = defaultManager.Close()
	if err != nil {
		t.Errorf("Không thể đóng log manager: %v", err)
	}
}

// TestServiceProviderBootInDetail kiểm tra chi tiết phương thức Boot
func TestServiceProviderBootInDetail(t *testing.T) {
	// Tạo mock application
	app := NewMockApplication()

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider
	provider.Register(app)

	// Khi kiểm tra Boot, tuy nó không làm gì trong triển khai hiện tại,
	// nhưng nên đảm bảo nó được chạy mà không gây ra lỗi
	provider.Boot(app)

	// Tạo một invalid app để kiểm tra Boot có xử lý đúng khi nhận invalid input
	invalidApp := struct{}{}

	// Boot không nên panic với invalid input
	provider.Boot(invalidApp)

	// Boot nhiều lần cũng không nên gây ra vấn đề
	provider.Boot(app)
	provider.Boot(app)
}

// TestRegisterBootChain kiểm tra chuỗi Register và Boot ở các trình tự khác nhau
func TestRegisterBootChain(t *testing.T) {
	// Trường hợp 1: Chuỗi bootRegisterBoot - xem Boot trước Register có gây ra vấn đề gì không
	app1 := NewMockApplication()
	provider1 := NewServiceProvider()

	// Boot trước Register
	provider1.Boot(app1)
	provider1.Register(app1)
	provider1.Boot(app1)

	// Kiểm tra log đã được đăng ký
	if !app1.Container().Bound("log") {
		t.Error("Binding 'log' không tồn tại sau Boot+Register+Boot")
	}

	// Trường hợp 2: Nhiều cycle Register-Boot
	app2 := NewMockApplication()
	provider2 := NewServiceProvider()

	// Nhiều lượt Register-Boot liên tiếp
	for i := 0; i < 3; i++ {
		provider2.Register(app2)
		provider2.Boot(app2)
	}

	// Kiểm tra log đã được đăng ký và có thể sử dụng
	manager, err := app2.Container().Make("log")
	if err != nil {
		t.Fatal("Không thể resolve binding 'log' sau nhiều lượt Register-Boot:", err)
	}

	// Kiểm tra kiểu
	_, ok := manager.(Manager)
	if !ok {
		t.Fatalf("Binding 'log' không phải Manager, nhận được kiểu %T", manager)
	}
}

// TestServiceProviderErrorHandling kiểm tra xử lý lỗi của ServiceProvider
func TestServiceProviderErrorHandling(t *testing.T) {
	// Tạo một mock application đặc biệt, trong đó BasePath trả về một đường dẫn không tồn tại
	app := &MockApplication{
		container: di.New(),
		basePath:  "/non/existent/path", // Đường dẫn không tồn tại
	}

	// Tạo service provider
	provider := NewServiceProvider()

	// Register vẫn nên hoạt động mà không panic, mặc dù có thể không tạo được thư mục logs
	provider.Register(app)

	// Kiểm tra xem log manager vẫn được đăng ký
	if !app.Container().Bound("log") {
		t.Error("Binding 'log' không tồn tại sau khi Register với đường dẫn không hợp lệ")
	}

	// Lấy manager để kiểm tra nó đã được cấu hình đúng không
	manager, err := app.Container().Make("log")
	if err != nil {
		t.Fatal("Không thể resolve binding 'log':", err)
	}

	// Kiểm tra kiểu
	defaultManager, ok := manager.(*DefaultManager)
	if !ok {
		t.Fatalf("Binding 'log' không phải DefaultManager, nhận được kiểu %T", manager)
	}

	// Vẫn nên có console handler
	if _, exists := defaultManager.GetHandler("console"); !exists {
		t.Error("ConsoleHandler không tồn tại khi Register với đường dẫn không hợp lệ")
	}
}

// TestServiceProviderWithoutContainer kiểm tra provider hoạt động với ứng dụng không có Container()
func TestServiceProviderWithoutContainer(t *testing.T) {
	// Tạo một mock application không có Container()
	appWithoutContainer := struct {
		BasePath func(paths ...string) string
	}{
		BasePath: func(paths ...string) string {
			return filepath.Join(os.TempDir(), filepath.Join(paths...))
		},
	}

	// Tạo service provider
	provider := NewServiceProvider()

	// Register không nên panic
	provider.Register(appWithoutContainer)

	// Boot không nên panic
	provider.Boot(appWithoutContainer)
}

// TestExhaustiveBootWithDifferentAppTypes kiểm tra Boot với nhiều loại app khác nhau
func TestExhaustiveBootWithDifferentAppTypes(t *testing.T) {
	// Tạo service provider
	provider := NewServiceProvider()

	// 1. App là nil
	var nilApp interface{} = nil
	provider.Boot(nilApp)

	// 2. App là empty struct
	emptyApp := struct{}{}
	provider.Boot(emptyApp)

	// 3. App là struct với một số phương thức nhưng không phải tất cả
	partialApp := struct {
		BasePath func(paths ...string) string
	}{
		BasePath: func(paths ...string) string {
			return os.TempDir()
		},
	}
	provider.Boot(partialApp)

	// 4. App là struct với container không hợp lệ
	invalidContainerApp := struct {
		Container func() interface{}
	}{
		Container: func() interface{} {
			return nil
		},
	}
	provider.Boot(invalidContainerApp)

	// 5. App là struct với container hợp lệ nhưng không có basepath
	containerOnlyApp := struct {
		Container func() *di.Container
	}{
		Container: func() *di.Container {
			return di.New()
		},
	}
	provider.Boot(containerOnlyApp)

	// 6. App đầy đủ, nhưng được gọi Boot nhiều lần
	fullApp := NewMockApplication()
	for i := 0; i < 5; i++ {
		provider.Boot(fullApp)
	}
}
