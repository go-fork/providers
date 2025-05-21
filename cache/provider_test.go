package cache

import (
	"testing"

	"github.com/go-fork/di"
	"github.com/go-fork/providers/cache/driver"
)

// mockApp mimics an application with Container and Environment
type mockApp struct {
	container *di.Container
	env       string
	basePath  string
}

func (m *mockApp) Container() *di.Container {
	return m.container
}

func (m *mockApp) Environment() string {
	return m.env
}

func (m *mockApp) BasePath(paths ...string) string {
	path := m.basePath
	for _, p := range paths {
		path += "/" + p
	}
	return path
}

func TestServiceProviderRegister(t *testing.T) {
	t.Run("registers cache manager to container", func(t *testing.T) {
		// Arrange
		container := di.New()
		app := &mockApp{container: container}
		provider := NewServiceProvider()

		// Act
		provider.Register(app)

		// Assert
		if !container.Bound("cache") {
			t.Errorf("Expected 'cache' to be bound in container, but it wasn't")
		}

		if !container.Bound("cache.manager") {
			t.Errorf("Expected 'cache.manager' to be bound in container, but it wasn't")
		}

		cacheService, err := container.Make("cache")
		if err != nil {
			t.Errorf("Expected no error when making 'cache', got %v", err)
		}

		_, ok := cacheService.(Manager)
		if !ok {
			t.Errorf("Expected 'cache' to be of type Manager, but it wasn't")
		}

		managerService, err := container.Make("cache.manager")
		if err != nil {
			t.Errorf("Expected no error when making 'cache.manager', got %v", err)
		}

		if cacheService != managerService {
			t.Errorf("Expected 'cache' and 'cache.manager' to be the same instance, but they weren't")
		}
	})

	t.Run("does nothing when app doesn't have container", func(t *testing.T) {
		// Arrange
		provider := NewServiceProvider()
		app := struct{}{}

		// Act & Assert
		// Should not panic
		provider.Register(app)
	})
}

func TestServiceProviderBoot(t *testing.T) {
	t.Run("configures memory and file drivers", func(t *testing.T) {
		// Arrange
		container := di.New()
		app := &mockApp{
			container: container,
			env:       "testing",
			basePath:  "/tmp",
		}
		provider := NewServiceProvider()
		provider.Register(app)

		// Act
		provider.Boot(app)

		// Assert
		cacheService, err := container.Make("cache")
		if err != nil {
			t.Errorf("Expected no error when making 'cache', got %v", err)
		}

		manager := cacheService.(Manager)

		// Test that memory driver is set as default
		manager.Set("test-key", "test-value", 0)
		value, found := manager.Get("test-key")

		if !found {
			t.Errorf("Expected to find key in memory driver, but didn't")
		}
		if value != "test-value" {
			t.Errorf("Expected value to be %v, got %v", "test-value", value)
		}

		// Test that file driver is registered
		fileDriver, err := manager.Driver("file")
		if err != nil {
			t.Errorf("Expected no error when getting file driver, got %v", err)
		}
		if fileDriver == nil {
			t.Errorf("Expected file driver to be registered, but it wasn't")
		}

		// Check the type of the file driver
		_, ok := fileDriver.(*driver.FileDriver)
		if !ok {
			t.Errorf("Expected file driver to be of type *driver.FileDriver, but it wasn't")
		}
	})

	t.Run("does nothing when app doesn't have container", func(t *testing.T) {
		// Arrange
		provider := NewServiceProvider()
		app := struct{}{}

		// Act & Assert
		// Should not panic
		provider.Boot(app)
	})

	t.Run("does nothing when cache service not found", func(t *testing.T) {
		// Arrange
		container := di.New()
		app := &mockApp{
			container: container,
			env:       "testing",
			basePath:  "/tmp",
		}
		provider := NewServiceProvider()
		// Intentionally NOT calling provider.Register(app)

		// Act & Assert
		// Should not panic
		provider.Boot(app)
	})
}
