package cache

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/di"
	dimocks "go.fork.vn/di/mocks"
	"go.fork.vn/providers/cache/config"
	configmocks "go.fork.vn/providers/config/mocks"
	redismocks "go.fork.vn/providers/redis/mocks"
)

func TestServiceProviderRegister(t *testing.T) {
	t.Run("registers cache manager to container", func(t *testing.T) {
		// Arrange
		container := di.New()
		mockApp := &dimocks.Application{}
		mockConfigManager := configmocks.NewMockManager(t)

		// Setup app mock - old style mocking
		mockApp.On("Container").Return(container)

		// Setup config manager
		container.Instance("config", mockConfigManager)

		// Configure cache config
		cacheConfig := config.Config{
			DefaultDriver: "memory",
			Drivers: config.DriversConfig{
				Memory: &config.DriverMemoryConfig{
					Enabled:         true,
					DefaultTTL:      3600,
					CleanupInterval: 600,
					MaxItems:        1000,
				},
			},
		}

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(func(key string, cfg interface{}) error {
			if c, ok := cfg.(*config.Config); ok {
				*c = cacheConfig
			}
			return nil
		})

		provider := NewServiceProvider()

		// Act
		provider.Register(mockApp)

		// Assert
		assert.True(t, container.Bound("cache"), "Expected 'cache' to be bound in container")

		cacheService, err := container.Make("cache")
		assert.NoError(t, err, "Expected no error when making 'cache'")
		assert.IsType(t, &manager{}, cacheService, "Expected cache service to be *manager")

		// Test that the manager has correct default driver
		cacheManager := cacheService.(*manager)
		assert.Equal(t, "memory", cacheManager.defaultDriver, "Expected default driver to be 'memory'")
	})

	t.Run("registers with file driver configuration", func(t *testing.T) {
		// Arrange
		container := di.New()
		mockApp := &dimocks.Application{}
		mockConfigManager := configmocks.NewMockManager(t)

		mockApp.On("Container").Return(container)
		container.Instance("config", mockConfigManager)

		cacheConfig := config.Config{
			DefaultDriver: "file",
			Drivers: config.DriversConfig{
				File: &config.DriverFileConfig{
					Enabled:         true,
					Path:            "/tmp/cache",
					DefaultTTL:      3600,
					CleanupInterval: 600,
				},
			},
		}

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(func(key string, cfg interface{}) error {
			if c, ok := cfg.(*config.Config); ok {
				*c = cacheConfig
			}
			return nil
		})

		provider := NewServiceProvider()

		// Act
		provider.Register(mockApp)

		// Assert
		assert.True(t, container.Bound("cache"), "Expected 'cache' to be bound in container")
	})

	t.Run("registers with redis driver configuration", func(t *testing.T) {
		// Arrange
		container := di.New()
		mockApp := &dimocks.Application{}
		mockConfigManager := configmocks.NewMockManager(t)
		mockRedisManager := redismocks.NewMockManager(t)

		mockApp.On("Container").Return(container)
		container.Instance("config", mockConfigManager)

		// Mock Redis manager and bind it to container
		mockRedisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
		mockRedisManager.EXPECT().Client().Return(mockRedisClient, nil)
		container.Instance("redis", mockRedisManager)

		cacheConfig := config.Config{
			DefaultDriver: "redis",
			Drivers: config.DriversConfig{
				Redis: &config.DriverRedisConfig{
					Enabled:    true,
					DefaultTTL: 3600,
					Serializer: "json",
				},
			},
		}

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(func(key string, cfg interface{}) error {
			if c, ok := cfg.(*config.Config); ok {
				*c = cacheConfig
			}
			return nil
		})

		provider := NewServiceProvider()

		// Act
		provider.Register(mockApp)

		// Assert
		assert.True(t, container.Bound("cache"), "Expected 'cache' to be bound in container")
	})

	t.Run("registers with mongodb driver configuration", func(t *testing.T) {
		// Skip this test as it requires complex MongoDB driver mocking
		// The MongoDB driver constructor calls DatabaseWithName().Collection() which
		// requires proper MongoDB client setup that's beyond the scope of this unit test
		t.Skip("MongoDB driver test requires integration testing with real MongoDB client setup")
	})
}

func TestServiceProviderBoot(t *testing.T) {
	t.Run("boots successfully with valid configuration", func(t *testing.T) {
		// Arrange
		container := di.New()
		mockApp := &dimocks.Application{}
		mockConfigManager := configmocks.NewMockManager(t)

		mockApp.On("Container").Return(container)
		container.Instance("config", mockConfigManager)

		cacheConfig := config.Config{
			DefaultDriver: "memory",
			Drivers: config.DriversConfig{
				Memory: &config.DriverMemoryConfig{
					Enabled:         true,
					DefaultTTL:      3600,
					CleanupInterval: 600,
					MaxItems:        1000,
				},
			},
		}

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(func(key string, cfg interface{}) error {
			if c, ok := cfg.(*config.Config); ok {
				*c = cacheConfig
			}
			return nil
		})

		provider := NewServiceProvider()
		provider.Register(mockApp)

		// Act
		provider.Boot(mockApp)

		// Assert
		// If no panic occurred, the boot was successful
		assert.True(t, true, "Boot should complete without errors")
	})
}

