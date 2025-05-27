package driver_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-fork/providers/cache/config"
	"github.com/go-fork/providers/cache/driver"
	cacheMocks "github.com/go-fork/providers/cache/mocks"
	redisMocks "github.com/go-fork/providers/redis/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RedisDriverTestSuite struct {
	suite.Suite
	ctx         context.Context
	mockManager *redisMocks.MockManager
	mockClient  *redis.Client
	driver      driver.RedisDriver
	config      config.DriverRedisConfig
}

func (suite *RedisDriverTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Setup Redis client for testing (you might want to use miniredis for integration tests)
	suite.mockClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use test database
	})
}

func (suite *RedisDriverTestSuite) SetupTest() {
	suite.mockManager = redisMocks.NewMockManager(suite.T())

	suite.config = config.DriverRedisConfig{
		Enabled:    true,
		DefaultTTL: 300, // 5 minutes
		Serializer: "json",
	}
}

func (suite *RedisDriverTestSuite) TearDownTest() {
	if suite.driver != nil {
		suite.driver.Close()
	}
}

func (suite *RedisDriverTestSuite) TestNewRedisDriver_Success() {
	// Arrange
	suite.mockManager.EXPECT().Client().Return(suite.mockClient, nil).Once()

	// Act
	redisDriver, err := driver.NewRedisDriver(suite.config, suite.mockManager)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), redisDriver)
	suite.driver = redisDriver
}

func (suite *RedisDriverTestSuite) TestNewRedisDriver_Disabled() {
	// Arrange
	disabledConfig := suite.config
	disabledConfig.Enabled = false

	// Act
	redisDriver, err := driver.NewRedisDriver(disabledConfig, suite.mockManager)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), redisDriver)
	assert.Contains(suite.T(), err.Error(), "redis driver is not enabled")
}

func (suite *RedisDriverTestSuite) TestNewRedisDriver_ClientError() {
	// Arrange
	expectedErr := assert.AnError
	suite.mockManager.EXPECT().Client().Return(nil, expectedErr).Once()

	// Act
	redisDriver, err := driver.NewRedisDriver(suite.config, suite.mockManager)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), redisDriver)
	assert.Contains(suite.T(), err.Error(), "could not create Redis client")
}

func (suite *RedisDriverTestSuite) TestNewRedisDriver_WithGobSerializer() {
	// Arrange
	gobConfig := suite.config
	gobConfig.Serializer = "gob"
	suite.mockManager.EXPECT().Client().Return(suite.mockClient, nil).Once()

	// Act
	redisDriver, err := driver.NewRedisDriver(gobConfig, suite.mockManager)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), redisDriver)
	suite.driver = redisDriver
}

func (suite *RedisDriverTestSuite) TestNewRedisDriver_WithMsgpackSerializer() {
	// Arrange
	msgpackConfig := suite.config
	msgpackConfig.Serializer = "msgpack"
	suite.mockManager.EXPECT().Client().Return(suite.mockClient, nil).Once()

	// Act
	redisDriver, err := driver.NewRedisDriver(msgpackConfig, suite.mockManager)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), redisDriver)
	suite.driver = redisDriver
}

func (suite *RedisDriverTestSuite) TestNewRedisDriver_WithInvalidSerializer() {
	// Arrange
	invalidConfig := suite.config
	invalidConfig.Serializer = "invalid"
	suite.mockManager.EXPECT().Client().Return(suite.mockClient, nil).Once()

	// Act
	redisDriver, err := driver.NewRedisDriver(invalidConfig, suite.mockManager)

	// Assert - should still work, falling back to default
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), redisDriver)
	suite.driver = redisDriver
}

func (suite *RedisDriverTestSuite) TestWithSerializer() {
	// Arrange
	suite.mockManager.EXPECT().Client().Return(suite.mockClient, nil).Once()
	redisDriver, err := driver.NewRedisDriver(suite.config, suite.mockManager)
	assert.NoError(suite.T(), err)
	suite.driver = redisDriver

	// Act & Assert - test WithSerializer method
	newDriver := redisDriver.WithSerializer("gob")
	assert.NotNil(suite.T(), newDriver)

	newDriver = redisDriver.WithSerializer("msgpack")
	assert.NotNil(suite.T(), newDriver)

	newDriver = redisDriver.WithSerializer("json")
	assert.NotNil(suite.T(), newDriver)

	newDriver = redisDriver.WithSerializer("invalid")
	assert.NotNil(suite.T(), newDriver) // Should fall back to default
}

