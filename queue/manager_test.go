package queue

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// TestNewManager tests creation of a new manager with default configuration
func TestNewManager(t *testing.T) {
	manager := NewManager()

	assert.NotNil(t, manager, "Manager should not be nil")
}

// TestNewManagerWithConfig tests creation of a new manager with custom configuration
func TestNewManagerWithConfig(t *testing.T) {
	config := Config{
		Adapter: AdapterConfig{
			Default: "memory",
			Memory: MemoryConfig{
				Prefix: "custom-prefix:",
			},
		},
	}

	manager := NewManagerWithConfig(config)

	assert.NotNil(t, manager, "Manager should not be nil")
}

// TestManagerRedisClient tests the RedisClient method
func TestManagerRedisClient(t *testing.T) {
	// Test with default config
	manager := NewManager()
	redisClient := manager.RedisClient()

	assert.NotNil(t, redisClient, "Redis client should not be nil")

	// Test with custom config that provides a client
	customRedisClient := redis.NewClient(&redis.Options{
		Addr: "custom-redis:6379",
	})

	config := Config{
		Adapter: AdapterConfig{
			Redis: RedisConfig{
				Client: customRedisClient,
			},
		},
	}

	managerWithCustomClient := NewManagerWithConfig(config)
	resultClient := managerWithCustomClient.RedisClient()

	assert.Same(t, customRedisClient, resultClient, "Should return the custom Redis client")
}

// TestManagerMemoryAdapter tests the MemoryAdapter method
func TestManagerMemoryAdapter(t *testing.T) {
	// Test with default config
	manager := NewManager()
	memoryAdapter := manager.MemoryAdapter()

	assert.NotNil(t, memoryAdapter, "Memory adapter should not be nil")
	// We already know it implements adapter.QueueAdapter as that's the return type

	// Test with custom config
	config := Config{
		Adapter: AdapterConfig{
			Memory: MemoryConfig{
				Prefix: "custom-prefix:",
			},
		},
	}

	customManager := NewManagerWithConfig(config)
	customMemoryAdapter := customManager.MemoryAdapter()

	assert.NotNil(t, customMemoryAdapter, "Memory adapter should not be nil")
}

// TestManagerRedisAdapter tests the RedisAdapter method
func TestManagerRedisAdapter(t *testing.T) {
	// Test with default config
	manager := NewManager()
	redisAdapter := manager.RedisAdapter()

	assert.NotNil(t, redisAdapter, "Redis adapter should not be nil")
	// We already know it implements adapter.QueueAdapter as that's the return type

	// Test with custom config
	config := Config{
		Adapter: AdapterConfig{
			Redis: RedisConfig{
				Prefix: "custom-prefix:",
			},
		},
	}

	customManager := NewManagerWithConfig(config)
	customRedisAdapter := customManager.RedisAdapter()

	assert.NotNil(t, customRedisAdapter, "Redis adapter should not be nil")
}

// TestManagerAdapter tests the Adapter method
func TestManagerAdapter(t *testing.T) {
	// Create a manager with memory as default
	memoryConfig := Config{
		Adapter: AdapterConfig{
			Default: "memory",
		},
	}
	memoryManager := NewManagerWithConfig(memoryConfig)

	// Get the default adapter
	defaultAdapter := memoryManager.Adapter("")
	assert.NotNil(t, defaultAdapter, "Default adapter should not be nil")

	// Get specific adapters
	memoryAdapter := memoryManager.Adapter("memory")
	assert.NotNil(t, memoryAdapter, "Memory adapter should not be nil")

	redisAdapter := memoryManager.Adapter("redis")
	assert.NotNil(t, redisAdapter, "Redis adapter should not be nil")

	// Test with unknown adapter type (should default to memory)
	unknownAdapter := memoryManager.Adapter("unknown")
	assert.NotNil(t, unknownAdapter, "Unknown adapter should not be nil")

	// Create a manager with redis as default
	redisConfig := Config{
		Adapter: AdapterConfig{
			Default: "redis",
		},
	}
	redisManager := NewManagerWithConfig(redisConfig)

	// Get the default adapter
	redisDefaultAdapter := redisManager.Adapter("")
	assert.NotNil(t, redisDefaultAdapter, "Default adapter should not be nil")
}

// TestManagerClient tests the Client method
func TestManagerClient(t *testing.T) {
	// Test with memory adapter as default
	memoryConfig := Config{
		Adapter: AdapterConfig{
			Default: "memory",
		},
	}
	memoryManager := NewManagerWithConfig(memoryConfig)

	memoryClient := memoryManager.Client()
	assert.NotNil(t, memoryClient, "Client should not be nil")

	// Test with redis adapter as default
	redisConfig := Config{
		Adapter: AdapterConfig{
			Default: "redis",
		},
	}
	redisManager := NewManagerWithConfig(redisConfig)

	redisClient := redisManager.Client()
	assert.NotNil(t, redisClient, "Client should not be nil")
}

// TestManagerServer tests the Server method
func TestManagerServer(t *testing.T) {
	// Test with memory adapter as default
	memoryConfig := Config{
		Adapter: AdapterConfig{
			Default: "memory",
		},
		Server: ServerConfig{
			Concurrency:     3,
			DefaultQueue:    "test-queue",
			ShutdownTimeout: 10,
		},
	}
	memoryManager := NewManagerWithConfig(memoryConfig)

	memoryServer := memoryManager.Server()
	assert.NotNil(t, memoryServer, "Server should not be nil")
	assert.IsType(t, &ServerImpl{}, memoryServer, "Should return a queue server")

	// Test with redis adapter as default
	redisConfig := Config{
		Adapter: AdapterConfig{
			Default: "redis",
		},
		Server: ServerConfig{
			Concurrency:     5,
			DefaultQueue:    "priority-queue",
			ShutdownTimeout: 15,
		},
	}
	redisManager := NewManagerWithConfig(redisConfig)

	redisServer := redisManager.Server()
	assert.NotNil(t, redisServer, "Server should not be nil")
	assert.IsType(t, &ServerImpl{}, redisServer, "Should return a queue server")
}
