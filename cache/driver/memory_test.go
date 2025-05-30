package driver_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.fork.vn/providers/cache/config"
	"go.fork.vn/providers/cache/driver"
	cacheMocks "go.fork.vn/providers/cache/mocks"
)

type MemoryDriverTestSuite struct {
	suite.Suite
	ctx    context.Context
	driver driver.MemoryDriver
	config config.DriverMemoryConfig
}

func (suite *MemoryDriverTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}

func (suite *MemoryDriverTestSuite) SetupTest() {
	suite.config = config.DriverMemoryConfig{
		DefaultTTL:      300, // 5 minutes
		CleanupInterval: 600, // 10 minutes
	}
}

func (suite *MemoryDriverTestSuite) TearDownTest() {
	if suite.driver != nil {
		suite.driver.Close()
	}
}

func (suite *MemoryDriverTestSuite) TestNewMemoryDriver_Success() {
	// Act
	memoryDriver := driver.NewMemoryDriver(suite.config)

	// Assert
	assert.NotNil(suite.T(), memoryDriver)
	suite.driver = memoryDriver
}

func (suite *MemoryDriverTestSuite) TestNewMemoryDriver_NoCleanup() {
	// Arrange
	noCleanupConfig := suite.config
	noCleanupConfig.CleanupInterval = 0

	// Act
	memoryDriver := driver.NewMemoryDriver(noCleanupConfig)

	// Assert
	assert.NotNil(suite.T(), memoryDriver)
	suite.driver = memoryDriver
}