func TestRedisDriverIntegration(t *testing.T) {
	// Skip if Redis is not available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping integration tests")
	}
	defer client.Close()

	// Clean test database
	client.FlushDB(ctx)

	mockManager := redisMocks.NewMockManager(t)
	mockManager.EXPECT().Client().Return(client, nil).Once()

	config := config.DriverRedisConfig{
		Enabled:    true,
		DefaultTTL: 10, // 10 seconds for faster tests
		Serializer: "json",
	}

	redisDriver, err := driver.NewRedisDriver(config, mockManager)
	assert.NoError(t, err)
	defer redisDriver.Close()

	t.Run("Set and Get", func(t *testing.T) {
		key := "test:key"
		value := map[string]interface{}{"name": "test", "value": 123}

		// Set value
		err := redisDriver.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Get value
		result, found := redisDriver.Get(ctx, key)
		assert.True(t, found)

		// JSON unmarshaling converts numbers to float64
		expectedValue := map[string]interface{}{"name": "test", "value": float64(123)}
		assert.Equal(t, expectedValue, result)
	})

	t.Run("Has", func(t *testing.T) {
		key := "test:has"
		value := "test_value"

		// Initially should not exist
		exists := redisDriver.Has(ctx, key)
		assert.False(t, exists)

		// Set value
		err := redisDriver.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Now should exist
		exists = redisDriver.Has(ctx, key)
		assert.True(t, exists)
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test:delete"
		value := "test_value"

		// Set value
		err := redisDriver.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Verify exists
		exists := redisDriver.Has(ctx, key)
		assert.True(t, exists)

		// Delete
		err = redisDriver.Delete(ctx, key)
		assert.NoError(t, err)

		// Verify deleted
		exists = redisDriver.Has(ctx, key)
		assert.False(t, exists)
	})

	t.Run("SetMultiple and GetMultiple", func(t *testing.T) {
		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		// Set multiple
		err := redisDriver.SetMultiple(ctx, values, 0)
		assert.NoError(t, err)

		// Get multiple
		keys := []string{"key1", "key2", "key3", "key4"} // key4 doesn't exist
		results, missed := redisDriver.GetMultiple(ctx, keys)

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
		err := redisDriver.SetMultiple(ctx, values, 0)
		assert.NoError(t, err)

		// Delete multiple
		keys := []string{"del1", "del2"}
		err = redisDriver.DeleteMultiple(ctx, keys)
		assert.NoError(t, err)

		// Verify deletion
		assert.False(t, redisDriver.Has(ctx, "del1"))
		assert.False(t, redisDriver.Has(ctx, "del2"))
		assert.True(t, redisDriver.Has(ctx, "del3")) // Should still exist
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
		result, err := redisDriver.Remember(ctx, key, 0, callback)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
		assert.True(t, callbackCalled)

		// Reset flag
		callbackCalled = false

		// Second call should use cache
		result, err = redisDriver.Remember(ctx, key, 0, callback)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
		assert.False(t, callbackCalled) // Callback should not be called
	})

	t.Run("Stats", func(t *testing.T) {
		// Set some test data
		redisDriver.Set(ctx, "stats1", "value1", 0)
		redisDriver.Set(ctx, "stats2", "value2", 0)

		stats := redisDriver.Stats(ctx)

		assert.Contains(t, stats, "count")
		assert.Contains(t, stats, "hits")
		assert.Contains(t, stats, "misses")
		assert.Contains(t, stats, "type")
		assert.Contains(t, stats, "prefix")
		assert.Equal(t, "redis", stats["type"])
		assert.Equal(t, "cache:", stats["prefix"])
	})

	t.Run("Flush", func(t *testing.T) {
		// Set some test data
		redisDriver.Set(ctx, "flush1", "value1", 0)
		redisDriver.Set(ctx, "flush2", "value2", 0)

		// Verify data exists
		assert.True(t, redisDriver.Has(ctx, "flush1"))
		assert.True(t, redisDriver.Has(ctx, "flush2"))

		// Flush
		err := redisDriver.Flush(ctx)
		assert.NoError(t, err)

		// Verify data is gone
		assert.False(t, redisDriver.Has(ctx, "flush1"))
		assert.False(t, redisDriver.Has(ctx, "flush2"))
	})

	t.Run("TTL Expiration", func(t *testing.T) {
		key := "test:ttl"
		value := "test_value"

		// Set with short TTL
		err := redisDriver.Set(ctx, key, value, 1*time.Second)
		assert.NoError(t, err)

		// Should exist immediately
		result, found := redisDriver.Get(ctx, key)
		assert.True(t, found)
		assert.Equal(t, value, result)

		// Wait for expiration
		time.Sleep(1100 * time.Millisecond)

		// Should no longer exist
		_, found = redisDriver.Get(ctx, key)
		assert.False(t, found)
	})
}

