package redis

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Check if client config is initialized with default values
	if config.Client == nil {
		t.Error("Expected client config to be initialized")
	}

	// Check default client values
	if config.Client.Host != "localhost" {
		t.Errorf("Expected default host to be 'localhost', got '%s'", config.Client.Host)
	}
	if config.Client.Port != 6379 {
		t.Errorf("Expected default port to be 6379, got %d", config.Client.Port)
	}
	if config.Client.DB != 0 {
		t.Errorf("Expected default DB to be 0, got %d", config.Client.DB)
	}
	if config.Client.Timeout != 5 {
		t.Errorf("Expected default timeout to be 5, got %d", config.Client.Timeout)
	}

	// Check if universal config is initialized
	if config.Universal == nil {
		t.Error("Expected universal config to be initialized")
	}

	// Check default universal values
	if len(config.Universal.Addresses) != 1 || config.Universal.Addresses[0] != "localhost:6379" {
		t.Errorf("Expected default addresses to be ['localhost:6379'], got %v", config.Universal.Addresses)
	}
	if config.Universal.DB != 0 {
		t.Errorf("Expected default DB to be 0, got %d", config.Universal.DB)
	}
	if config.Universal.ClusterMode != false {
		t.Errorf("Expected default cluster mode to be false, got %v", config.Universal.ClusterMode)
	}
}

func TestGetClientOptions(t *testing.T) {
	config := &ClientConfig{
		Host:         "redis.example.com",
		Port:         6380,
		Password:     "secret",
		DB:           1,
		DialTimeout:  10,
		ReadTimeout:  5,
		WriteTimeout: 5,
		PoolSize:     20,
		MinIdleConns: 10,
	}

	options := config.GetClientOptions()

	if options.Addr != "redis.example.com:6380" {
		t.Errorf("Expected addr to be 'redis.example.com:6380', got '%s'", options.Addr)
	}
	if options.Password != "secret" {
		t.Errorf("Expected password to be 'secret', got '%s'", options.Password)
	}
	if options.DB != 1 {
		t.Errorf("Expected DB to be 1, got %d", options.DB)
	}
	if options.DialTimeout != 10*time.Second {
		t.Errorf("Expected dial timeout to be 10s, got %v", options.DialTimeout)
	}
	if options.PoolSize != 20 {
		t.Errorf("Expected pool size to be 20, got %d", options.PoolSize)
	}
}

func TestGetUniversalOptions(t *testing.T) {
	config := &UniversalConfig{
		Addresses:       []string{"redis1:6379", "redis2:6379"},
		Password:        "cluster-secret",
		DB:              2,
		DialTimeout:     15,
		ReadTimeout:     7,
		WriteTimeout:    7,
		MaxRetries:      5,
		MinRetryBackoff: 10,
		MaxRetryBackoff: 1000,
		PoolSize:        30,
		MinIdleConns:    15,
		ClusterMode:     true,
		MasterName:      "redis-master",
	}

	options := config.GetUniversalOptions()

	if len(options.Addrs) != 2 || options.Addrs[0] != "redis1:6379" || options.Addrs[1] != "redis2:6379" {
		t.Errorf("Expected addrs to be ['redis1:6379', 'redis2:6379'], got %v", options.Addrs)
	}
	if options.Password != "cluster-secret" {
		t.Errorf("Expected password to be 'cluster-secret', got '%s'", options.Password)
	}
	if options.DB != 2 {
		t.Errorf("Expected DB to be 2, got %d", options.DB)
	}
	if options.DialTimeout != 15*time.Second {
		t.Errorf("Expected dial timeout to be 15s, got %v", options.DialTimeout)
	}
	if options.MaxRetries != 5 {
		t.Errorf("Expected max retries to be 5, got %d", options.MaxRetries)
	}
	if options.MasterName != "redis-master" {
		t.Errorf("Expected master name to be 'redis-master', got '%s'", options.MasterName)
	}
}
