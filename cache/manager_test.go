package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-fork/providers/cache/mocks"
)

func TestManagerGet(t *testing.T) {
	t.Run("returns value when default driver is set", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		mockDriver.GetFunc = func(ctx context.Context, key string) (interface{}, bool) {
			if key == "test-key" {
				return "test-value", true
			}
			return nil, false
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, found := manager.Get("test-key")

		// Assert
		if !found {
			t.Errorf("Expected to find key, but didn't")
		}
		if value != "test-value" {
			t.Errorf("Expected value to be %v, got %v", "test-value", value)
		}
	})

	t.Run("returns not found when default driver is set but key doesn't exist", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		_, found := manager.Get("nonexistent-key")

		// Assert
		if found {
			t.Errorf("Expected not to find key, but did")
		}
	})

	t.Run("returns not found when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		_, found := manager.Get("test-key")

		// Assert
		if found {
			t.Errorf("Expected not to find key, but did")
		}
	})
}

func TestManagerSet(t *testing.T) {
	t.Run("sets value when default driver is set", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		var setKey string
		var setValue interface{}
		var setTTL time.Duration

		mockDriver.SetFunc = func(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
			setKey = key
			setValue = value
			setTTL = ttl
			return nil
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if setKey != "test-key" {
			t.Errorf("Expected key to be %v, got %v", "test-key", setKey)
		}
		if setValue != "test-value" {
			t.Errorf("Expected value to be %v, got %v", "test-value", setValue)
		}
		if setTTL != 5*time.Minute {
			t.Errorf("Expected TTL to be %v, got %v", 5*time.Minute, setTTL)
		}
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("returns error when driver set fails", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		expectedErr := errors.New("set error")
		mockDriver.SetFunc = func(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
			return expectedErr
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestManagerHas(t *testing.T) {
	t.Run("returns true when key exists", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		mockDriver.HasFunc = func(ctx context.Context, key string) bool {
			return key == "test-key"
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		exists := manager.Has("test-key")

		// Assert
		if !exists {
			t.Errorf("Expected key to exist, but it doesn't")
		}
	})

	t.Run("returns false when key doesn't exist", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		mockDriver.HasFunc = func(ctx context.Context, key string) bool {
			return false
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		exists := manager.Has("nonexistent-key")

		// Assert
		if exists {
			t.Errorf("Expected key not to exist, but it does")
		}
	})

	t.Run("returns false when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		exists := manager.Has("test-key")

		// Assert
		if exists {
			t.Errorf("Expected key not to exist, but it does")
		}
	})
}

func TestManagerDelete(t *testing.T) {
	t.Run("deletes key successfully", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		var deletedKey string
		mockDriver.DeleteFunc = func(ctx context.Context, key string) error {
			deletedKey = key
			return nil
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Delete("test-key")

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if deletedKey != "test-key" {
			t.Errorf("Expected deleted key to be %v, got %v", "test-key", deletedKey)
		}
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		err := manager.Delete("test-key")

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("returns error when driver delete fails", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		expectedErr := errors.New("delete error")
		mockDriver.DeleteFunc = func(ctx context.Context, key string) error {
			return expectedErr
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Delete("test-key")

		// Assert
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestManagerFlush(t *testing.T) {
	t.Run("flushes cache successfully", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		flushed := false
		mockDriver.FlushFunc = func(ctx context.Context) error {
			flushed = true
			return nil
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Flush()

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !flushed {
			t.Errorf("Expected cache to be flushed, but it wasn't")
		}
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		err := manager.Flush()

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("returns error when driver flush fails", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		expectedErr := errors.New("flush error")
		mockDriver.FlushFunc = func(ctx context.Context) error {
			return expectedErr
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Flush()

		// Assert
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestManagerGetMultiple(t *testing.T) {
	t.Run("gets multiple values successfully", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		expectedValues := map[string]interface{}{
			"key1": "value1",
			"key3": "value3",
		}
		expectedMissing := []string{"key2"}

		mockDriver.GetMultipleFunc = func(ctx context.Context, keys []string) (map[string]interface{}, []string) {
			return expectedValues, expectedMissing
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		values, missing := manager.GetMultiple([]string{"key1", "key2", "key3"})

		// Assert
		if len(values) != len(expectedValues) {
			t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
		}
		for k, v := range expectedValues {
			if values[k] != v {
				t.Errorf("Expected value for key %s to be %v, got %v", k, v, values[k])
			}
		}
		if len(missing) != len(expectedMissing) {
			t.Errorf("Expected %d missing keys, got %d", len(expectedMissing), len(missing))
		}
		for i, k := range expectedMissing {
			if missing[i] != k {
				t.Errorf("Expected missing key at index %d to be %s, got %s", i, k, missing[i])
			}
		}
	})

	t.Run("returns empty result when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		keys := []string{"key1", "key2", "key3"}

		// Act
		values, missing := manager.GetMultiple(keys)

		// Assert
		if len(values) != 0 {
			t.Errorf("Expected 0 values, got %d", len(values))
		}
		if len(missing) != len(keys) {
			t.Errorf("Expected %d missing keys, got %d", len(keys), len(missing))
		}
		for i, k := range keys {
			if missing[i] != k {
				t.Errorf("Expected missing key at index %d to be %s, got %s", i, k, missing[i])
			}
		}
	})
}

func TestManagerSetMultiple(t *testing.T) {
	t.Run("sets multiple values successfully", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		var setValues map[string]interface{}
		var setTTL time.Duration

		mockDriver.SetMultipleFunc = func(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
			setValues = values
			setTTL = ttl
			return nil
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		ttl := 10 * time.Minute

		// Act
		err := manager.SetMultiple(values, ttl)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(setValues) != len(values) {
			t.Errorf("Expected %d values, got %d", len(values), len(setValues))
		}
		for k, v := range values {
			if setValues[k] != v {
				t.Errorf("Expected value for key %s to be %v, got %v", k, v, setValues[k])
			}
		}
		if setTTL != ttl {
			t.Errorf("Expected TTL to be %v, got %v", ttl, setTTL)
		}
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}

		// Act
		err := manager.SetMultiple(values, 10*time.Minute)

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("returns error when driver setMultiple fails", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		expectedErr := errors.New("setMultiple error")
		mockDriver.SetMultipleFunc = func(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
			return expectedErr
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}

		// Act
		err := manager.SetMultiple(values, 10*time.Minute)

		// Assert
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestManagerDeleteMultiple(t *testing.T) {
	t.Run("deletes multiple keys successfully", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		var deletedKeys []string

		mockDriver.DeleteMultipleFunc = func(ctx context.Context, keys []string) error {
			deletedKeys = keys
			return nil
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		keys := []string{"key1", "key2", "key3"}

		// Act
		err := manager.DeleteMultiple(keys)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(deletedKeys) != len(keys) {
			t.Errorf("Expected %d keys, got %d", len(keys), len(deletedKeys))
		}
		for i, k := range keys {
			if deletedKeys[i] != k {
				t.Errorf("Expected key at index %d to be %s, got %s", i, k, deletedKeys[i])
			}
		}
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		keys := []string{"key1", "key2", "key3"}

		// Act
		err := manager.DeleteMultiple(keys)

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("returns error when driver deleteMultiple fails", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		expectedErr := errors.New("deleteMultiple error")
		mockDriver.DeleteMultipleFunc = func(ctx context.Context, keys []string) error {
			return expectedErr
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		keys := []string{"key1", "key2", "key3"}

		// Act
		err := manager.DeleteMultiple(keys)

		// Assert
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestManagerRemember(t *testing.T) {
	t.Run("remembers value successfully", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		expectedValue := "computed-value"
		callbackCalled := false

		mockDriver.RememberFunc = func(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
			callbackCalled = true
			return expectedValue, nil
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, err := manager.Remember("test-key", 5*time.Minute, func() (interface{}, error) {
			return expectedValue, nil
		})

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if value != expectedValue {
			t.Errorf("Expected value to be %v, got %v", expectedValue, value)
		}
		if !callbackCalled {
			t.Errorf("Expected callback to be called, but it wasn't")
		}
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		_, err := manager.Remember("test-key", 5*time.Minute, func() (interface{}, error) {
			return "value", nil
		})

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("returns error when driver remember fails", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver()
		expectedErr := errors.New("remember error")
		mockDriver.RememberFunc = func(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
			return nil, expectedErr
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		_, err := manager.Remember("test-key", 5*time.Minute, func() (interface{}, error) {
			return "value", nil
		})

		// Assert
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestManagerAddDriver(t *testing.T) {
	t.Run("adds driver successfully", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		mockDriver := mocks.NewMockDriver()

		// Act
		manager.AddDriver("mock", mockDriver)

		// Assert
		driver, err := manager.Driver("mock")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if driver != mockDriver {
			t.Errorf("Expected to get the same driver instance")
		}
	})

	t.Run("sets first driver as default", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		mockDriver := mocks.NewMockDriver()
		mockDriver.HasFunc = func(ctx context.Context, key string) bool {
			return key == "test-key"
		}

		// Act
		manager.AddDriver("mock", mockDriver)

		// Assert
		exists := manager.Has("test-key")
		if !exists {
			t.Errorf("Expected key to exist, but it doesn't")
		}
	})

	t.Run("doesn't change default when adding subsequent drivers", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		mockDriver1 := mocks.NewMockDriver()
		mockDriver1.HasFunc = func(ctx context.Context, key string) bool {
			return key == "test-key-1"
		}
		mockDriver2 := mocks.NewMockDriver()
		mockDriver2.HasFunc = func(ctx context.Context, key string) bool {
			return key == "test-key-2"
		}

		// Act
		manager.AddDriver("mock1", mockDriver1)
		manager.AddDriver("mock2", mockDriver2)

		// Assert
		exists1 := manager.Has("test-key-1")
		exists2 := manager.Has("test-key-2")
		if !exists1 {
			t.Errorf("Expected key1 to exist, but it doesn't")
		}
		if exists2 {
			t.Errorf("Expected key2 not to exist, but it does")
		}
	})
}

func TestManagerSetDefaultDriver(t *testing.T) {
	t.Run("sets default driver successfully", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		mockDriver1 := mocks.NewMockDriver()
		mockDriver1.HasFunc = func(ctx context.Context, key string) bool {
			return key == "test-key-1"
		}
		mockDriver2 := mocks.NewMockDriver()
		mockDriver2.HasFunc = func(ctx context.Context, key string) bool {
			return key == "test-key-2"
		}

		manager.AddDriver("mock1", mockDriver1)
		manager.AddDriver("mock2", mockDriver2)

		// Act
		manager.SetDefaultDriver("mock2")

		// Assert
		exists1 := manager.Has("test-key-1")
		exists2 := manager.Has("test-key-2")
		if exists1 {
			t.Errorf("Expected key1 not to exist, but it does")
		}
		if !exists2 {
			t.Errorf("Expected key2 to exist, but it doesn't")
		}
	})

	t.Run("does nothing if driver doesn't exist", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		mockDriver := mocks.NewMockDriver()
		mockDriver.HasFunc = func(ctx context.Context, key string) bool {
			return key == "test-key"
		}

		manager.AddDriver("mock", mockDriver)

		// Act
		manager.SetDefaultDriver("nonexistent")

		// Assert
		exists := manager.Has("test-key")
		if !exists {
			t.Errorf("Expected key to exist, but it doesn't")
		}
	})
}

func TestManagerDriver(t *testing.T) {
	t.Run("returns driver successfully", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		mockDriver := mocks.NewMockDriver()
		manager.AddDriver("mock", mockDriver)

		// Act
		driver, err := manager.Driver("mock")

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if driver != mockDriver {
			t.Errorf("Expected to get the same driver instance")
		}
	})

	t.Run("returns error if driver doesn't exist", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		_, err := manager.Driver("nonexistent")

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestManagerStats(t *testing.T) {
	t.Run("returns stats for all drivers", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		mockDriver1 := mocks.NewMockDriver()
		mockDriver1Stats := map[string]interface{}{
			"items": 10,
			"hits":  100,
		}
		mockDriver1.StatsFunc = func(ctx context.Context) map[string]interface{} {
			return mockDriver1Stats
		}

		mockDriver2 := mocks.NewMockDriver()
		mockDriver2Stats := map[string]interface{}{
			"items": 20,
			"hits":  200,
		}
		mockDriver2.StatsFunc = func(ctx context.Context) map[string]interface{} {
			return mockDriver2Stats
		}

		manager.AddDriver("mock1", mockDriver1)
		manager.AddDriver("mock2", mockDriver2)

		// Act
		stats := manager.Stats()

		// Assert
		if len(stats) != 2 {
			t.Errorf("Expected stats for 2 drivers, got %d", len(stats))
		}

		if stats["mock1"]["items"] != mockDriver1Stats["items"] {
			t.Errorf("Expected mock1 items to be %v, got %v", mockDriver1Stats["items"], stats["mock1"]["items"])
		}

		if stats["mock2"]["hits"] != mockDriver2Stats["hits"] {
			t.Errorf("Expected mock2 hits to be %v, got %v", mockDriver2Stats["hits"], stats["mock2"]["hits"])
		}
	})

	t.Run("returns empty map when no drivers", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		stats := manager.Stats()

		// Assert
		if len(stats) != 0 {
			t.Errorf("Expected empty stats, got %v", stats)
		}
	})
}

func TestManagerClose(t *testing.T) {
	t.Run("closes all drivers successfully", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		mockDriver1 := mocks.NewMockDriver()
		mockDriver1Closed := false
		mockDriver1.CloseFunc = func() error {
			mockDriver1Closed = true
			return nil
		}

		mockDriver2 := mocks.NewMockDriver()
		mockDriver2Closed := false
		mockDriver2.CloseFunc = func() error {
			mockDriver2Closed = true
			return nil
		}

		manager.AddDriver("mock1", mockDriver1)
		manager.AddDriver("mock2", mockDriver2)

		// Act
		err := manager.Close()

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !mockDriver1Closed {
			t.Errorf("Expected mock1 to be closed, but it wasn't")
		}
		if !mockDriver2Closed {
			t.Errorf("Expected mock2 to be closed, but it wasn't")
		}
	})

	t.Run("returns error if any driver fails to close", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		mockDriver1 := mocks.NewMockDriver()
		mockDriver1.CloseFunc = func() error {
			return nil
		}

		mockDriver2 := mocks.NewMockDriver()
		expectedErr := errors.New("close error")
		mockDriver2.CloseFunc = func() error {
			return expectedErr
		}

		manager.AddDriver("mock1", mockDriver1)
		manager.AddDriver("mock2", mockDriver2)

		// Act
		err := manager.Close()

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("returns nil when no drivers", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		err := manager.Close()

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