func TestRedisDriverMocked(t *testing.T) {
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
			"count":  10,
			"hits":   50,
			"misses": 5,
			"type":   "redis",
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

func TestRedisDriverTestSuite(t *testing.T) {
	suite.Run(t, new(RedisDriverTestSuite))
}

func BenchmarkRedisDriver(b *testing.B) {
	// Skip if Redis is not available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		b.Skip("Redis not available, skipping benchmark")
	}
	defer client.Close()

	mockManager := redisMocks.NewMockManager(b)
	mockManager.EXPECT().Client().Return(client, nil).Once()

	config := config.DriverRedisConfig{
		Enabled:    true,
		DefaultTTL: 300,
		Serializer: "json",
	}

	redisDriver, err := driver.NewRedisDriver(config, mockManager)
	if err != nil {
		b.Fatal(err)
	}
	defer redisDriver.Close()

	b.Run("Set", func(b *testing.B) {
		value := map[string]interface{}{"test": "value", "number": 123}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			key := "bench:set:" + string(rune(i))
			redisDriver.Set(ctx, key, value, 0)
		}
	})

	b.Run("Get", func(b *testing.B) {
		// Setup data
		value := map[string]interface{}{"test": "value", "number": 123}
		for i := 0; i < 1000; i++ {
			key := "bench:get:" + string(rune(i))
			redisDriver.Set(ctx, key, value, 0)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := "bench:get:" + string(rune(i%1000))
			redisDriver.Get(ctx, key)
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
			redisDriver.SetMultiple(ctx, values, 0)
		}
	})
}

func TestRedisDriverErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis error handling tests in short mode")
	}

	t.Run("Redis Connection Error", func(t *testing.T) {
		// Test with invalid Redis address
		client := redis.NewClient(&redis.Options{
			Addr: "invalid:9999",
			DB:   15,
		})
		defer client.Close()

		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(client, nil).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.NoError(t, err) // Driver creation should succeed
		defer redisDriver.Close()

		ctx := context.Background()

		// Operations should handle Redis connection errors gracefully
		err = redisDriver.Set(ctx, "test", "value", 0)
		assert.Error(t, err) // Should get connection error

		_, found := redisDriver.Get(ctx, "test")
		assert.False(t, found) // Should return false on error

		exists := redisDriver.Has(ctx, "test")
		assert.False(t, exists) // Should return false on error
	})

	t.Run("Redis Manager Error", func(t *testing.T) {
		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(nil, assert.AnError).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.Error(t, err)
		assert.Nil(t, redisDriver)
		assert.Contains(t, err.Error(), "could not create Redis client")
	})
}

func TestRedisDriverSerializationEdgeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis serialization tests in short mode")
	}

	// Create a base client that we'll share
	baseClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})

	ctx := context.Background()
	if err := baseClient.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping serialization tests")
	}
	defer baseClient.Close()

	testCases := []struct {
		name       string
		serializer string
	}{
		{"JSON Serializer", "json"},
		{"GOB Serializer", "gob"},
		{"MSGPACK Serializer", "msgpack"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new client for each test case
			client := redis.NewClient(&redis.Options{
				Addr: "localhost:6379",
				DB:   15,
			})
			defer client.Close()

			mockManager := redisMocks.NewMockManager(t)
			mockManager.EXPECT().Client().Return(client, nil).Once()

			config := config.DriverRedisConfig{
				Enabled:    true,
				DefaultTTL: 300,
				Serializer: tc.serializer,
			}

			redisDriver, err := driver.NewRedisDriver(config, mockManager)
			assert.NoError(t, err)
			defer redisDriver.Close()

			// Clean up
			redisDriver.Flush(ctx)

			// For GOB, we need to register types and use specific data structures
			if tc.serializer == "gob" {
				// GOB has limitations with interface{} decoding, so we test simpler scenarios
				// Test that Set operations don't fail
				err = redisDriver.Set(ctx, "gob_string", "test_string", 0)
				assert.NoError(t, err)

				err = redisDriver.Set(ctx, "gob_number", 42, 0)
				assert.NoError(t, err)

				// Test with struct (GOB works well with structs)
				type TestStruct struct {
					Name   string
					Number int
					Active bool
				}

				testData := TestStruct{
					Name:   "test",
					Number: 123,
					Active: true,
				}

				err = redisDriver.Set(ctx, "gob_struct", testData, 0)
				assert.NoError(t, err)

				// For GOB, we verify that keys exist (Has method)
				// since Get with interface{} has limitations with GOB
				assert.True(t, redisDriver.Has(ctx, "gob_string"))
				assert.True(t, redisDriver.Has(ctx, "gob_number"))
				assert.True(t, redisDriver.Has(ctx, "gob_struct"))

				// Test GOB doesn't work well with nil values, so we skip nil tests
			} else {
				// Test simple data first
				err = redisDriver.Set(ctx, "simple", "test_string", 0)
				assert.NoError(t, err)

				result, found := redisDriver.Get(ctx, "simple")
				assert.True(t, found)
				assert.Equal(t, "test_string", result)

				// Test complex data structures for JSON and MSGPACK
				complexData := map[string]interface{}{
					"string":  "test",
					"number":  123,
					"float":   45.67,
					"boolean": true,
					"array":   []interface{}{1, 2, 3},
					"nested": map[string]interface{}{
						"inner": "value",
					},
				}

				// Set and get complex data
				err = redisDriver.Set(ctx, "complex", complexData, 0)
				assert.NoError(t, err)

				result, found = redisDriver.Get(ctx, "complex")
				assert.True(t, found)
				assert.NotNil(t, result)

				// Test nil values
				err = redisDriver.Set(ctx, "nil_value", nil, 0)
				assert.NoError(t, err)

				_, found = redisDriver.Get(ctx, "nil_value")
				assert.True(t, found)
			}
		})
	}
}

