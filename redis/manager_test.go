package redis

import (
	"context"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
)

// skipIntegrationTest kiểm tra nếu thử nghiệm tích hợp nên được bỏ qua
func skipIntegrationTest(t *testing.T) {
	if os.Getenv("REDIS_INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test. Set REDIS_INTEGRATION_TEST=1 to run.")
	}
}

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Error("Expected manager to be initialized, got nil")
	}

	config := manager.GetConfig()
	if config == nil {
		t.Error("Expected config to be initialized, got nil")
		return // Return early to avoid nil dereference
	}

	if config.Client == nil {
		t.Error("Expected client config to be initialized, got nil")
	}

	if config.Universal == nil {
		t.Error("Expected universal config to be initialized, got nil")
	}
}

func TestNewManagerWithConfig(t *testing.T) {
	customConfig := DefaultConfig()
	customConfig.Client.Host = "custom-redis-host"
	customConfig.Client.Port = 6380
	customConfig.Client.Password = "custom-password"

	manager := NewManagerWithConfig(customConfig)
	config := manager.GetConfig()

	if config.Client.Host != "custom-redis-host" {
		t.Errorf("Expected host to be 'custom-redis-host', got '%s'", config.Client.Host)
	}
	if config.Client.Port != 6380 {
		t.Errorf("Expected port to be 6380, got %d", config.Client.Port)
	}
	if config.Client.Password != "custom-password" {
		t.Errorf("Expected password to be 'custom-password', got '%s'", config.Client.Password)
	}
}

func TestManagerSetConfig(t *testing.T) {
	// Create a manager with initial config
	initialConfig := DefaultConfig()
	manager := NewManagerWithConfig(initialConfig).(*manager)

	// First create a client and universal client
	client, _ := manager.Client()
	if client == nil {
		t.Fatal("Failed to create client")
	}

	universalClient, _ := manager.UniversalClient()
	if universalClient == nil {
		t.Fatal("Failed to create universal client")
	}

	// Create a new config to set
	newConfig := DefaultConfig()
	newConfig.Client.Host = "new-host"
	newConfig.Client.Port = 7777

	// Set the new config
	manager.SetConfig(newConfig)

	// Check if config was updated
	if manager.config != newConfig {
		t.Error("Expected config to be updated to new config")
	}

	// Check if clients were reset
	if manager.client != nil {
		t.Error("Expected client to be reset to nil")
	}

	if manager.universalClient != nil {
		t.Error("Expected universal client to be reset to nil")
	}
}

func TestManagerClient(t *testing.T) {
	manager := NewManager()

	client, err := manager.Client()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client == nil {
		t.Error("Expected client to be initialized, got nil")
	}

	// Second call should return cached client
	cachedClient, err := manager.Client()
	if err != nil {
		t.Fatalf("Expected no error on cached client, got %v", err)
	}

	if client != cachedClient {
		t.Error("Expected cached client to be the same instance")
	}
}

func TestManagerUniversalClient(t *testing.T) {
	manager := NewManager()

	client, err := manager.UniversalClient()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client == nil {
		t.Error("Expected universal client to be initialized, got nil")
	}

	// Second call should return cached client
	cachedClient, err := manager.UniversalClient()
	if err != nil {
		t.Fatalf("Expected no error on cached universal client, got %v", err)
	}

	if client != cachedClient {
		t.Error("Expected cached universal client to be the same instance")
	}
}

func TestManagerClose(t *testing.T) {
	manager := NewManager()

	// Create clients to be closed
	_, err := manager.Client()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = manager.UniversalClient()
	if err != nil {
		t.Fatalf("Failed to create universal client: %v", err)
	}

	// Close all connections
	err = manager.Close()
	if err != nil {
		t.Fatalf("Failed to close connections: %v", err)
	}
}

func TestManagerPing_Integration(t *testing.T) {
	skipIntegrationTest(t)

	manager := NewManager()
	err := manager.Ping(context.Background())

	if err != nil {
		t.Fatalf("Failed to ping Redis server: %v", err)
	}
}

func TestManagerClusterPing_Integration(t *testing.T) {
	skipIntegrationTest(t)

	// For this test, you need a Redis cluster or a single Redis instance
	// that works with the universal client
	config := DefaultConfig()
	manager := NewManagerWithConfig(config)

	err := manager.ClusterPing(context.Background())
	if err != nil {
		t.Fatalf("Failed to ping Redis cluster: %v", err)
	}
}

// MockRedisClient tạo một client giả để test
func MockRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

// Additional Ping and ClusterPing tests are covered by the integration tests
// when REDIS_INTEGRATION_TEST=1 is set
