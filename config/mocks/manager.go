// Package mocks provides mock implementations of interfaces from the config package for testing purposes.
package mocks

import (
	"io"
	"time"

	"github.com/fsnotify/fsnotify"
)

// MockManager is a mock implementation of the Manager interface for testing purposes.
type MockManager struct {
	// Internal state for mocking
	hasKey       bool
	config       map[string]interface{}
	unmarshalErr error
}

// NewMockManager creates a new instance of MockManager with default configuration.
func NewMockManager() *MockManager {
	return &MockManager{
		config: make(map[string]interface{}),
	}
}

// SetHasKey sets the value to be returned by the Has method.
func (m *MockManager) SetHasKey(value bool) {
	m.hasKey = value
}

// SetConfig sets the internal configuration map.
func (m *MockManager) SetConfig(config map[string]interface{}) {
	m.config = config
}

// SetUnmarshalError sets the error to be returned by Unmarshal methods.
func (m *MockManager) SetUnmarshalError(err error) {
	m.unmarshalErr = err
}

// GetString returns a string value for the given key.
func (m *MockManager) GetString(key string) (string, bool) {
	if val, ok := m.config[key]; ok {
		if strVal, typeOk := val.(string); typeOk {
			return strVal, true
		}
	}
	return "", false
}

// GetInt returns an int value for the given key.
func (m *MockManager) GetInt(key string) (int, bool) {
	if val, ok := m.config[key]; ok {
		if intVal, typeOk := val.(int); typeOk {
			return intVal, true
		}
	}
	return 0, false
}

// GetBool returns a bool value for the given key.
func (m *MockManager) GetBool(key string) (bool, bool) {
	if val, ok := m.config[key]; ok {
		if boolVal, typeOk := val.(bool); typeOk {
			return boolVal, true
		}
	}
	return false, false
}

// GetFloat returns a float64 value for the given key.
func (m *MockManager) GetFloat(key string) (float64, bool) {
	if val, ok := m.config[key]; ok {
		if floatVal, typeOk := val.(float64); typeOk {
			return floatVal, true
		}
	}
	return 0, false
}

// GetDuration returns a time.Duration value for the given key.
func (m *MockManager) GetDuration(key string) (time.Duration, bool) {
	if val, ok := m.config[key]; ok {
		if durVal, typeOk := val.(time.Duration); typeOk {
			return durVal, true
		}
	}
	return 0, false
}

// GetTime returns a time.Time value for the given key.
func (m *MockManager) GetTime(key string) (time.Time, bool) {
	if val, ok := m.config[key]; ok {
		if timeVal, typeOk := val.(time.Time); typeOk {
			return timeVal, true
		}
	}
	return time.Time{}, false
}

// GetSlice returns a slice value for the given key.
func (m *MockManager) GetSlice(key string) ([]interface{}, bool) {
	if val, ok := m.config[key]; ok {
		if sliceVal, typeOk := val.([]interface{}); typeOk {
			return sliceVal, true
		}
	}
	return nil, false
}

// GetStringSlice returns a string slice value for the given key.
func (m *MockManager) GetStringSlice(key string) ([]string, bool) {
	if val, ok := m.config[key]; ok {
		if sliceVal, typeOk := val.([]string); typeOk {
			return sliceVal, true
		}
	}
	return nil, false
}

// GetIntSlice returns an int slice value for the given key.
func (m *MockManager) GetIntSlice(key string) ([]int, bool) {
	if val, ok := m.config[key]; ok {
		if sliceVal, typeOk := val.([]int); typeOk {
			return sliceVal, true
		}
	}
	return nil, false
}

// GetMap returns a map value for the given key.
func (m *MockManager) GetMap(key string) (map[string]interface{}, bool) {
	if val, ok := m.config[key]; ok {
		if mapVal, typeOk := val.(map[string]interface{}); typeOk {
			return mapVal, true
		}
	}
	return nil, false
}

// GetStringMap returns a map[string]interface{} value for the given key.
func (m *MockManager) GetStringMap(key string) (map[string]interface{}, bool) {
	if val, ok := m.config[key]; ok {
		if mapVal, typeOk := val.(map[string]interface{}); typeOk {
			return mapVal, true
		}
	}
	return nil, false
}