func TestServiceProviderHandlesConfigurationErrors(t *testing.T) {
	t.Run("panics when config unmarshal errors occur", func(t *testing.T) {
		// Arrange
		container := di.New()
		mockApp := &dimocks.Application{}
		mockConfigManager := configmocks.NewMockManager(t)

		mockApp.On("Container").Return(container)
		container.Instance("config", mockConfigManager)

		// Configure config manager to return an error
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).Return(assert.AnError)

		provider := NewServiceProvider()

		// Act & Assert
		// Should panic with config errors (this is the expected behavior)
		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Register should panic on config errors")
	})
}

func TestServiceProviderWithNilApp(t *testing.T) {
	t.Run("handles nil app gracefully", func(t *testing.T) {
		// Arrange
		provider := NewServiceProvider()

		// Act & Assert
		assert.NotPanics(t, func() {
			provider.Register(nil)
		}, "Register should not panic with nil app")

		assert.NotPanics(t, func() {
			provider.Boot(nil)
		}, "Boot should not panic with nil app")
	})
}

func TestServiceProviderInterface(t *testing.T) {
	t.Run("implements service provider interface", func(t *testing.T) {
		// Arrange & Act
		provider := NewServiceProvider()

		// Assert
		assert.Implements(t, (*di.ServiceProvider)(nil), provider, "Should implement ServiceProvider interface")
	})
}

func TestServiceProviderWithMissingContainer(t *testing.T) {
	t.Run("panics when container is nil", func(t *testing.T) {
		// Arrange
		mockApp := &dimocks.Application{}
		mockApp.On("Container").Return(nil)

		provider := NewServiceProvider()

		// Act & Assert
		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Register should panic with nil container")
	})
}

func TestServiceProviderCacheManagerCreation(t *testing.T) {
	t.Run("creates cache manager with all driver types", func(t *testing.T) {
		// Arrange
		container := di.New()
		mockApp := &dimocks.Application{}
		mockConfigManager := configmocks.NewMockManager(t)
		mockRedisManager := redismocks.NewMockManager(t)

		mockApp.On("Container").Return(container)
		container.Instance("config", mockConfigManager)

		// Mock Redis manager and bind it to container - needed for Redis driver
		mockRedisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
		mockRedisManager.EXPECT().Client().Return(mockRedisClient, nil)
		container.Instance("redis", mockRedisManager)

		// Test configuration with all drivers (MongoDB disabled)
		cacheConfig := config.Config{
			DefaultDriver: "memory",
			DefaultTTL:    1800,
			Prefix:        "test_",
			Drivers: config.DriversConfig{
				Memory: &config.DriverMemoryConfig{
					Enabled:         true,
					DefaultTTL:      3600,
					CleanupInterval: 600,
					MaxItems:        1000,
				},
				File: &config.DriverFileConfig{
					Enabled:         true,
					Path:            "/tmp/cache",
					DefaultTTL:      3600,
					CleanupInterval: 600,
				},
				Redis: &config.DriverRedisConfig{
					Enabled:    true,
					DefaultTTL: 3600,
					Serializer: "json",
				},
				MongoDB: &config.DriverMongodbConfig{
					Enabled:    false, // Disable MongoDB to avoid complexity
					DefaultTTL: 3600,
				},
			},
		}

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(func(key string, cfg interface{}) error {
			if c, ok := cfg.(*config.Config); ok {
				*c = cacheConfig
			}
			return nil
		})

		provider := NewServiceProvider().(*serviceProvider)

		// Act
		provider.Register(mockApp)

		// Assert
		cacheService, err := container.Make("cache")
		assert.NoError(t, err)

		cacheManager := cacheService.(*manager)
		// Verify drivers are configured
		assert.Contains(t, cacheManager.drivers, "memory")
		assert.Contains(t, cacheManager.drivers, "file")
		assert.Contains(t, cacheManager.drivers, "redis")
		// MongoDB intentionally not checked (disabled)

		// Test that we can get a driver instance (verifies the driver type)
		memDriver, err := cacheManager.Driver("memory")
		assert.NoError(t, err)
		assert.NotNil(t, memDriver)
		fileDriver, err := cacheManager.Driver("file")
		assert.NoError(t, err)
		assert.NotNil(t, fileDriver)
		redisDriver, err := cacheManager.Driver("redis")
		assert.NoError(t, err)
		assert.NotNil(t, redisDriver)
	})
}

func TestServiceProviderEdgeCases(t *testing.T) {
	t.Run("handles empty configuration", func(t *testing.T) {
		// Arrange
		container := di.New()
		mockApp := &dimocks.Application{}
		mockConfigManager := configmocks.NewMockManager(t)

		mockApp.On("Container").Return(container)
		container.Instance("config", mockConfigManager)

		// Empty configuration
		emptyConfig := config.Config{}

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(func(key string, cfg interface{}) error {
			if c, ok := cfg.(*config.Config); ok {
				*c = emptyConfig
			}
			return nil
		})

		provider := NewServiceProvider()

		// Act & Assert
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		}, "Should handle empty configuration gracefully")
	})

	t.Run("handles missing config service", func(t *testing.T) {
		// Arrange
		container := di.New()
		mockApp := &dimocks.Application{}

		mockApp.On("Container").Return(container)
		// Don't bind config service to container

		provider := NewServiceProvider()

		// Act & Assert
		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Should panic when missing config service")
	})
}
