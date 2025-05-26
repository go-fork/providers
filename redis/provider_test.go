package redis

import (
	"testing"

	"github.com/go-fork/di"
	"github.com/go-fork/providers/config/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockApp implements the container interface for testing
type mockApp struct {
	container *di.Container
}

func (a *mockApp) Container() *di.Container {
	return a.container
}

// setupTestRedisConfig creates a Redis config for testing
func setupTestRedisConfig() *Config {
	return &Config{
		Client: &ClientConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
			Prefix:   "test:",
		},
		Universal: &UniversalConfig{
			Addresses:   []string{"localhost:6379"},
			Password:    "",
			DB:          0,
			Prefix:      "test:",
			ClusterMode: false,
		},
	}
}

func TestNewServiceProvider(t *testing.T) {
	provider := NewServiceProvider()
	assert.NotNil(t, provider, "Expected service provider to be initialized")
}

func TestServiceProviderRegister(t *testing.T) {
	t.Run("registers redis services to container with config", func(t *testing.T) {
		// Arrange
		container := di.New()
		mockConfig := mocks.NewMockManager(t)

		// Setup mock config to handle UnmarshalKey for our test Redis config
		testRedisConfig := setupTestRedisConfig()
		mockConfig.EXPECT().UnmarshalKey("redis", mock.Anything).Run(func(_ string, out interface{}) {
			// Copy our test config to the output parameter
			if cfg, ok := out.(*Config); ok {
				*cfg = *testRedisConfig
			}
		}).Return(nil)

		container.Instance("config", mockConfig)
		app := &mockApp{container: container}
		provider := NewServiceProvider()

		// Act - since we're testing dynamic providers, we have an empty list initially
		initialProviders := provider.Providers()
		assert.Empty(t, initialProviders, "Expected 0 initial providers")

		provider.Register(app)

		// Assert - check services were registered
		assert.True(t, container.Bound("redis.manager"), "Expected 'redis.manager' to be bound")
		assert.True(t, container.Bound("redis.client"), "Expected 'redis.client' to be bound")
		assert.True(t, container.Bound("redis.universal"), "Expected 'redis.universal' to be bound")

		// Check that providers were dynamically added
		finalProviders := provider.Providers()
		assert.Len(t, finalProviders, 3, "Expected 3 providers after registration")

		// Test manager resolution
		managerService, err := container.Make("redis.manager")
		assert.NoError(t, err, "Expected no error when resolving redis.manager")
		assert.IsType(t, &manager{}, managerService, "Expected redis.manager to be of type *manager")

		// Test client resolution
		clientService, err := container.Make("redis.client")
		assert.NoError(t, err, "Expected no error when resolving redis.client")
		assert.IsType(t, &redis.Client{}, clientService, "Expected redis.client to be of type *redis.Client")

		// Test universal client resolution
		universalService, err := container.Make("redis.universal")
		assert.NoError(t, err, "Expected no error when resolving redis.universal")
		_, ok := universalService.(redis.UniversalClient)
		assert.True(t, ok, "Expected redis.universal to be of type redis.UniversalClient")
	})

	t.Run("panics when config service is missing", func(t *testing.T) {
		// Arrange
		container := di.New()
		app := &mockApp{container: container}
		provider := NewServiceProvider()

		// Act & Assert - should panic when config is missing
		assert.Panics(t, func() {
			provider.Register(app)
		}, "Expected provider.Register to panic when config is missing")
	})

	t.Run("does nothing when app doesn't have container", func(t *testing.T) {
		// Arrange
		app := &mockApp{container: nil}
		provider := NewServiceProvider()

		// Act & Assert - should not panic
		assert.NotPanics(t, func() {
			provider.Register(app)
		}, "Should not panic when app doesn't have container")
	})
}

func TestServiceProviderBoot(t *testing.T) {
	t.Run("Boot doesn't panic", func(t *testing.T) {
		// Create DI container with config
		container := di.New()
		mockConfig := mocks.NewMockManager(t)

		// Setup expectations for UnmarshalKey
		mockConfig.EXPECT().UnmarshalKey("redis", mock.Anything).Run(func(_ string, out interface{}) {
			// Copy our test config to the output parameter
			if cfg, ok := out.(*Config); ok {
				*cfg = *setupTestRedisConfig()
			}
		}).Return(nil)

		container.Instance("config", mockConfig)

		// Create app and provider
		app := &mockApp{container: container}
		provider := NewServiceProvider()

		// First register the provider
		provider.Register(app)

		// Then test that boot doesn't panic
		assert.NotPanics(t, func() {
			provider.Boot(app)
		}, "Boot should not panic with valid configuration")
	})

	t.Run("Boot works without container", func(t *testing.T) {
		// Test with no container
		provider := NewServiceProvider()
		app := &mockApp{container: nil}

		// Should not panic
		assert.NotPanics(t, func() {
			provider.Boot(app)
		}, "Boot should not panic with nil container")
	})
}

