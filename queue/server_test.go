package queue

import (
	"context"
	"testing"
	"time"

	"github.com/go-fork/providers/queue/adapter"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewServer tests the creation of a server with Redis adapter
func TestNewServer(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	opts := ServerOptions{
		Concurrency:     5,
		PollingInterval: 1000,
		DefaultQueue:    "default",
		Queues:          []string{"default", "emails", "critical"},
	}

	srv := NewServer(redisClient, opts)

	assert.NotNil(t, srv, "Server should not be nil")
	queueServer, ok := srv.(*queueServer)
	assert.True(t, ok, "Server should be of type *queueServer")
	assert.Equal(t, opts.Concurrency, queueServer.options.Concurrency)
}

// TestNewServerWithAdapter tests the creation of a server with a custom adapter
func TestNewServerWithAdapter(t *testing.T) {
	memoryAdapter := adapter.NewMemoryQueue("test:")

	opts := ServerOptions{
		Concurrency:     5,
		PollingInterval: 1000,
		DefaultQueue:    "default",
	}

	server := NewServerWithAdapter(memoryAdapter, opts)

	assert.NotNil(t, server, "Server should not be nil")
	queueServer, ok := server.(*queueServer)
	assert.True(t, ok, "Server should be of type *queueServer")
	assert.Equal(t, opts.Concurrency, queueServer.options.Concurrency)
}

// TestServerRegisterHandler tests registering a single handler
func TestServerRegisterHandler(t *testing.T) {
	memoryAdapter := adapter.NewMemoryQueue("test:")

	opts := ServerOptions{
		Concurrency:  2,
		DefaultQueue: "default",
	}

	server := NewServerWithAdapter(memoryAdapter, opts).(*queueServer)

	// Register a handler
	handler := func(ctx context.Context, task *Task) error {
		return nil
	}

	server.RegisterHandler("test_task", handler)

	// Verify handler was registered
	value, ok := server.handlers.Load("test_task")
	assert.True(t, ok, "Handler should be stored in handlers map")
	assert.NotNil(t, value, "Handler function should not be nil")
}

// TestServerRegisterHandlers tests registering multiple handlers at once
func TestServerRegisterHandlers(t *testing.T) {
	memoryAdapter := adapter.NewMemoryQueue("test:")

	opts := ServerOptions{
		Concurrency:  2,
		DefaultQueue: "default",
	}

	server := NewServerWithAdapter(memoryAdapter, opts).(*queueServer)

	// Create some handlers
	handler1 := func(ctx context.Context, task *Task) error { return nil }
	handler2 := func(ctx context.Context, task *Task) error { return nil }

	// Register multiple handlers
	handlers := map[string]HandlerFunc{
		"task1": handler1,
		"task2": handler2,
	}

	server.RegisterHandlers(handlers)

	// Verify handlers were registered
	value1, ok1 := server.handlers.Load("task1")
	assert.True(t, ok1, "Handler for task1 should be stored")
	assert.NotNil(t, value1, "Handler function for task1 should not be nil")

	value2, ok2 := server.handlers.Load("task2")
	assert.True(t, ok2, "Handler for task2 should be stored")
	assert.NotNil(t, value2, "Handler function for task2 should not be nil")
}

// TestServerOptionsDefaults tests that default values are applied correctly
func TestServerOptionsDefaults(t *testing.T) {
	memoryAdapter := adapter.NewMemoryQueue("test:")

	// Create server with no options specified
	opts := ServerOptions{}

	server := NewServerWithAdapter(memoryAdapter, opts).(*queueServer)

	// Check that the default queue was applied
	assert.Contains(t, server.queues, "default", "Default queue should be 'default'")
}

// TestServerStart tests the Start method
func TestServerStart(t *testing.T) {
	memoryAdapter := adapter.NewMemoryQueue("test:")

	opts := ServerOptions{
		Concurrency:  2,
		DefaultQueue: "default",
	}

	server := NewServerWithAdapter(memoryAdapter, opts)

	// Start the server
	err := server.Start()
	assert.NoError(t, err, "Start should not return an error")

	// Try starting again - should fail
	err = server.Start()
	assert.Error(t, err, "Starting an already started server should return an error")
	assert.Contains(t, err.Error(), "already started", "Error should mention server is already started")

	// Clean up
	_ = server.Stop()
}

// TestServerStop tests the Stop method
func TestServerStop(t *testing.T) {
	memoryAdapter := adapter.NewMemoryQueue("test:")

	opts := ServerOptions{
		Concurrency:  2,
		DefaultQueue: "default",
	}

	server := NewServerWithAdapter(memoryAdapter, opts)

	// Start the server first
	err := server.Start()
	require.NoError(t, err, "Start should not return an error")

	// Stop the server
	err = server.Stop()
	assert.NoError(t, err, "Stop should not return an error")

	// Try stopping again - should fail
	err = server.Stop()
	assert.Error(t, err, "Stopping an already stopped server should return an error")
	assert.Contains(t, err.Error(), "not started", "Error should mention server is not started")
}

// TestServerWithRedisAdapter tests server with Redis adapter using redismock
func TestServerWithRedisAdapter(t *testing.T) {
	// Create redis client and mock
	redisClient, _ := redismock.NewClientMock()

	opts := ServerOptions{
		Concurrency:     5,
		PollingInterval: 1000,
		DefaultQueue:    "test-queue",
		ShutdownTimeout: 5 * time.Second,
	}

	server := NewServer(redisClient, opts)

	assert.NotNil(t, server, "Server should not be nil")

	// Register a handler
	server.RegisterHandler("test_task", func(ctx context.Context, task *Task) error {
		return nil
	})

	// We can test start/stop cycle
	err := server.Start()
	assert.NoError(t, err, "Start should not return an error")

	err = server.Stop()
	assert.NoError(t, err, "Stop should not return an error")
}

// TestServerWithCustomQueues tests server with custom queue names
func TestServerWithCustomQueues(t *testing.T) {
	memoryAdapter := adapter.NewMemoryQueue("test:")

	customQueues := []string{"high", "medium", "low"}
	opts := ServerOptions{
		Queues: customQueues,
	}

	server := NewServerWithAdapter(memoryAdapter, opts).(*queueServer)

	// Verify the queues were set correctly
	assert.Equal(t, len(customQueues), len(server.queues), "Server should have the correct number of queues")
	for i, queue := range customQueues {
		assert.Equal(t, queue, server.queues[i], "Queue name should match")
	}
}

// TestNewServerWithUniversalClient tests server creation with a non-standard Redis client
func TestNewServerWithUniversalClient(t *testing.T) {
	// Create a Redis Cluster client (which is not a standard *redis.Client)
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:6379"},
	})

	opts := ServerOptions{
		Concurrency:     5,
		PollingInterval: 1000,
		DefaultQueue:    "cluster-queue",
	}

	// This should trigger the fallback code path
	server := NewServer(clusterClient, opts)

	assert.NotNil(t, server, "Server should not be nil")

	// Start and stop to verify it's functional
	err := server.Start()
	assert.NoError(t, err, "Start should not return an error")

	err = server.Stop()
	assert.NoError(t, err, "Stop should not return an error")
}