func TestMemoryDriverIntegration(t *testing.T) {
	ctx := context.Background()

	memoryConfig := config.DriverMemoryConfig{
		DefaultTTL:      10, // 10 seconds for faster tests
		CleanupInterval: 5,  // 5 seconds cleanup
	}

	memoryDriver := driver.NewMemoryDriver(memoryConfig)
	defer memoryDriver.Close()

	t.Run("Set and Get", func(t *testing.T) {
		key := "test:key"
		value := map[string]interface{}{"name": "test", "value": 123}

		// Set value
		err := memoryDriver.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Get value
		result, found := memoryDriver.Get(ctx, key)
		assert.True(t, found)
		assert.Equal(t, value, result)
	})

	t.Run("Has", func(t *testing.T) {
		key := "test:has"
		value := "test_value"

		// Initially should not exist
		exists := memoryDriver.Has(ctx, key)
		assert.False(t, exists)

		// Set value
		err := memoryDriver.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Now should exist
		exists = memoryDriver.Has(ctx, key)
		assert.True(t, exists)
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test:delete"
		value := "test_value"

		// Set value
		err := memoryDriver.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Verify exists
		exists := memoryDriver.Has(ctx, key)
		assert.True(t, exists)

		// Delete
		err = memoryDriver.Delete(ctx, key)
		assert.NoError(t, err)

		// Verify deleted
		exists = memoryDriver.Has(ctx, key)
		assert.False(t, exists)
	})

	t.Run("SetMultiple and GetMultiple", func(t *testing.T) {
		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		// Set multiple
		err := memoryDriver.SetMultiple(ctx, values, 0)
		assert.NoError(t, err)

		// Get multiple
		keys := []string{"key1", "key2", "key3", "key4"} // key4 doesn't exist
		results, missed := memoryDriver.GetMultiple(ctx, keys)

		assert.Len(t, results, 3)
		assert.Len(t, missed, 1)
		assert.Contains(t, missed, "key4")
		assert.Equal(t, "value1", results["key1"])
		assert.Equal(t, "value2", results["key2"])
		assert.Equal(t, "value3", results["key3"])
	})

	t.Run("DeleteMultiple", func(t *testing.T) {
		values := map[string]interface{}{
			"del1": "value1",
			"del2": "value2",
			"del3": "value3",
		}

		// Set multiple
		err := memoryDriver.SetMultiple(ctx, values, 0)
		assert.NoError(t, err)

		// Delete multiple
		keys := []string{"del1", "del2"}
		err = memoryDriver.DeleteMultiple(ctx, keys)
		assert.NoError(t, err)

		// Verify deletion
		assert.False(t, memoryDriver.Has(ctx, "del1"))
		assert.False(t, memoryDriver.Has(ctx, "del2"))
		assert.True(t, memoryDriver.Has(ctx, "del3")) // Should still exist
	})

	t.Run("Remember", func(t *testing.T) {
		key := "test:remember"
		expectedValue := "computed_value"
		callbackCalled := false

		callback := func() (interface{}, error) {
			callbackCalled = true
			return expectedValue, nil
		}

		// First call should execute callback
		result, err := memoryDriver.Remember(ctx, key, 0, callback)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
		assert.True(t, callbackCalled)

		// Reset flag
		callbackCalled = false

		// Second call should use cache
		result, err = memoryDriver.Remember(ctx, key, 0, callback)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
		assert.False(t, callbackCalled) // Callback should not be called
	})

	t.Run("Stats", func(t *testing.T) {
		// Set some test data
		memoryDriver.Set(ctx, "stats1", "value1", 0)
		memoryDriver.Set(ctx, "stats2", "value2", 0)

		stats := memoryDriver.Stats(ctx)

		assert.Contains(t, stats, "count")
		assert.Contains(t, stats, "hits")
		assert.Contains(t, stats, "misses")
		assert.Contains(t, stats, "type")
		assert.Equal(t, "memory", stats["type"])
		assert.GreaterOrEqual(t, stats["count"], 2) // At least 2 items
	})

	t.Run("Flush", func(t *testing.T) {
		// Set some test data
		memoryDriver.Set(ctx, "flush1", "value1", 0)
		memoryDriver.Set(ctx, "flush2", "value2", 0)

		// Verify data exists
		assert.True(t, memoryDriver.Has(ctx, "flush1"))
		assert.True(t, memoryDriver.Has(ctx, "flush2"))

		// Flush
		err := memoryDriver.Flush(ctx)
		assert.NoError(t, err)

		// Verify data is gone
		assert.False(t, memoryDriver.Has(ctx, "flush1"))
		assert.False(t, memoryDriver.Has(ctx, "flush2"))
	})

	t.Run("TTL Expiration", func(t *testing.T) {
		key := "test:ttl"
		value := "test_value"

		// Set with short TTL
		err := memoryDriver.Set(ctx, key, value, 500*time.Millisecond)
		assert.NoError(t, err)

		// Should exist immediately
		result, found := memoryDriver.Get(ctx, key)
		assert.True(t, found)
		assert.Equal(t, value, result)

		// Wait for expiration
		time.Sleep(600 * time.Millisecond)

		// Should no longer exist
		_, found = memoryDriver.Get(ctx, key)
		assert.False(t, found)
	})

	t.Run("Item Expired Method", func(t *testing.T) {
		// Test Item.Expired() method directly
		now := time.Now()

		// Item with no expiration
		item1 := driver.Item{
			Value:      "test",
			Expiration: 0,
		}
		assert.False(t, item1.Expired())

		// Item that expired
		item2 := driver.Item{
			Value:      "test",
			Expiration: now.Add(-1 * time.Hour).UnixNano(),
		}
		assert.True(t, item2.Expired())

		// Item not yet expired
		item3 := driver.Item{
			Value:      "test",
			Expiration: now.Add(1 * time.Hour).UnixNano(),
		}
		assert.False(t, item3.Expired())
	})

	t.Run("Automatic Cleanup", func(t *testing.T) {
		// Create driver with very short cleanup interval
		quickCleanupConfig := config.DriverMemoryConfig{
			DefaultTTL:      1, // 1 second
			CleanupInterval: 1, // 1 second cleanup
		}
		quickDriver := driver.NewMemoryDriver(quickCleanupConfig)
		defer quickDriver.Close()

		// Set item with short TTL
		key := "cleanup:test"
		err := quickDriver.Set(ctx, key, "value", 500*time.Millisecond)
		assert.NoError(t, err)

		// Verify exists
		assert.True(t, quickDriver.Has(ctx, key))

		// Wait for automatic cleanup
		time.Sleep(2 * time.Second)

		// Item should be cleaned up automatically
		_, found := quickDriver.Get(ctx, key)
		assert.False(t, found)
	})
}

