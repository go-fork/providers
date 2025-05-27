package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDefaultConfig tests the DefaultConfig function
func TestDefaultConfig(t *testing.T) {
	t.Run("returns valid default configuration", func(t *testing.T) {
		// Act
		config := DefaultConfig()

		// Assert
		assert.NotNil(t, config)
		assert.Equal(t, "memory", config.DefaultDriver)
		assert.Equal(t, 3600, config.DefaultTTL)
		assert.Equal(t, "cache:", config.Prefix)

		// Test drivers configuration
		assert.NotNil(t, config.Drivers)
		assert.NotNil(t, config.Drivers.Memory)
		assert.NotNil(t, config.Drivers.File)
		assert.NotNil(t, config.Drivers.Redis)
		assert.NotNil(t, config.Drivers.MongoDB)
	})

	t.Run("memory driver has correct default values", func(t *testing.T) {
		// Arrange & Act
		config := DefaultConfig()

		// Assert
		memory := config.Drivers.Memory
		assert.Equal(t, 3600, memory.DefaultTTL)
		assert.Equal(t, 600, memory.CleanupInterval)
		assert.Equal(t, 10000, memory.MaxItems)
	})

	t.Run("file driver has correct default values", func(t *testing.T) {
		// Arrange & Act
		config := DefaultConfig()

		// Assert
		file := config.Drivers.File
		assert.Equal(t, "./storage/cache", file.Path)
		assert.Equal(t, 3600, file.DefaultTTL)
		assert.Equal(t, ".cache", file.Extension)
		assert.Equal(t, 600, file.CleanupInterval)
	})

	t.Run("redis driver has correct default values", func(t *testing.T) {
		// Arrange & Act
		config := DefaultConfig()

		// Assert
		redis := config.Drivers.Redis
		assert.True(t, redis.Enabled)
		assert.Equal(t, 3600, redis.DefaultTTL)
		assert.Equal(t, "json", redis.Serializer)
	})

	t.Run("mongodb driver has correct default values", func(t *testing.T) {
		// Arrange & Act
		config := DefaultConfig()

		// Assert
		mongodb := config.Drivers.MongoDB
		assert.True(t, mongodb.Enabled)
		assert.Equal(t, "cache_db", mongodb.Database)
		assert.Equal(t, "cache_items", mongodb.Collection)
		assert.Equal(t, 3600, mongodb.DefaultTTL)
		assert.Equal(t, int64(0), mongodb.Hits)
		assert.Equal(t, int64(0), mongodb.Misses)
	})
}

// TestConfigGetDefaultExpiration tests the GetDefaultExpiration method
func TestConfigGetDefaultExpiration(t *testing.T) {
	testCases := []struct {
		name        string
		defaultTTL  int
		expectedDur time.Duration
	}{
		{
			name:        "zero TTL returns zero duration",
			defaultTTL:  0,
			expectedDur: 0,
		},
		{
			name:        "positive TTL returns correct duration",
			defaultTTL:  3600,
			expectedDur: time.Hour,
		},
		{
			name:        "negative TTL returns negative duration",
			defaultTTL:  -1,
			expectedDur: -time.Second,
		},
		{
			name:        "large TTL returns correct duration",
			defaultTTL:  86400,
			expectedDur: 24 * time.Hour,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			config := &Config{DefaultTTL: tc.defaultTTL}

			// Act
			duration := config.GetDefaultExpiration()

			// Assert
			assert.Equal(t, tc.expectedDur, duration)
		})
	}
}

// TestDriverMemoryConfigMethods tests DriverMemoryConfig methods
func TestDriverMemoryConfigMethods(t *testing.T) {
	t.Run("GetDefaultExpiration returns correct duration", func(t *testing.T) {
		// Arrange
		config := &DriverMemoryConfig{DefaultTTL: 1800}

		// Act
		duration := config.GetDefaultExpiration()

		// Assert
		assert.Equal(t, 30*time.Minute, duration)
	})

	t.Run("GetDefaultExpiration with zero TTL", func(t *testing.T) {
		// Arrange
		config := &DriverMemoryConfig{DefaultTTL: 0}

		// Act
		duration := config.GetDefaultExpiration()

		// Assert
		assert.Equal(t, time.Duration(0), duration)
	})

	t.Run("GetCleanupInterval returns correct duration", func(t *testing.T) {
		// Arrange
		config := &DriverMemoryConfig{CleanupInterval: 300}

		// Act
		duration := config.GetCleanupInterval()

		// Assert
		assert.Equal(t, 5*time.Minute, duration)
	})

	t.Run("GetCleanupInterval with zero interval", func(t *testing.T) {
		// Arrange
		config := &DriverMemoryConfig{CleanupInterval: 0}

		// Act
		duration := config.GetCleanupInterval()

		// Assert
		assert.Equal(t, time.Duration(0), duration)
	})
}