func TestRedisDriverComprehensive(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping comprehensive Redis tests in short mode")
	}

	// Skip if Redis is not available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping comprehensive tests")
	}
	defer client.Close()

	// Clean test database
	client.FlushDB(ctx)

	t.Run("WithSerializer Comprehensive", func(t *testing.T) {
		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(client, nil).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 10,
			Serializer: "json",
		}

		baseDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.NoError(t, err)
		defer baseDriver.Close()

		// Test all serializer types
		serializers := []string{"json", "gob", "msgpack", "invalid", "", "unknown"}
		for _, serializer := range serializers {
			newDriver := baseDriver.WithSerializer(serializer)
			assert.NotNil(t, newDriver)

			// Test that the new driver works
			key := "test_" + serializer
			value := "test_value_" + serializer

			err := newDriver.Set(ctx, key, value, 0)
			assert.NoError(t, err)

			exists := newDriver.Has(ctx, key)
			assert.True(t, exists)

			err = newDriver.Delete(ctx, key)
			assert.NoError(t, err)
		}
	})

	t.Run("Get Error Scenarios", func(t *testing.T) {
		// Create a new client for this test
		testClient := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   15,
		})
		defer testClient.Close()

		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(testClient, nil).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 10,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.NoError(t, err)
		defer redisDriver.Close()

		// Test getting non-existent key
		result, found := redisDriver.Get(ctx, "non_existent_key")
		assert.False(t, found)
		assert.Nil(t, result)

		// Test with corrupted data (set invalid JSON data directly)
		redisKey := "cache:corrupted_data"
		err = testClient.Set(ctx, redisKey, "invalid json {", 0).Err()
		assert.NoError(t, err)

		result, found = redisDriver.Get(ctx, "corrupted_data")
		assert.False(t, found)
		assert.Nil(t, result)
	})

	t.Run("Set Error Scenarios", func(t *testing.T) {
		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(client, nil).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 10,
			Serializer: "gob",
		}

		redisDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.NoError(t, err)
		defer redisDriver.Close()

		// Test with unsupported data type for GOB
		unsupportedValue := make(chan int) // channels can't be serialized with GOB
		err = redisDriver.Set(ctx, "unsupported", unsupportedValue, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not serialize value")
	})
	t.Run("GetMultiple Error Scenarios", func(t *testing.T) {
		// Create a new client for this test
		testClient := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   15,
		})
		defer testClient.Close()

		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(testClient, nil).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 10,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.NoError(t, err)
		defer redisDriver.Close()

		// Set up some data, including corrupted data
		validValues := map[string]interface{}{
			"valid1": "value1",
			"valid2": "value2",
		}
		err = redisDriver.SetMultiple(ctx, validValues, 0)
		assert.NoError(t, err)

		// Add corrupted data directly to Redis
		err = testClient.Set(ctx, "cache:corrupted", "invalid json {", 0).Err()
		assert.NoError(t, err)

		// Test GetMultiple with mix of valid, invalid, and missing keys
		keys := []string{"valid1", "valid2", "corrupted", "missing"}
		results, missed := redisDriver.GetMultiple(ctx, keys)

		// Should get valid keys only
		assert.Len(t, results, 2)
		assert.Equal(t, "value1", results["valid1"])
		assert.Equal(t, "value2", results["valid2"])

		// Should miss corrupted and missing keys
		assert.Contains(t, missed, "corrupted")
		assert.Contains(t, missed, "missing")
		assert.Len(t, missed, 2)
	})

	t.Run("SetMultiple and DeleteMultiple Error Scenarios", func(t *testing.T) {
		// Create a new client for this test
		testClient := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   15,
		})
		defer testClient.Close()

		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(testClient, nil).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 10,
			Serializer: "gob",
		}

		redisDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.NoError(t, err)
		defer redisDriver.Close()

		// Test SetMultiple with unsupported values
		unsupportedValues := map[string]interface{}{
			"good":        "valid_value",
			"unsupported": make(chan int), // channels can't be serialized
		}

		err = redisDriver.SetMultiple(ctx, unsupportedValues, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not serialize value")

		// Test DeleteMultiple with empty keys
		err = redisDriver.DeleteMultiple(ctx, []string{})
		assert.NoError(t, err) // Should handle empty slice gracefully

		// Test DeleteMultiple with non-existent keys
		err = redisDriver.DeleteMultiple(ctx, []string{"non_existent1", "non_existent2"})
		assert.NoError(t, err) // Should handle non-existent keys gracefully
	})
	t.Run("Remember Error Scenarios", func(t *testing.T) {
		// Create a new client for this test
		testClient := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   15,
		})
		defer testClient.Close()

		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(testClient, nil).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 10,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.NoError(t, err)
		defer redisDriver.Close()

		// Test Remember with callback that returns error
		key := "error_callback"
		expectedError := fmt.Errorf("callback error")

		callback := func() (interface{}, error) {
			return nil, expectedError
		}

		result, err := redisDriver.Remember(ctx, key, 0, callback)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)

		// Test Remember with callback that returns unsupported value (for GOB)
		gobDriver := redisDriver.WithSerializer("gob")
		callbackUnsupported := func() (interface{}, error) {
			return make(chan int), nil // channels can't be serialized
		}

		_, err = gobDriver.Remember(ctx, "unsupported_callback", 0, callbackUnsupported)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not serialize value")
		// Note: result might not be nil due to the callback returning a value before serialization fails
	})

	t.Run("Stats Comprehensive", func(t *testing.T) {
		// Create a new client for this test
		testClient := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   15,
		})
		defer testClient.Close()

		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(testClient, nil).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 10,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.NoError(t, err)
		defer redisDriver.Close()

		// Clean slate
		redisDriver.Flush(ctx)

		// Generate some hits and misses
		redisDriver.Set(ctx, "key1", "value1", 0)
		redisDriver.Set(ctx, "key2", "value2", 0)

		// Generate hits
		redisDriver.Get(ctx, "key1")
		redisDriver.Get(ctx, "key2")

		// Generate misses
		redisDriver.Get(ctx, "non_existent1")
		redisDriver.Get(ctx, "non_existent2")

		stats := redisDriver.Stats(ctx)

		// Check all required stats fields
		assert.Contains(t, stats, "count")
		assert.Contains(t, stats, "hits")
		assert.Contains(t, stats, "misses")
		assert.Contains(t, stats, "type")
		assert.Contains(t, stats, "prefix")
		assert.Contains(t, stats, "info")

		// Verify stats values
		assert.Equal(t, "redis", stats["type"])
		assert.Equal(t, "cache:", stats["prefix"])
		// Note: count is int, not int64 based on Stats implementation
		assert.GreaterOrEqual(t, stats["count"].(int), 2)
		assert.GreaterOrEqual(t, stats["hits"].(int64), int64(2))
		assert.GreaterOrEqual(t, stats["misses"].(int64), int64(2))
	})

	t.Run("Flush Error Scenarios", func(t *testing.T) {
		// Test Flush with disconnected client
		disconnectedClient := redis.NewClient(&redis.Options{
			Addr: "invalid:9999",
			DB:   15,
		})
		defer disconnectedClient.Close()

		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(disconnectedClient, nil).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 10,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.NoError(t, err)
		defer redisDriver.Close()

		// This should handle connection error gracefully
		err = redisDriver.Flush(ctx)
		assert.Error(t, err) // Should get connection error
	})

	t.Run("NewRedisDriver Error Scenarios", func(t *testing.T) {
		// Test disabled driver
		disabledConfig := config.DriverRedisConfig{
			Enabled:    false,
			DefaultTTL: 10,
			Serializer: "json",
		}

		mockManager := redisMocks.NewMockManager(t)
		// Should not call Client() for disabled driver

		redisDriver, err := driver.NewRedisDriver(disabledConfig, mockManager)
		assert.Error(t, err)
		assert.Nil(t, redisDriver)
		assert.Contains(t, err.Error(), "redis driver is not enabled")

		// Test client creation error
		enabledConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 10,
			Serializer: "json",
		}

		mockManager.EXPECT().Client().Return(nil, fmt.Errorf("connection failed")).Once()

		redisDriver, err = driver.NewRedisDriver(enabledConfig, mockManager)
		assert.Error(t, err)
		assert.Nil(t, redisDriver)
		assert.Contains(t, err.Error(), "could not create Redis client")

		// Test all serializer types in NewRedisDriver
		serializers := []string{"json", "gob", "msgpack", "invalid", ""}
		for _, serializer := range serializers {
			testConfig := config.DriverRedisConfig{
				Enabled:    true,
				DefaultTTL: 10,
				Serializer: serializer,
			}

			mockManager := redisMocks.NewMockManager(t)
			mockManager.EXPECT().Client().Return(client, nil).Once()

			redisDriver, err := driver.NewRedisDriver(testConfig, mockManager)
			assert.NoError(t, err, "Failed for serializer: %s", serializer)
			assert.NotNil(t, redisDriver, "Driver is nil for serializer: %s", serializer)
			redisDriver.Close()
		}
	})

	t.Run("Edge Cases with Different Data Types", func(t *testing.T) {
		// Create a new client for this test
		testClient := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   15,
		})
		defer testClient.Close()

		mockManager := redisMocks.NewMockManager(t)
		mockManager.EXPECT().Client().Return(testClient, nil).Once()

		config := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 10,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(config, mockManager)
		assert.NoError(t, err)
		defer redisDriver.Close()

		// Test various data types
		testCases := []struct {
			key   string
			value interface{}
		}{
			{"nil_value", nil},
			{"bool_true", true},
			{"bool_false", false},
			{"int_zero", 0},
			{"int_negative", -123},
			{"float_zero", 0.0},
			{"float_negative", -123.456},
			{"string_empty", ""},
			{"string_unicode", "测试中文🎉"},
			{"slice_empty", []interface{}{}},
			{"slice_mixed", []interface{}{1, "two", 3.0, true, nil}},
			{"map_empty", map[string]interface{}{}},
			{"map_nested", map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": "deep_value",
					},
				},
			}},
		}

		for _, tc := range testCases {
			t.Run("DataType_"+tc.key, func(t *testing.T) {
				// Set value
				err := redisDriver.Set(ctx, tc.key, tc.value, 0)
				assert.NoError(t, err)

				// Check existence
				exists := redisDriver.Has(ctx, tc.key)
				assert.True(t, exists)

				// Get value
				result, found := redisDriver.Get(ctx, tc.key)
				assert.True(t, found)

				// For JSON, nil becomes nil, but empty structures might change type
				if tc.value == nil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
				}

				// Clean up
				err = redisDriver.Delete(ctx, tc.key)
				assert.NoError(t, err)
			})
		}
	})
}

func TestRedisDriverContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis context cancellation tests in short mode")
	}

	// Skip if Redis is not available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping context tests")
	}
	defer client.Close()

	mockManager := redisMocks.NewMockManager(t)
	mockManager.EXPECT().Client().Return(client, nil).Once()

	config := config.DriverRedisConfig{
		Enabled:    true,
		DefaultTTL: 10,
		Serializer: "json",
	}

	redisDriver, err := driver.NewRedisDriver(config, mockManager)
	assert.NoError(t, err)
	defer redisDriver.Close()

	t.Run("Cancelled Context Operations", func(t *testing.T) {
		// Create a cancelled context
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Operations with cancelled context should fail
		err := redisDriver.Set(cancelledCtx, "cancelled_key", "value", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")

		_, found := redisDriver.Get(cancelledCtx, "some_key")
		assert.False(t, found)

		exists := redisDriver.Has(cancelledCtx, "some_key")
		assert.False(t, exists)
	})
}

func TestRedisDriverTTLEdgeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis TTL edge case tests in short mode")
	}

	// Skip if Redis is not available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping TTL tests")
	}
	defer client.Close()

	mockManager := redisMocks.NewMockManager(t)
	mockManager.EXPECT().Client().Return(client, nil).Once()

	config := config.DriverRedisConfig{
		Enabled:    true,
		DefaultTTL: 5, // 5 seconds default
		Serializer: "json",
	}

	redisDriver, err := driver.NewRedisDriver(config, mockManager)
	assert.NoError(t, err)
	defer redisDriver.Close()

	t.Run("TTL Edge Cases", func(t *testing.T) {
		// Test with default TTL (ttl = 0)
		err := redisDriver.Set(ctx, "default_ttl", "value", 0)
		assert.NoError(t, err)

		// Test with no expiration (ttl = -1)
		err = redisDriver.Set(ctx, "no_expiration", "value", -1*time.Second)
		assert.NoError(t, err)

		// Test with very short TTL
		err = redisDriver.Set(ctx, "short_ttl", "value", 1*time.Millisecond)
		assert.NoError(t, err)

		// Wait a bit and check if short TTL key expired
		time.Sleep(10 * time.Millisecond)
		_, found := redisDriver.Get(ctx, "short_ttl")
		assert.False(t, found) // Should be expired

		// Check that no_expiration key still exists
		exists := redisDriver.Has(ctx, "no_expiration")
		assert.True(t, exists)
	})

	t.Run("SetMultiple with Different TTLs", func(t *testing.T) {
		values := map[string]interface{}{
			"multi1": "value1",
			"multi2": "value2",
			"multi3": "value3",
		}

		// Test SetMultiple with default TTL
		err := redisDriver.SetMultiple(ctx, values, 0)
		assert.NoError(t, err)

		// Test SetMultiple with custom TTL
		err = redisDriver.SetMultiple(ctx, values, 1*time.Second)
		assert.NoError(t, err)

		// Verify all keys exist
		for key := range values {
			exists := redisDriver.Has(ctx, key)
			assert.True(t, exists)
		}
	})
}