func TestMemoryDriverMocked(t *testing.T) {
	mockDriver := cacheMocks.NewMockDriver(t)
	ctx := context.Background()

	t.Run("Mock Driver Interface", func(t *testing.T) {
		key := "test_key"
		value := "test_value"

		// Setup expectations
		mockDriver.EXPECT().Set(ctx, key, value, time.Duration(0)).Return(nil).Once()
		mockDriver.EXPECT().Get(ctx, key).Return(value, true).Once()
		mockDriver.EXPECT().Has(ctx, key).Return(true).Once()
		mockDriver.EXPECT().Delete(ctx, key).Return(nil).Once()
		mockDriver.EXPECT().Close().Return(nil).Once()

		// Test operations
		err := mockDriver.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		result, found := mockDriver.Get(ctx, key)
		assert.True(t, found)
		assert.Equal(t, value, result)

		exists := mockDriver.Has(ctx, key)
		assert.True(t, exists)

		err = mockDriver.Delete(ctx, key)
		assert.NoError(t, err)

		err = mockDriver.Close()
		assert.NoError(t, err)
	})

	t.Run("Mock Multiple Operations", func(t *testing.T) {
		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		keys := []string{"key1", "key2", "key3"}
		expectedResults := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		expectedMissed := []string{"key3"}

		mockDriver.EXPECT().SetMultiple(ctx, values, time.Duration(0)).Return(nil).Once()
		mockDriver.EXPECT().GetMultiple(ctx, keys).Return(expectedResults, expectedMissed).Once()
		mockDriver.EXPECT().DeleteMultiple(ctx, keys).Return(nil).Once()

		err := mockDriver.SetMultiple(ctx, values, 0)
		assert.NoError(t, err)

		results, missed := mockDriver.GetMultiple(ctx, keys)
		assert.Equal(t, expectedResults, results)
		assert.Equal(t, expectedMissed, missed)

		err = mockDriver.DeleteMultiple(ctx, keys)
		assert.NoError(t, err)
	})

	t.Run("Mock Remember Operation", func(t *testing.T) {
		key := "remember_key"
		expectedValue := "computed_value"
		callback := func() (interface{}, error) {
			return expectedValue, nil
		}

		mockDriver.EXPECT().Remember(ctx, key, time.Duration(0), mock.MatchedBy(func(cb func() (interface{}, error)) bool {
			return cb != nil
		})).Return(expectedValue, nil).Once()

		result, err := mockDriver.Remember(ctx, key, 0, callback)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
	})

	t.Run("Mock Stats Operation", func(t *testing.T) {
		expectedStats := map[string]interface{}{
			"count":  5,
			"hits":   25,
			"misses": 3,
			"type":   "memory",
		}

		mockDriver.EXPECT().Stats(ctx).Return(expectedStats).Once()

		stats := mockDriver.Stats(ctx)
		assert.Equal(t, expectedStats, stats)
	})

	t.Run("Mock Flush Operation", func(t *testing.T) {
		mockDriver.EXPECT().Flush(ctx).Return(nil).Once()

		err := mockDriver.Flush(ctx)
		assert.NoError(t, err)
	})
}

func TestMemoryDriverTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryDriverTestSuite))
}

func TestMemoryDriverConcurrency(t *testing.T) {
	ctx := context.Background()
	memoryConfig := config.DriverMemoryConfig{
		DefaultTTL:      300,
		CleanupInterval: 60,
	}

	memoryDriver := driver.NewMemoryDriver(memoryConfig)
	defer memoryDriver.Close()

	t.Run("Concurrent Operations", func(t *testing.T) {
		// Test concurrent reads and writes
		done := make(chan bool, 100)

		// Start multiple goroutines for writing
		for i := 0; i < 50; i++ {
			go func(id int) {
				for j := 0; j < 10; j++ {
					key := fmt.Sprintf("concurrent:write:%d:%d", id, j)
					value := fmt.Sprintf("value_%d_%d", id, j)
					memoryDriver.Set(ctx, key, value, 0)
				}
				done <- true
			}(i)
		}

		// Start multiple goroutines for reading
		for i := 0; i < 50; i++ {
			go func(id int) {
				for j := 0; j < 10; j++ {
					key := fmt.Sprintf("concurrent:read:%d:%d", id, j)
					memoryDriver.Set(ctx, key, "read_value", 0)
					memoryDriver.Get(ctx, key)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 100; i++ {
			<-done
		}

		// Verify some data exists
		stats := memoryDriver.Stats(ctx)
		assert.Greater(t, stats["count"], 0)
	})
}

func BenchmarkMemoryDriver(b *testing.B) {
	ctx := context.Background()
	memoryConfig := config.DriverMemoryConfig{
		DefaultTTL:      300,
		CleanupInterval: 0, // Disable cleanup for benchmarks
	}

	memoryDriver := driver.NewMemoryDriver(memoryConfig)
	defer memoryDriver.Close()

	b.Run("Set", func(b *testing.B) {
		value := map[string]interface{}{"test": "value", "number": 123}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("bench:set:%d", i)
			memoryDriver.Set(ctx, key, value, 0)
		}
	})

	b.Run("Get", func(b *testing.B) {
		// Setup data
		value := map[string]interface{}{"test": "value", "number": 123}
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("bench:get:%d", i)
			memoryDriver.Set(ctx, key, value, 0)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("bench:get:%d", i%1000)
			memoryDriver.Get(ctx, key)
		}
	})

	b.Run("SetMultiple", func(b *testing.B) {
		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			memoryDriver.SetMultiple(ctx, values, 0)
		}
	})

	b.Run("GetMultiple", func(b *testing.B) {
		// Setup data
		values := map[string]interface{}{
			"bench1": "value1",
			"bench2": "value2",
			"bench3": "value3",
		}
		memoryDriver.SetMultiple(ctx, values, 0)

		keys := []string{"bench1", "bench2", "bench3"}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			memoryDriver.GetMultiple(ctx, keys)
		}
	})

	b.Run("Has", func(b *testing.B) {
		// Setup data
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("bench:has:%d", i)
			memoryDriver.Set(ctx, key, "value", 0)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("bench:has:%d", i%1000)
			memoryDriver.Has(ctx, key)
		}
	})
}