// TestDriverFileConfigMethods tests DriverFileConfig methods
func TestDriverFileConfigMethods(t *testing.T) {
	t.Run("GetDefaultExpiration returns correct duration", func(t *testing.T) {
		// Arrange
		config := &DriverFileConfig{DefaultTTL: 7200}

		// Act
		duration := config.GetDefaultExpiration()

		// Assert
		assert.Equal(t, 2*time.Hour, duration)
	})

	t.Run("GetFileCleanupInterval returns correct duration", func(t *testing.T) {
		// Arrange
		config := &DriverFileConfig{CleanupInterval: 1200}

		// Act
		duration := config.GetFileCleanupInterval()

		// Assert
		assert.Equal(t, 20*time.Minute, duration)
	})

	t.Run("GetFileCleanupInterval with negative interval", func(t *testing.T) {
		// Arrange
		config := &DriverFileConfig{CleanupInterval: -1}

		// Act
		duration := config.GetFileCleanupInterval()

		// Assert
		assert.Equal(t, -time.Second, duration)
	})
}

// TestDriverRedisConfigMethods tests DriverRedisConfig methods
func TestDriverRedisConfigMethods(t *testing.T) {
	t.Run("GetDefaultExpiration returns correct duration", func(t *testing.T) {
		// Arrange
		config := &DriverRedisConfig{DefaultTTL: 900}

		// Act
		duration := config.GetDefaultExpiration()

		// Assert
		assert.Equal(t, 15*time.Minute, duration)
	})

	t.Run("GetDefaultExpiration with large TTL", func(t *testing.T) {
		// Arrange
		config := &DriverRedisConfig{DefaultTTL: 604800} // 1 week

		// Act
		duration := config.GetDefaultExpiration()

		// Assert
		assert.Equal(t, 7*24*time.Hour, duration)
	})
}

// TestDriverMongodbConfigMethods tests DriverMongodbConfig methods
func TestDriverMongodbConfigMethods(t *testing.T) {
	t.Run("GetDefaultExpiration returns correct duration", func(t *testing.T) {
		// Arrange
		config := &DriverMongodbConfig{DefaultTTL: 2700}

		// Act
		duration := config.GetDefaultExpiration()

		// Assert
		assert.Equal(t, 45*time.Minute, duration)
	})

	t.Run("GetDefaultExpiration with minimum positive value", func(t *testing.T) {
		// Arrange
		config := &DriverMongodbConfig{DefaultTTL: 1}

		// Act
		duration := config.GetDefaultExpiration()

		// Assert
		assert.Equal(t, time.Second, duration)
	})
}

// TestConfigStructValidation tests config struct validation
func TestConfigStructValidation(t *testing.T) {
	t.Run("empty config struct", func(t *testing.T) {
		// Arrange
		config := &Config{}

		// Act & Assert
		assert.Empty(t, config.DefaultDriver)
		assert.Zero(t, config.DefaultTTL)
		assert.Empty(t, config.Prefix)
		assert.Zero(t, config.Drivers)
	})

	t.Run("partial config initialization", func(t *testing.T) {
		// Arrange
		config := &Config{
			DefaultDriver: "redis",
			DefaultTTL:    1800,
		}

		// Act & Assert
		assert.Equal(t, "redis", config.DefaultDriver)
		assert.Equal(t, 1800, config.DefaultTTL)
		assert.Empty(t, config.Prefix)
		assert.Zero(t, config.Drivers)
	})

	t.Run("full config initialization", func(t *testing.T) {
		// Arrange
		config := &Config{
			DefaultDriver: "file",
			DefaultTTL:    600,
			Prefix:        "myapp:",
			Drivers: DriversConfig{
				Memory: &DriverMemoryConfig{
					Enabled:         true,
					DefaultTTL:      600,
					CleanupInterval: 120,
					MaxItems:        5000,
				},
			},
		}

		// Act & Assert
		assert.Equal(t, "file", config.DefaultDriver)
		assert.Equal(t, 600, config.DefaultTTL)
		assert.Equal(t, "myapp:", config.Prefix)
		assert.NotNil(t, config.Drivers.Memory)
		assert.True(t, config.Drivers.Memory.Enabled)
		assert.Equal(t, 600, config.Drivers.Memory.DefaultTTL)
		assert.Equal(t, 120, config.Drivers.Memory.CleanupInterval)
		assert.Equal(t, 5000, config.Drivers.Memory.MaxItems)
	})
}