// GetStringMapString returns a map[string]string value for the given key.
func (m *MockManager) GetStringMapString(key string) (map[string]string, bool) {
	if val, ok := m.config[key]; ok {
		if mapVal, typeOk := val.(map[string]string); typeOk {
			return mapVal, true
		}
	}
	return nil, false
}

// GetStringMapStringSlice returns a map[string][]string value for the given key.
func (m *MockManager) GetStringMapStringSlice(key string) (map[string][]string, bool) {
	if val, ok := m.config[key]; ok {
		if mapVal, typeOk := val.(map[string][]string); typeOk {
			return mapVal, true
		}
	}
	return nil, false
}

// Get returns a raw interface{} value for the given key.
func (m *MockManager) Get(key string) (interface{}, bool) {
	val, ok := m.config[key]
	return val, ok
}

// Set sets a value for the given key.
func (m *MockManager) Set(key string, value interface{}) error {
	m.config[key] = value
	return nil
}

// SetDefault sets a default value for the given key.
func (m *MockManager) SetDefault(key string, value interface{}) {
	if _, exists := m.config[key]; !exists {
		m.config[key] = value
	}
}

// Has checks if the key exists in the configuration.
func (m *MockManager) Has(key string) bool {
	_, exists := m.config[key]
	return exists || m.hasKey
}

// AllSettings returns all settings as a map[string]interface{}.
func (m *MockManager) AllSettings() map[string]interface{} {
	return m.config
}

// AllKeys returns all keys for the configuration.
func (m *MockManager) AllKeys() []string {
	keys := make([]string, 0, len(m.config))
	for k := range m.config {
		keys = append(keys, k)
	}
	return keys
}

// Unmarshal unmarshals the config into a struct.
func (m *MockManager) Unmarshal(target interface{}) error {
	return m.unmarshalErr
}

// UnmarshalKey unmarshals a specific key into a struct.
func (m *MockManager) UnmarshalKey(key string, target interface{}) error {
	return m.unmarshalErr
}

// SetConfigFile sets the path to a config file.
func (m *MockManager) SetConfigFile(path string) {
	// Mock implementation, do nothing
}

// SetConfigType sets the type of the configuration.
func (m *MockManager) SetConfigType(configType string) {
	// Mock implementation, do nothing
}

// SetConfigName sets the name for the config file.
func (m *MockManager) SetConfigName(name string) {
	// Mock implementation, do nothing
}

// AddConfigPath adds a path to search for the config file in.
func (m *MockManager) AddConfigPath(path string) {
	// Mock implementation, do nothing
}

// ReadInConfig loads the config file.
func (m *MockManager) ReadInConfig() error {
	return nil
}

// MergeInConfig merges a new config file with the current config.
func (m *MockManager) MergeInConfig() error {
	return nil
}

// WriteConfig writes the current configuration to a file.
func (m *MockManager) WriteConfig() error {
	return nil
}

// SafeWriteConfig writes the current configuration to a file if it doesn't exist.
func (m *MockManager) SafeWriteConfig() error {
	return nil
}

// WriteConfigAs writes the current configuration to a file with the specified name.
func (m *MockManager) WriteConfigAs(filename string) error {
	return nil
}

// SafeWriteConfigAs writes the current configuration to a file with the specified name if it doesn't exist.
func (m *MockManager) SafeWriteConfigAs(filename string) error {
	return nil
}

// WatchConfig watches the config file and reloads it when it changes.
func (m *MockManager) WatchConfig() {
	// Mock implementation, do nothing
}

// OnConfigChange sets the function to run when a config change is detected.
func (m *MockManager) OnConfigChange(callback func(event fsnotify.Event)) {
	// Mock implementation, do nothing
}

// SetEnvPrefix sets the prefix that environment variables will use.
func (m *MockManager) SetEnvPrefix(prefix string) {
	// Mock implementation, do nothing
}

// AutomaticEnv enables automatic environment variable support.
func (m *MockManager) AutomaticEnv() {
	// Mock implementation, do nothing
}

// BindEnv binds a config key to an environment variable.
func (m *MockManager) BindEnv(input ...string) error {
	return nil
}

// MergeConfig merges a new config with the current config.
func (m *MockManager) MergeConfig(in io.Reader) error {
	return nil
}