func TestServiceProviderBootWithNil(t *testing.T) {
	// Test Boot with nil app parameter
	provider := NewServiceProvider()

	// Should not panic with nil app
	assert.NotPanics(t, func() {
		provider.Boot(nil)
	}, "Boot should not panic with nil app parameter")
}

func TestServiceProviderProviders(t *testing.T) {
	// In the new implementation, providers are dynamically added during Register
	// So a freshly created provider should have an empty providers list
	provider := NewServiceProvider()
	providers := provider.Providers()

	assert.Empty(t, providers, "Expected empty providers list initially")

	// We test the dynamic registration of providers in TestServiceProviderRegister
}

func TestServiceProviderRequires(t *testing.T) {
	provider := NewServiceProvider()
	requires := provider.Requires()

	// Redis provider requires the config provider
	assert.Len(t, requires, 1, "Expected 1 required dependency")
	assert.Equal(t, "config", requires[0], "Expected required dependency to be 'config'")
}

func TestDynamicProvidersList(t *testing.T) {
	// This test verifies that providers are correctly registered in the dynamic list
	container := di.New()
	mockConfig := mocks.NewMockManager(t)

	// Setup expectations for UnmarshalKey
	mockConfig.EXPECT().UnmarshalKey("redis", mock.Anything).Run(func(_ string, out interface{}) {
		// Copy our test config to the output parameter
		if cfg, ok := out.(*Config); ok {
			*cfg = *setupTestRedisConfig()
		}
	}).Return(nil)

	container.Instance("config", mockConfig)
	app := &mockApp{container: container}
	provider := NewServiceProvider()

	// Initially empty providers list
	initialProviders := provider.Providers()
	assert.Empty(t, initialProviders, "Expected 0 initial providers")

	// Register provider
	provider.Register(app)

	// Check providers list after registration
	providers := provider.Providers()

	// We expect 3 entries: redis, redis.client, redis.universal
	expectedItems := []string{"redis", "redis.client", "redis.universal"}
	for _, expected := range expectedItems {
		assert.Contains(t, providers, expected, "Expected to find '%s' in providers list", expected)
	}

	// Length should match too
	assert.Len(t, providers, len(expectedItems), "Expected %d providers", len(expectedItems))
}

func TestServiceProviderInterfaceCompliance(t *testing.T) {
	// This test verifies that our concrete type implements the interface
	var _ ServiceProvider = (*serviceProvider)(nil)
	var _ di.ServiceProvider = (*serviceProvider)(nil)
}

func TestMockConfigManagerWithRedisConfig(t *testing.T) {
	// This test verifies that our mock config manager can be used with Redis config
	mockConfig := mocks.NewMockManager(t)
	testConfig := setupTestRedisConfig()

	// Setup expectations for the Has method
	mockConfig.EXPECT().Has("redis").Return(true)

	// Setup expectations for the Get method
	mockConfig.EXPECT().Get("redis").Return(testConfig, true)

	// Setup expectations for UnmarshalKey
	mockConfig.EXPECT().UnmarshalKey("redis", mock.Anything).Run(func(_ string, out interface{}) {
		// Copy our test config to the output parameter
		if cfg, ok := out.(*Config); ok {
			*cfg = *testConfig
		}
	}).Return(nil)

	// Test Has method
	assert.True(t, mockConfig.Has("redis"), "Has should return true for redis key")

	// Test Get method
	value, exists := mockConfig.Get("redis")
	assert.True(t, exists, "Should find the redis key")
	assert.Equal(t, testConfig, value, "Should return our test config")

	// Test UnmarshalKey method
	var outConfig Config
	err := mockConfig.UnmarshalKey("redis", &outConfig)
	assert.NoError(t, err, "UnmarshalKey should not return an error")

	// Verify our mock expectations were met
	mockConfig.AssertExpectations(t)
}
