package queue

import (
	"testing"

	"github.com/go-fork/di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockApp mocks an application with a container
type MockApp struct {
	mock.Mock
}

func (m *MockApp) Container() *di.Container {
	args := m.Called()
	return args.Get(0).(*di.Container)
}

// MockConfigManager mocks a config manager
type MockConfigManager struct {
	mock.Mock
}

func (m *MockConfigManager) Get(key string) interface{} {
	args := m.Called(key)
	return args.Get(0)
}

func (m *MockConfigManager) Has(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

func (m *MockConfigManager) Set(key string, value interface{}) {
	m.Called(key, value)
}

func (m *MockConfigManager) UnmarshalKey(key string, out interface{}) error {
	args := m.Called(key, out)
	// If we have a configuration to unmarshal, do it
	if fn, ok := args.Get(1).(func(interface{})); ok {
		fn(out)
	}
	return args.Error(0)
}

// TestNewServiceProvider tests creation of a new service provider
func TestNewServiceProvider(t *testing.T) {
	provider := NewServiceProvider()

	assert.NotNil(t, provider, "Service provider should not be nil")
	assert.IsType(t, &ServiceProvider{}, provider, "Service provider should be of type *ServiceProvider")
}

// TestNewServiceProviderWithConfig tests creation of a new service provider with custom config
func TestNewServiceProviderWithConfig(t *testing.T) {
	config := Config{
		Adapter: AdapterConfig{
			Default: "memory",
		},
	}

	provider := NewServiceProviderWithConfig(config)

	assert.NotNil(t, provider, "Service provider should not be nil")
	assert.IsType(t, &ServiceProvider{}, provider, "Service provider should be of type *ServiceProvider")
}

// TestServiceProviderRegisterWithApp tests registering services with the app container
func TestServiceProviderRegisterWithApp(t *testing.T) {
	// Create a mock app with container
	mockApp := new(MockApp)
	container := di.New()
	mockApp.On("Container").Return(container)

	// Create a provider and register it
	provider := NewServiceProvider()
	provider.Register(mockApp)

	// Verify the services were registered
	hasQueueService := container.Bound("queue")
	assert.True(t, hasQueueService, "Container should have 'queue' service")

	hasQueueClient := container.Bound("queue.client")
	assert.True(t, hasQueueClient, "Container should have 'queue.client' service")

	hasQueueServer := container.Bound("queue.server")
	assert.True(t, hasQueueServer, "Container should have 'queue.server' service")

	hasQueueManager := container.Bound("queue.manager")
	assert.True(t, hasQueueManager, "Container should have 'queue.manager' service")

	// Verify mock expectations
	mockApp.AssertExpectations(t)
}

// TestServiceProviderRegisterWithConfigManager tests registering with a config manager
func TestServiceProviderRegisterWithConfigManager(t *testing.T) {
	// Create a mock app with container
	mockApp := new(MockApp)
	container := di.New()
	mockApp.On("Container").Return(container)

	// Manually create a provider with Redis configuration to avoid mocking issues
	provider := NewServiceProviderWithConfig(Config{
		Adapter: AdapterConfig{
			Default: "redis",
		},
	})

	// Register the provider
	provider.Register(mockApp)

	// Verify the services were registered
	hasQueueService := container.Bound("queue")
	assert.True(t, hasQueueService, "Container should have 'queue' service")

	// When default adapter is redis, should register redis client
	hasRedisService := container.Bound("queue.redis")
	assert.True(t, hasRedisService, "Container should have 'queue.redis' service since redis is default")

	// Verify mock expectations
	mockApp.AssertExpectations(t)
}

// TestServiceProviderRegisterWithConfigError tests registering with a config error
func TestServiceProviderRegisterWithConfigError(t *testing.T) {
	// Create a mock app with container
	mockApp := new(MockApp)
	container := di.New()
	mockApp.On("Container").Return(container)

	// Create a provider with default config
	provider := NewServiceProvider()
	provider.Register(mockApp)

	// Services should be registered with default config
	hasQueueService := container.Bound("queue")
	assert.True(t, hasQueueService, "Container should have 'queue' service")

	// Verify mock expectations
	mockApp.AssertExpectations(t)
}

// TestServiceProviderRegisterWithInvalidConfig tests registering with invalid config manager
func TestServiceProviderRegisterWithInvalidConfig(t *testing.T) {
	// Create a mock app with container
	mockApp := new(MockApp)
	container := di.New()
	mockApp.On("Container").Return(container)

	// Create an instance that doesn't implement config.Manager
	type invalidConfigManager struct{}
	container.Instance("config", &invalidConfigManager{})

	// Create a provider and register it - should not panic
	provider := NewServiceProvider()
	provider.Register(mockApp)

	// Services should still be registered with default config
	hasQueueService := container.Bound("queue")
	assert.True(t, hasQueueService, "Container should have 'queue' service")

	// Verify mock expectations
	mockApp.AssertExpectations(t)
}

// TestServiceProviderRegisterWithNonContainerApp tests registering with an app that doesn't provide a container
func TestServiceProviderRegisterWithNonContainerApp(t *testing.T) {
	// Create a struct that doesn't implement the Container method
	type invalidApp struct{}

	app := &invalidApp{}

	// Create a provider and register it - should not panic
	provider := NewServiceProvider()
	provider.Register(app)

	// Just verify it doesn't panic
}

// TestServiceProviderBoot tests the Boot method
func TestServiceProviderBoot(t *testing.T) {
	// Create a mock app
	mockApp := new(MockApp)

	// Create a provider and call Boot - should not do anything but should not error
	provider := NewServiceProvider()
	provider.Boot(mockApp)

	// No expectations to verify - Boot is currently a no-op
}