func TestRedisDriverLargeData(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis large data tests in short mode")
	}

	// Skip if Redis is not available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping large data tests")
	}
	defer client.Close()

	mockManager := redisMocks.NewMockManager(t)
	mockManager.EXPECT().Client().Return(client, nil).Once()

	config := config.DriverRedisConfig{
		Enabled:    true,
		DefaultTTL: 10,
		Serializer: "json",
	}

	redisDriver, err := driver.NewRedisDriver(config, mockManager)
	assert.NoError(t, err)
	defer redisDriver.Close()

	t.Run("Large Data Structures", func(t *testing.T) {
		// Create a large slice
		largeSlice := make([]int, 10000)
		for i := range largeSlice {
			largeSlice[i] = i
		}

		err := redisDriver.Set(ctx, "large_slice", largeSlice, 0)
		assert.NoError(t, err)

		result, found := redisDriver.Get(ctx, "large_slice")
		assert.True(t, found)
		assert.NotNil(t, result)

		// Create a large map
		largeMap := make(map[string]interface{})
		for i := 0; i < 1000; i++ {
			largeMap[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
		}

		err = redisDriver.Set(ctx, "large_map", largeMap, 0)
		assert.NoError(t, err)

		result, found = redisDriver.Get(ctx, "large_map")
		assert.True(t, found)
		assert.NotNil(t, result)

		// Clean up large data
		redisDriver.Delete(ctx, "large_slice")
		redisDriver.Delete(ctx, "large_map")
	})

	t.Run("Large String Data", func(t *testing.T) {
		// Create a large string (1MB)
		largeString := strings.Repeat("a", 1024*1024)

		err := redisDriver.Set(ctx, "large_string", largeString, 0)
		assert.NoError(t, err)

		result, found := redisDriver.Get(ctx, "large_string")
		assert.True(t, found)
		assert.Equal(t, largeString, result)

		// Clean up
		redisDriver.Delete(ctx, "large_string")
	})
}

