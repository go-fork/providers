package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/providers/cache/mocks"
)

// TestNewManager tests the NewManager constructor
func TestNewManager(t *testing.T) {
	manager := NewManager()
	assert.NotNil(t, manager)

	// Test that newly created manager has no default driver
	_, found := manager.Get("any-key")
	assert.False(t, found)
}

// TestManagerGet tests the Get method with various scenarios
func TestManagerGet(t *testing.T) {
	t.Run("returns value when default driver is set and key exists", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Get(context.Background(), "test-key").Return("test-value", true)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, found := manager.Get("test-key")

		// Assert
		assert.True(t, found)
		assert.Equal(t, "test-value", value)
	})

	t.Run("returns not found when default driver is set but key doesn't exist", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Get(context.Background(), "nonexistent-key").Return(nil, false)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, found := manager.Get("nonexistent-key")

		// Assert
		assert.False(t, found)
		assert.Nil(t, value)
	})

	t.Run("returns not found when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		value, found := manager.Get("any-key")

		// Assert
		assert.False(t, found)
		assert.Nil(t, value)
	})

	t.Run("returns not found when default driver doesn't exist", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		manager.SetDefaultDriver("nonexistent")

		// Act
		value, found := manager.Get("any-key")

		// Assert
		assert.False(t, found)
		assert.Nil(t, value)
	})
}

// TestManagerSet tests the Set method with various scenarios
func TestManagerSet(t *testing.T) {
	t.Run("sets value successfully when default driver is configured", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Set(context.Background(), "test-key", "test-value", 5*time.Minute).Return(nil)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns error when default driver set operation fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("driver set error")
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Set(context.Background(), "test-key", "test-value", 5*time.Minute).Return(expectedError)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})

	t.Run("returns error when default driver doesn't exist", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		manager.SetDefaultDriver("nonexistent")

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManagerHas tests the Has method with various scenarios
func TestManagerHas(t *testing.T) {
	t.Run("returns true when key exists", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Has(context.Background(), "existing-key").Return(true)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		exists := manager.Has("existing-key")

		// Assert
		assert.True(t, exists)
	})

	t.Run("returns false when key doesn't exist", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Has(context.Background(), "nonexistent-key").Return(false)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		exists := manager.Has("nonexistent-key")

		// Assert
		assert.False(t, exists)
	})

	t.Run("returns false when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		exists := manager.Has("any-key")

		// Assert
		assert.False(t, exists)
	})

	t.Run("returns false when default driver doesn't exist", func(t *testing.T) {
		// Arrange
		manager := NewManager()
		manager.SetDefaultDriver("nonexistent")

		// Act
		exists := manager.Has("any-key")

		// Assert
		assert.False(t, exists)
	})
}

