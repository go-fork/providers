package scheduler

import (
	"testing"

	"github.com/go-fork/di"
)

// MockApp implements the interface required by ServiceProvider
type MockApp struct {
	container *di.Container
}

func (m *MockApp) Container() *di.Container {
	return m.container
}

func TestServiceProviderRegister(t *testing.T) {
	// Tạo DI container
	container := di.New()
	app := &MockApp{container: container}

	// Tạo service provider
	provider := NewServiceProvider()

	// Test Register
	provider.Register(app)

	// Kiểm tra scheduler đã được đăng ký vào container
	instance, err := container.Make("scheduler")
	if err != nil {
		t.Fatalf("Failed to get scheduler from container: %v", err)
	}

	scheduler, ok := instance.(Manager)
	if !ok {
		t.Fatal("Registered instance is not a Manager")
	}

	if scheduler == nil {
		t.Fatal("Scheduler is nil")
	}
}

func TestServiceProviderRegisterWithNilApp(t *testing.T) {
	provider := NewServiceProvider()

	// Test với app không implement Container()
	provider.Register("invalid-app")

	// Test này không nên panic
}

func TestServiceProviderRegisterWithNilContainer(t *testing.T) {
	app := &MockApp{container: nil}
	provider := NewServiceProvider()

	// Test với container nil - nên panic
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic when container is nil")
		}
	}()

	provider.Register(app)
}

func TestServiceProviderBoot(t *testing.T) {
	// Tạo DI container và đăng ký scheduler
	container := di.New()
	app := &MockApp{container: container}
	provider := NewServiceProvider()

	// Register trước
	provider.Register(app)

	// Test Boot
	provider.Boot(app)

	// Lấy scheduler và kiểm tra nó đã được start
	instance, err := container.Make("scheduler")
	if err != nil {
		t.Fatalf("Failed to get scheduler from container: %v", err)
	}

	scheduler, ok := instance.(Manager)
	if !ok {
		t.Fatal("Instance is not a Manager")
	}

	// Kiểm tra scheduler đã được start
	if !scheduler.IsRunning() {
		t.Fatal("Scheduler should be running after Boot")
	}

	// Cleanup
	scheduler.Stop()
}

func TestServiceProviderBootWithNilApp(t *testing.T) {
	provider := NewServiceProvider()

	// Test với app không implement Container()
	provider.Boot("invalid-app")

	// Test này không nên panic
}

func TestServiceProviderBootWithNilContainer(t *testing.T) {
	app := &MockApp{container: nil}
	provider := NewServiceProvider()

	// Test với container nil
	provider.Boot(app)

	// Test này không nên panic
}

func TestServiceProviderBootWithoutScheduler(t *testing.T) {
	// Tạo container nhưng không đăng ký scheduler
	container := di.New()
	app := &MockApp{container: container}
	provider := NewServiceProvider()

	// Test Boot mà không có scheduler trong container
	provider.Boot(app)

	// Test này không nên panic
}

func TestServiceProviderBootTwice(t *testing.T) {
	// Tạo DI container và đăng ký scheduler
	container := di.New()
	app := &MockApp{container: container}
	provider := NewServiceProvider()

	// Register trước
	provider.Register(app)

	// Boot lần đầu
	provider.Boot(app)

	// Boot lần hai - không nên start lại scheduler
	provider.Boot(app)

	// Lấy scheduler và kiểm tra nó vẫn đang chạy
	instance, err := container.Make("scheduler")
	if err != nil {
		t.Fatalf("Failed to get scheduler from container: %v", err)
	}

	scheduler, ok := instance.(Manager)
	if !ok {
		t.Fatal("Instance is not a Manager")
	}

	// Kiểm tra scheduler vẫn đang chạy
	if !scheduler.IsRunning() {
		t.Fatal("Scheduler should still be running after second Boot")
	}

	// Cleanup
	scheduler.Stop()
}

func TestNewServiceProvider(t *testing.T) {
	provider := NewServiceProvider()

	if provider == nil {
		t.Fatal("NewServiceProvider() returned nil")
	}

	// Kiểm tra provider implement đúng interface
	_, ok := provider.(di.ServiceProvider)
	if !ok {
		t.Fatal("Provider does not implement di.ServiceProvider interface")
	}
}