// TestDriverConfigStructValidation tests individual driver config structs
func TestDriverConfigStructValidation(t *testing.T) {
	t.Run("DriverMemoryConfig with all fields", func(t *testing.T) {
		// Arrange
		config := &DriverMemoryConfig{
			Enabled:         true,
			DefaultTTL:      3600,
			CleanupInterval: 600,
			MaxItems:        10000,
		}

		// Act & Assert
		assert.True(t, config.Enabled)
		assert.Equal(t, 3600, config.DefaultTTL)
		assert.Equal(t, 600, config.CleanupInterval)
		assert.Equal(t, 10000, config.MaxItems)
	})

	t.Run("DriverFileConfig with all fields", func(t *testing.T) {
		// Arrange
		config := &DriverFileConfig{
			Enabled:         true,
			Path:            "/tmp/cache",
			DefaultTTL:      1800,
			Extension:       ".tmp",
			CleanupInterval: 300,
		}

		// Act & Assert
		assert.True(t, config.Enabled)
		assert.Equal(t, "/tmp/cache", config.Path)
		assert.Equal(t, 1800, config.DefaultTTL)
		assert.Equal(t, ".tmp", config.Extension)
		assert.Equal(t, 300, config.CleanupInterval)
	})

	t.Run("DriverRedisConfig with all fields", func(t *testing.T) {
		// Arrange
		config := &DriverRedisConfig{
			Enabled:    false,
			DefaultTTL: 7200,
			Serializer: "msgpack",
		}

		// Act & Assert
		assert.False(t, config.Enabled)
		assert.Equal(t, 7200, config.DefaultTTL)
		assert.Equal(t, "msgpack", config.Serializer)
	})

	t.Run("DriverMongodbConfig with all fields", func(t *testing.T) {
		// Arrange
		config := &DriverMongodbConfig{
			Enabled:    true,
			Database:   "test_db",
			Collection: "test_collection",
			DefaultTTL: 4800,
			Hits:       100,
			Misses:     25,
		}

		// Act & Assert
		assert.True(t, config.Enabled)
		assert.Equal(t, "test_db", config.Database)
		assert.Equal(t, "test_collection", config.Collection)
		assert.Equal(t, 4800, config.DefaultTTL)
		assert.Equal(t, int64(100), config.Hits)
		assert.Equal(t, int64(25), config.Misses)
	})
}

// TestDurationConversionEdgeCases tests edge cases in duration conversion
func TestDurationConversionEdgeCases(t *testing.T) {
	t.Run("large positive value", func(t *testing.T) {
		// Arrange - use a large but safe value
		const largeTTL = 2147483647 // max int32, safe for duration conversion
		config := &Config{DefaultTTL: largeTTL}

		// Act
		duration := config.GetDefaultExpiration()

		// Assert
		expected := time.Duration(largeTTL) * time.Second
		assert.Equal(t, expected, duration)
	})

	t.Run("large negative value", func(t *testing.T) {
		// Arrange - use a large negative but safe value
		const largeTTL = -2147483647
		config := &Config{DefaultTTL: largeTTL}

		// Act
		duration := config.GetDefaultExpiration()

		// Assert
		expected := time.Duration(largeTTL) * time.Second
		assert.Equal(t, expected, duration)
	})

	t.Run("common values consistency", func(t *testing.T) {
		testCases := []struct {
			name     string
			seconds  int
			expected time.Duration
		}{
			{"one minute", 60, time.Minute},
			{"one hour", 3600, time.Hour},
			{"one day", 86400, 24 * time.Hour},
			{"one week", 604800, 7 * 24 * time.Hour},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Test Config
				config := &Config{DefaultTTL: tc.seconds}
				assert.Equal(t, tc.expected, config.GetDefaultExpiration())

				// Test DriverMemoryConfig
				memConfig := &DriverMemoryConfig{DefaultTTL: tc.seconds}
				assert.Equal(t, tc.expected, memConfig.GetDefaultExpiration())

				// Test DriverFileConfig
				fileConfig := &DriverFileConfig{DefaultTTL: tc.seconds}
				assert.Equal(t, tc.expected, fileConfig.GetDefaultExpiration())

				// Test DriverRedisConfig
				redisConfig := &DriverRedisConfig{DefaultTTL: tc.seconds}
				assert.Equal(t, tc.expected, redisConfig.GetDefaultExpiration())

				// Test DriverMongodbConfig
				mongoConfig := &DriverMongodbConfig{DefaultTTL: tc.seconds}
				assert.Equal(t, tc.expected, mongoConfig.GetDefaultExpiration())
			})
		}
	})
}

// TestConfigDeepCopy tests that configurations can be safely copied
func TestConfigDeepCopy(t *testing.T) {
	t.Run("modifying copy does not affect original", func(t *testing.T) {
		// Arrange
		original := DefaultConfig()
		copy := *original

		// Act - modify the copy
		copy.DefaultDriver = "redis"
		copy.DefaultTTL = 1800
		copy.Prefix = "modified:"

		// Assert - original should be unchanged
		assert.Equal(t, "memory", original.DefaultDriver)
		assert.Equal(t, 3600, original.DefaultTTL)
		assert.Equal(t, "cache:", original.Prefix)

		// Assert - copy should be modified
		assert.Equal(t, "redis", copy.DefaultDriver)
		assert.Equal(t, 1800, copy.DefaultTTL)
		assert.Equal(t, "modified:", copy.Prefix)
	})

	t.Run("modifying driver config in copy affects original", func(t *testing.T) {
		// Arrange
		original := DefaultConfig()
		copy := *original

		// Act - modify driver config in copy (this is a shallow copy issue)
		copy.Drivers.Memory.DefaultTTL = 999

		// Assert - both should be affected due to shared pointer
		assert.Equal(t, 999, original.Drivers.Memory.DefaultTTL)
		assert.Equal(t, 999, copy.Drivers.Memory.DefaultTTL)
	})
}