// TestManagerDelete tests the Delete method with various scenarios
func TestManagerDelete(t *testing.T) {
	t.Run("deletes key successfully when default driver is configured", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Delete(context.Background(), "test-key").Return(nil)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Delete("test-key")

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns error when driver delete operation fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("driver delete error")
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Delete(context.Background(), "test-key").Return(expectedError)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Delete("test-key")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		err := manager.Delete("test-key")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManagerFlush tests the Flush method with various scenarios
func TestManagerFlush(t *testing.T) {
	t.Run("flushes cache successfully when default driver is configured", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Flush(context.Background()).Return(nil)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Flush()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns error when driver flush operation fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("driver flush error")
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Flush(context.Background()).Return(expectedError)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Flush()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		err := manager.Flush()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManagerGetMultiple tests the GetMultiple method with various scenarios
func TestManagerGetMultiple(t *testing.T) {
	t.Run("gets multiple values successfully when default driver is configured", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2", "key3"}
		expectedFound := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		expectedMissing := []string{"key3"}

		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().GetMultiple(context.Background(), keys).Return(expectedFound, expectedMissing)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		found, missing := manager.GetMultiple(keys)

		// Assert
		assert.Equal(t, expectedFound, found)
		assert.Equal(t, expectedMissing, missing)
	})

	t.Run("returns empty map and all keys as missing when no default driver is set", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2", "key3"}
		manager := NewManager()

		// Act
		found, missing := manager.GetMultiple(keys)

		// Assert
		assert.Empty(t, found)
		assert.Equal(t, keys, missing)
	})

	t.Run("returns empty map and all keys as missing when default driver doesn't exist", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2", "key3"}
		manager := NewManager()
		manager.SetDefaultDriver("nonexistent")

		// Act
		found, missing := manager.GetMultiple(keys)

		// Assert
		assert.Empty(t, found)
		assert.Equal(t, keys, missing)
	})
}

// TestManagerSetMultiple tests the SetMultiple method with various scenarios
func TestManagerSetMultiple(t *testing.T) {
	t.Run("sets multiple values successfully when default driver is configured", func(t *testing.T) {
		// Arrange
		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		ttl := 10 * time.Minute

		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().SetMultiple(context.Background(), values, ttl).Return(nil)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.SetMultiple(values, ttl)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns error when driver setMultiple operation fails", func(t *testing.T) {
		// Arrange
		values := map[string]interface{}{
			"key1": "value1",
		}
		ttl := 10 * time.Minute
		expectedError := errors.New("driver setMultiple error")

		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().SetMultiple(context.Background(), values, ttl).Return(expectedError)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.SetMultiple(values, ttl)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		values := map[string]interface{}{
			"key1": "value1",
		}
		manager := NewManager()

		// Act
		err := manager.SetMultiple(values, 10*time.Minute)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManagerDeleteMultiple tests the DeleteMultiple method with various scenarios
func TestManagerDeleteMultiple(t *testing.T) {
	t.Run("deletes multiple keys successfully when default driver is configured", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2", "key3"}

		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().DeleteMultiple(context.Background(), keys).Return(nil)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.DeleteMultiple(keys)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns error when driver deleteMultiple operation fails", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2"}
		expectedError := errors.New("driver deleteMultiple error")

		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().DeleteMultiple(context.Background(), keys).Return(expectedError)

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.DeleteMultiple(keys)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2"}
		manager := NewManager()

		// Act
		err := manager.DeleteMultiple(keys)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManagerRemember tests the Remember method with various scenarios
func TestManagerRemember(t *testing.T) {
	t.Run("returns cached value when key exists", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Remember(context.Background(), "cached-key", 10*time.Minute, mock.AnythingOfType("func() (interface {}, error)")).Return("cached-value", nil)

		callback := func() (interface{}, error) {
			return "new-value", nil
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, err := manager.Remember("cached-key", 10*time.Minute, callback)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "cached-value", value)
	})

	t.Run("calls callback and caches result when key doesn't exist", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Remember(context.Background(), "new-key", 10*time.Minute, mock.AnythingOfType("func() (interface {}, error)")).Return("callback-value", nil)

		callback := func() (interface{}, error) {
			return "callback-value", nil
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, err := manager.Remember("new-key", 10*time.Minute, callback)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "callback-value", value)
	})

	t.Run("returns callback error when callback fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("callback error")
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Remember(context.Background(), "new-key", 10*time.Minute, mock.AnythingOfType("func() (interface {}, error)")).Return(nil, expectedError)

		callback := func() (interface{}, error) {
			return nil, expectedError
		}

		manager := NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, err := manager.Remember("new-key", 10*time.Minute, callback)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, value)
	})

	t.Run("returns error when no default driver is set", func(t *testing.T) {
		// Arrange
		callback := func() (interface{}, error) {
			return "value", nil
		}
		manager := NewManager()

		// Act
		value, err := manager.Remember("key", 10*time.Minute, callback)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, value)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManagerAddDriver tests the AddDriver method with various scenarios
func TestManagerAddDriver(t *testing.T) {
	t.Run("adds driver successfully", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		manager := NewManager()

		// Act
		manager.AddDriver("test-driver", mockDriver)

		// Assert
		driver, err := manager.Driver("test-driver")
		assert.NoError(t, err)
		assert.Equal(t, mockDriver, driver)
	})

	t.Run("sets first added driver as default", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Get(context.Background(), "test-key").Return("test-value", true)

		manager := NewManager()

		// Act
		manager.AddDriver("first-driver", mockDriver)

		// Assert - should be able to use operations without explicitly setting default
		value, found := manager.Get("test-key")
		assert.True(t, found)
		assert.Equal(t, "test-value", value)
	})

	t.Run("replaces existing driver with same name", func(t *testing.T) {
		// Arrange
		oldDriver := mocks.NewMockDriver(t)
		newDriver := mocks.NewMockDriver(t)
		manager := NewManager()

		// Act
		manager.AddDriver("same-name", oldDriver)
		manager.AddDriver("same-name", newDriver)

		// Assert
		driver, err := manager.Driver("same-name")
		assert.NoError(t, err)
		assert.Equal(t, newDriver, driver)
	})
}

// TestManagerSetDefaultDriver tests the SetDefaultDriver method
func TestManagerSetDefaultDriver(t *testing.T) {
	t.Run("sets default driver successfully when driver exists", func(t *testing.T) {
		// Arrange
		driver1 := mocks.NewMockDriver(t)
		driver2 := mocks.NewMockDriver(t)
		driver2.EXPECT().Get(context.Background(), "test-key").Return("value-from-driver2", true)

		manager := NewManager()
		manager.AddDriver("driver1", driver1)
		manager.AddDriver("driver2", driver2)

		// Act
		manager.SetDefaultDriver("driver2")

		// Assert - operations should use driver2
		value, found := manager.Get("test-key")
		assert.True(t, found)
		assert.Equal(t, "value-from-driver2", value)
	})

	t.Run("setting non-existent driver as default doesn't cause immediate error", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act - this should not panic or error immediately
		manager.SetDefaultDriver("non-existent")

		// Assert - but operations should fail
		_, found := manager.Get("test-key")
		assert.False(t, found)
	})
}

