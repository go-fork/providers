package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/di"
	"go.fork.vn/providers/config/mocks"
)

// mockApp implements the container interface for testing
type mockApp struct {
	container *di.Container
}

func (a *mockApp) Container() *di.Container {
	return a.container
}

// setupTestMongoConfig creates a MongoDB config for testing
func setupTestMongoConfig() *Config {
	return &Config{
		URI:                    "mongodb://localhost:27017",
		Database:               "testdb",
		ConnectTimeout:         10000,
		MaxPoolSize:            10,
		MinPoolSize:            1,
		MaxConnIdleTime:        300000,
		HeartbeatInterval:      30000,
		ServerSelectionTimeout: 30000,
		SocketTimeout:          5000,
		LocalThreshold:         15000,
		Auth: AuthConfig{
			Username:      "",
			Password:      "",
			AuthSource:    "",
			AuthMechanism: "",
		},
		TLS: TLSConfig{
			InsecureSkipVerify: false,
		},
		ReadPreference: ReadPreferenceConfig{
			Mode: "primary",
		},
		ReadConcern: ReadConcernConfig{
			Level: "",
		},
		WriteConcern: WriteConcernConfig{
			W:        1,
			WTimeout: 0,
			Journal:  false,
		},
		AppName:      "",
		Direct:       false,
		ReplicaSet:   "",
		Compressors:  []string{},
		RetryWrites:  true,
		RetryReads:   true,
		LoadBalanced: false,
	}
}

func TestNewServiceProvider(t *testing.T) {
	provider := NewServiceProvider()
	assert.NotNil(t, provider, "Expected service provider to be initialized")
}

func TestServiceProviderRegister(t *testing.T) {
	t.Run("registers mongodb services to container with config", func(t *testing.T) {
		// Arrange
		container := di.New()
		mockConfig := mocks.NewMockManager(t)

		// Setup mock config to handle UnmarshalKey for our test MongoDB config
		testMongoConfig := setupTestMongoConfig()
		mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
			// Copy our test config to the output parameter
			if cfg, ok := out.(*Config); ok {
				*cfg = *testMongoConfig
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
		assert.True(t, container.Bound("mongodb.manager"), "Expected 'mongodb.manager' to be bound")
		assert.True(t, container.Bound("mongodb.client"), "Expected 'mongodb.client' to be bound")
		assert.True(t, container.Bound("mongodb.database"), "Expected 'mongodb.database' to be bound")

		// Check that providers were dynamically added
		finalProviders := provider.Providers()
		assert.Len(t, finalProviders, 3, "Expected 3 providers after registration")

		// Test manager resolution
		managerService, err := container.Make("mongodb.manager")
		assert.NoError(t, err, "Expected no error when resolving mongodb.manager")
		assert.NotNil(t, managerService, "Expected mongodb.manager to be non-nil")

		// Test client resolution
		clientService, err := container.Make("mongodb.client")
		assert.NoError(t, err, "Expected no error when resolving mongodb.client")
		assert.NotNil(t, clientService, "Expected mongodb.client to be non-nil")

		// Test database resolution
		databaseService, err := container.Make("mongodb.database")
		assert.NoError(t, err, "Expected no error when resolving mongodb.database")
		assert.NotNil(t, databaseService, "Expected mongodb.database to be non-nil")
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
		mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
			// Copy our test config to the output parameter
			if cfg, ok := out.(*Config); ok {
				*cfg = *setupTestMongoConfig()
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

	// MongoDB provider requires the config provider
	assert.Len(t, requires, 1, "Expected 1 required dependency")
	assert.Equal(t, "config", requires[0], "Expected required dependency to be 'config'")
}

func TestDynamicProvidersList(t *testing.T) {
	// This test verifies that providers are correctly registered in the dynamic list
	container := di.New()
	mockConfig := mocks.NewMockManager(t)

	// Setup expectations for UnmarshalKey
	mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
		// Copy our test config to the output parameter
		if cfg, ok := out.(*Config); ok {
			*cfg = *setupTestMongoConfig()
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

	// We expect 3 entries: mongodb, mongodb.client, mongodb.database
	expectedItems := []string{"mongodb", "mongodb.client", "mongodb.database"}
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

func TestMockConfigManagerWithMongoConfig(t *testing.T) {
	// This test verifies that our mock config manager can be used with MongoDB config
	mockConfig := mocks.NewMockManager(t)
	testConfig := setupTestMongoConfig()

	// Setup expectations for the Has method
	mockConfig.EXPECT().Has("mongodb").Return(true)

	// Setup expectations for the Get method
	mockConfig.EXPECT().Get("mongodb").Return(testConfig, true)

	// Setup expectations for UnmarshalKey
	mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
		// Copy our test config to the output parameter
		if cfg, ok := out.(*Config); ok {
			*cfg = *testConfig
		}
	}).Return(nil)

	// Test Has method
	assert.True(t, mockConfig.Has("mongodb"), "Has should return true for mongodb key")

	// Test Get method
	value, exists := mockConfig.Get("mongodb")
	assert.True(t, exists, "Should find the mongodb key")
	assert.Equal(t, testConfig, value, "Should return our test config")

	// Test UnmarshalKey method
	var outConfig Config
	err := mockConfig.UnmarshalKey("mongodb", &outConfig)
	assert.NoError(t, err, "UnmarshalKey should not return an error")

	// Verify our mock expectations were met
	mockConfig.AssertExpectations(t)
}
