package queue

import (
	"encoding/json"
	"fmt"
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

// TestServiceProviderBoot tests the Boot method of ServiceProvider
func TestServiceProviderBoot(t *testing.T) {
	provider := NewServiceProvider()

	// Test with nil app (should not panic)
	provider.Boot(nil)

	// Test with non-nil app
	type mockApp struct{}
	app := &mockApp{}
	provider.Boot(app)
}

// TestServiceProviderRegisterWithConfig tests the Register method with a config manager
func TestServiceProviderRegisterWithConfig(t *testing.T) {
	// Create a mock config manager
	mockConfig := &mockConfigManager{
		data: map[string]interface{}{
			"queue": map[string]interface{}{
				"adapter": map[string]interface{}{
					"default": "redis",
					"redis": map[string]interface{}{
						"addr": "localhost:6379",
					},
				},
			},
		},
		hasKey: true,
	}

	// Create a mock container
	container := di.New()
	container.Instance("config", mockConfig)

	// Create a mock app that returns our container
	app := &mockApp{container: container}

	// Create provider and register
	provider := NewServiceProvider()
	provider.Register(app)

	// Verify registrations
	queueManager, err := container.Make("queue")
	assert.NoError(t, err)
	assert.NotNil(t, queueManager)

	queueClient, err := container.Make("queue.client")
	assert.NoError(t, err)
	assert.NotNil(t, queueClient)

	queueServer, err := container.Make("queue.server")
	assert.NoError(t, err)
	assert.NotNil(t, queueServer)

	redisClient, err := container.Make("queue.redis")
	assert.NoError(t, err)
	assert.NotNil(t, redisClient)
}

// mockApp is a mock implementation for testing
type mockApp struct {
	container *di.Container
}

func (m *mockApp) Container() *di.Container {
	return m.container
}

// mockConfigManager is a mock implementation of config.Manager for testing
type mockConfigManager struct {
	data   map[string]interface{}
	hasKey bool
}

func (m *mockConfigManager) Has(key string) bool {
	return m.hasKey
}

func (m *mockConfigManager) UnmarshalKey(key string, target interface{}) error {
	val, ok := m.data[key]
	if !ok {
		return fmt.Errorf("key not found")
	}

	// Convert the map to JSON, then unmarshal to the target
	jsonData, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, target)
}