// TestManagerDriver tests the Driver method with various scenarios
func TestManagerDriver(t *testing.T) {
	t.Run("returns driver when it exists", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		manager := NewManager()
		manager.AddDriver("test-driver", mockDriver)

		// Act
		driver, err := manager.Driver("test-driver")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, mockDriver, driver)
	})

	t.Run("returns error when driver doesn't exist", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		driver, err := manager.Driver("non-existent")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, driver)
		assert.Contains(t, err.Error(), "driver 'non-existent' not found")
	})
}

// TestManagerStats tests the Stats method with various scenarios
func TestManagerStats(t *testing.T) {
	t.Run("returns stats for all drivers", func(t *testing.T) {
		// Arrange
		driver1Stats := map[string]interface{}{
			"hits":   100,
			"misses": 20,
		}
		driver2Stats := map[string]interface{}{
			"hits":   50,
			"misses": 10,
		}

		mockDriver1 := mocks.NewMockDriver(t)
		mockDriver1.EXPECT().Stats(context.Background()).Return(driver1Stats)

		mockDriver2 := mocks.NewMockDriver(t)
		mockDriver2.EXPECT().Stats(context.Background()).Return(driver2Stats)

		manager := NewManager()
		manager.AddDriver("driver1", mockDriver1)
		manager.AddDriver("driver2", mockDriver2)

		// Act
		stats := manager.Stats()

		// Assert
		assert.Len(t, stats, 2)
		assert.Equal(t, driver1Stats, stats["driver1"])
		assert.Equal(t, driver2Stats, stats["driver2"])
	})

	t.Run("returns empty map when no drivers are registered", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		stats := manager.Stats()

		// Assert
		assert.Empty(t, stats)
	})
}

// TestManagerClose tests the Close method with various scenarios
func TestManagerClose(t *testing.T) {
	t.Run("closes all drivers successfully", func(t *testing.T) {
		// Arrange
		mockDriver1 := mocks.NewMockDriver(t)
		mockDriver1.EXPECT().Close().Return(nil)

		mockDriver2 := mocks.NewMockDriver(t)
		mockDriver2.EXPECT().Close().Return(nil)

		manager := NewManager()
		manager.AddDriver("driver1", mockDriver1)
		manager.AddDriver("driver2", mockDriver2)

		// Act
		err := manager.Close()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns error when any driver close fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("close error")

		mockDriver1 := mocks.NewMockDriver(t)
		mockDriver1.EXPECT().Close().Return(nil)

		mockDriver2 := mocks.NewMockDriver(t)
		mockDriver2.EXPECT().Close().Return(expectedError)

		manager := NewManager()
		manager.AddDriver("driver1", mockDriver1)
		manager.AddDriver("driver2", mockDriver2)

		// Act
		err := manager.Close()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "close error")
	})

	t.Run("succeeds when no drivers are registered", func(t *testing.T) {
		// Arrange
		manager := NewManager()

		// Act
		err := manager.Close()

		// Assert
		assert.NoError(t, err)
	})
}

// TestManagerConcurrency tests concurrent access to manager
func TestManagerConcurrency(t *testing.T) {
	t.Run("concurrent operations don't cause race conditions", func(t *testing.T) {
		// Arrange
		mockDriver := mocks.NewMockDriver(t)
		mockDriver.EXPECT().Get(mock.Anything, mock.Anything).Return("value", true).Maybe()
		mockDriver.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
		mockDriver.EXPECT().Has(mock.Anything, mock.Anything).Return(true).Maybe()

		manager := NewManager()
		manager.AddDriver("concurrent-driver", mockDriver)
		manager.SetDefaultDriver("concurrent-driver")

		// Act - perform concurrent operations
		done := make(chan bool, 3)

		go func() {
			for i := 0; i < 100; i++ {
				manager.Get("key")
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 100; i++ {
				manager.Set("key", "value", time.Minute)
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 100; i++ {
				manager.Has("key")
			}
			done <- true
		}()

		// Wait for all goroutines to complete
		for i := 0; i < 3; i++ {
			<-done
		}

		// Assert - if we reach here without panic, the test passes
		assert.True(t, true)
	})
}
