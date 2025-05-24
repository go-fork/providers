package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Test Adapter config
	assert.Equal(t, "memory", config.Adapter.Default)
	assert.Equal(t, "queue:", config.Adapter.Memory.Prefix)

	// Test Redis config
	assert.Equal(t, "localhost:6379", config.Adapter.Redis.Address)
	assert.Equal(t, "", config.Adapter.Redis.Password)
	assert.Equal(t, 0, config.Adapter.Redis.DB)
	assert.False(t, config.Adapter.Redis.TLS)
	assert.Equal(t, "queue:", config.Adapter.Redis.Prefix)
	assert.False(t, config.Adapter.Redis.Cluster.Enabled)
	assert.Len(t, config.Adapter.Redis.Cluster.Addresses, 3)

	// Test Server config
	assert.Equal(t, 10, config.Server.Concurrency)
	assert.Equal(t, 1000, config.Server.PollingInterval)
	assert.Equal(t, "default", config.Server.DefaultQueue)
	assert.True(t, config.Server.StrictPriority)
	assert.Len(t, config.Server.Queues, 4)
	assert.Equal(t, 30, config.Server.ShutdownTimeout)
	assert.Equal(t, 1, config.Server.LogLevel)
	assert.Equal(t, 3, config.Server.RetryLimit)

	// Test Client config
	assert.Equal(t, "default", config.Client.DefaultOptions.Queue)
	assert.Equal(t, 3, config.Client.DefaultOptions.MaxRetry)
	assert.Equal(t, 30, config.Client.DefaultOptions.Timeout)
}
