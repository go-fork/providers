package config

import (
	"testing"

	"github.com/go-fork/di"
	"github.com/stretchr/testify/assert"
)

// mockApp implements the Container() method required by ServiceProvider
type mockApp struct {
	container *di.Container
}

func (m *mockApp) Container() *di.Container {
	return m.container
}

// mockAppInvalid doesn't implement Container()
type mockAppInvalid struct{}

func TestNewServiceProvider(t *testing.T) {
	provider := NewServiceProvider()
	assert.NotNil(t, provider, "NewServiceProvider should return a non-nil provider")
}

func TestServiceProvider_Register(t *testing.T) {
	provider := NewServiceProvider()

	// Test with valid app
	diContainer := di.New()
	app := &mockApp{container: diContainer}

	provider.Register(app)

	// Check that config was registered in the container
	cfg, err := diContainer.Make("config")
	assert.NoError(t, err)
	assert.NotNil(t, cfg, "Config should be registered in container")

	// Test with app that doesn't implement Container()
	invalidApp := &mockAppInvalid{}
	// This should not panic
	provider.Register(invalidApp)

	// Test with nil container (should panic)
	app.container = nil
	assert.Panics(t, func() {
		provider.Register(app)
	}, "Register should panic when container is nil")
}

func TestServiceProvider_Boot(t *testing.T) {
	provider := NewServiceProvider()

	// Test with valid app
	app := &mockApp{container: di.New()}
	assert.NotPanics(t, func() {
		provider.Boot(app)
	}, "Boot should not panic with valid app")

	// Test with invalid app type
	invalidApp := &mockAppInvalid{}
	assert.NotPanics(t, func() {
		provider.Boot(invalidApp)
	}, "Boot should not panic with invalid app")

	// Test with nil app
	assert.NotPanics(t, func() {
		provider.Boot(nil)
	}, "Boot should not panic with nil app")
}
