package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Manager là interface chính để quản lý kết nối Redis.
//
// Manager cung cấp các phương thức để tạo và quản lý các kết nối Redis,
// bao gồm cả standard client và universal client.
type Manager interface {
	// Client trả về một Redis Client mới hoặc đã lưu trong cache.
	Client() (*redis.Client, error)

	// UniversalClient trả về một Redis Universal Client mới hoặc đã lưu trong cache.
	UniversalClient() (redis.UniversalClient, error)

	// GetConfig trả về cấu hình hiện tại của Manager.
	GetConfig() *Config

	// SetConfig đặt cấu hình cho Manager.
	SetConfig(config *Config)

	// Close đóng tất cả các kết nối Redis.
	Close() error

	// Ping kiểm tra kết nối tới Redis server.
	Ping(ctx context.Context) error

	// ClusterPing kiểm tra kết nối tới Redis Cluster.
	ClusterPing(ctx context.Context) error
}

// manager là implementation mặc định của Manager.
type manager struct {
	config          *Config
	client          *redis.Client
	universalClient redis.UniversalClient
}

// NewManager tạo một Manager mới với cấu hình mặc định.
func NewManager() Manager {
	return NewManagerWithConfig(DefaultConfig())
}

// NewManagerWithConfig tạo một Manager mới với cấu hình tùy chỉnh.
func NewManagerWithConfig(config *Config) Manager {
	return &manager{
		config: config,
	}
}

// Client trả về một Redis Client mới hoặc đã lưu trong cache.
func (m *manager) Client() (*redis.Client, error) {
	if m.client != nil {
		return m.client, nil
	}

	if m.config.Client == nil {
		return nil, fmt.Errorf("redis client configuration is missing")
	}

	options := m.config.Client.GetClientOptions()
	m.client = redis.NewClient(&redis.Options{
		Addr:         options.Addr,
		Password:     options.Password,
		DB:           options.DB,
		DialTimeout:  options.DialTimeout,
		ReadTimeout:  options.ReadTimeout,
		WriteTimeout: options.WriteTimeout,
		PoolSize:     options.PoolSize,
		MinIdleConns: options.MinIdleConns,
	})

	return m.client, nil
}

// UniversalClient trả về một Redis Universal Client mới hoặc đã lưu trong cache.
func (m *manager) UniversalClient() (redis.UniversalClient, error) {
	if m.universalClient != nil {
		return m.universalClient, nil
	}

	if m.config.Universal == nil {
		return nil, fmt.Errorf("redis universal client configuration is missing")
	}

	options := m.config.Universal.GetUniversalOptions()
	m.universalClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:           options.Addrs,
		MasterName:      options.MasterName,
		Password:        options.Password,
		DB:              options.DB,
		DialTimeout:     options.DialTimeout,
		ReadTimeout:     options.ReadTimeout,
		WriteTimeout:    options.WriteTimeout,
		MaxRetries:      options.MaxRetries,
		MinRetryBackoff: options.MinRetryBackoff,
		MaxRetryBackoff: options.MaxRetryBackoff,
		PoolSize:        options.PoolSize,
		MinIdleConns:    options.MinIdleConns,
		RouteRandomly:   options.RouteRandomly,
	})

	return m.universalClient, nil
}

// GetConfig trả về cấu hình hiện tại của Manager.
func (m *manager) GetConfig() *Config {
	return m.config
}

// SetConfig đặt cấu hình cho Manager.
func (m *manager) SetConfig(config *Config) {
	m.config = config

	// Reset các client hiện tại để chúng sẽ được tạo lại với cấu hình mới
	if m.client != nil {
		m.client.Close()
		m.client = nil
	}

	if m.universalClient != nil {
		m.universalClient.Close()
		m.universalClient = nil
	}
}

// Close đóng tất cả các kết nối Redis.
func (m *manager) Close() error {
	var clientErr, universalErr error

	if m.client != nil {
		clientErr = m.client.Close()
		m.client = nil
	}

	if m.universalClient != nil {
		universalErr = m.universalClient.Close()
		m.universalClient = nil
	}

	if clientErr != nil {
		return fmt.Errorf("failed to close Redis client: %w", clientErr)
	}

	if universalErr != nil {
		return fmt.Errorf("failed to close Redis universal client: %w", universalErr)
	}

	return nil
}

// Ping kiểm tra kết nối tới Redis server.
func (m *manager) Ping(ctx context.Context) error {
	client, err := m.Client()
	if err != nil {
		return err
	}

	_, err = client.Ping(ctx).Result()
	return err
}

// ClusterPing kiểm tra kết nối tới Redis Cluster.
func (m *manager) ClusterPing(ctx context.Context) error {
	client, err := m.UniversalClient()
	if err != nil {
		return err
	}

	_, err = client.Ping(ctx).Result()
	return err
}