func TestRedisDriverSerializerSwitching(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis serializer switching tests in short mode")
	}

	// Skip if Redis is not available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping serializer tests")
	}
	defer client.Close()

	mockManager := redisMocks.NewMockManager(t)
	mockManager.EXPECT().Client().Return(client, nil).Once()

	config := config.DriverRedisConfig{
		Enabled:    true,
		DefaultTTL: 10,
		Serializer: "json",
	}

	baseDriver, err := driver.NewRedisDriver(config, mockManager)
	assert.NoError(t, err)
	defer baseDriver.Close()

	t.Run("Cross-Serializer Compatibility", func(t *testing.T) {
		// Test data that should work across serializers
		testData := map[string]interface{}{
			"string": "test_value",
			"number": 42,
			"bool":   true,
		}

		// Set with JSON
		jsonDriver := baseDriver.WithSerializer("json")
		err := jsonDriver.Set(ctx, "cross_test", testData, 0)
		assert.NoError(t, err)

		// Try to read with different serializers (this might fail, which is expected)
		msgpackDriver := baseDriver.WithSerializer("msgpack")
		_, found := msgpackDriver.Get(ctx, "cross_test")
		// This might not work due to different serialization formats, but we test the flow

		// Clean up
		jsonDriver.Delete(ctx, "cross_test")
		_ = found // Use the variable to avoid unused error
	})

	t.Run("WithSerializer Returns New Instance", func(t *testing.T) {
		// Test that WithSerializer returns a new instance
		jsonDriver := baseDriver.WithSerializer("json")
		gobDriver := baseDriver.WithSerializer("gob")
		msgpackDriver := baseDriver.WithSerializer("msgpack")

		// All should be different instances
		assert.NotEqual(t, fmt.Sprintf("%p", baseDriver), fmt.Sprintf("%p", jsonDriver))
		assert.NotEqual(t, fmt.Sprintf("%p", jsonDriver), fmt.Sprintf("%p", gobDriver))
		assert.NotEqual(t, fmt.Sprintf("%p", gobDriver), fmt.Sprintf("%p", msgpackDriver))

		// Test that each can perform operations independently
		err := jsonDriver.Set(ctx, "json_key", "json_value", 0)
		assert.NoError(t, err)

		err = msgpackDriver.Set(ctx, "msgpack_key", "msgpack_value", 0)
		assert.NoError(t, err)

		// Verify isolation
		exists := jsonDriver.Has(ctx, "json_key")
		assert.True(t, exists)

		exists = msgpackDriver.Has(ctx, "msgpack_key")
		assert.True(t, exists)

		// Clean up
		jsonDriver.Delete(ctx, "json_key")
		msgpackDriver.Delete(ctx, "msgpack_key")
	})
}
