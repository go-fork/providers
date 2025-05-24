package queue

import (
	"time"

	"github.com/go-fork/providers/queue/adapter"
	"github.com/redis/go-redis/v9"
)

// Manager định nghĩa interface cho việc quản lý các thành phần queue.
type Manager interface {
	// RedisClient trả về Redis client.
	RedisClient() redis.UniversalClient

	// MemoryAdapter trả về memory queue adapter.
	MemoryAdapter() adapter.QueueAdapter

	// RedisAdapter trả về redis queue adapter.
	RedisAdapter() adapter.QueueAdapter

	// Adapter trả về queue adapter dựa trên cấu hình.
	Adapter(name string) adapter.QueueAdapter

	// Client trả về Client.
	Client() Client

	// Server trả về Server.
	Server() Server
}

// manager quản lý các thành phần trong queue.
type manager struct {
	config      Config
	client      Client
	server      Server
	redisClient redis.UniversalClient
	memoryQueue adapter.QueueAdapter
	redisQueue  adapter.QueueAdapter
}

// NewManager tạo một manager mới với cấu hình mặc định.
func NewManager() Manager {
	return NewManagerWithConfig(DefaultConfig())
}

// NewManagerWithConfig tạo một manager mới với cấu hình tùy chỉnh.
func NewManagerWithConfig(config Config) Manager {
	return &manager{
		config: config,
	}
}

// RedisClient trả về Redis client.
func (m *manager) RedisClient() redis.UniversalClient {
	if m.redisClient == nil {
		// Nếu client đã được cung cấp trong cấu hình, sử dụng nó
		if m.config.Adapter.Redis.Client != nil {
			m.redisClient = m.config.Adapter.Redis.Client
		} else {
			// Tạo client mới dựa trên cấu hình
			m.redisClient = redis.NewClient(&redis.Options{
				Addr:     m.config.Adapter.Redis.Address,
				Password: m.config.Adapter.Redis.Password,
				DB:       m.config.Adapter.Redis.DB,
			})
		}
	}
	return m.redisClient
}

// MemoryAdapter trả về memory queue adapter.
func (m *manager) MemoryAdapter() adapter.QueueAdapter {
	if m.memoryQueue == nil {
		m.memoryQueue = adapter.NewMemoryQueue(m.config.Adapter.Memory.Prefix)
	}
	return m.memoryQueue
}

// RedisAdapter trả về redis queue adapter.
func (m *manager) RedisAdapter() adapter.QueueAdapter {
	if m.redisQueue == nil {
		// Khởi tạo với client thông thường
		redisClient := m.RedisClient()
		redisStdClient, ok := redisClient.(*redis.Client)
		if !ok {
			// Fallback cho các trường hợp client không phải là *redis.Client
			redisStdClient = redis.NewClient(&redis.Options{
				Addr: m.config.Adapter.Redis.Address,
			})
		}
		m.redisQueue = adapter.NewRedisQueue(redisStdClient, m.config.Adapter.Redis.Prefix)
	}
	return m.redisQueue
}

// Adapter trả về queue adapter dựa trên cấu hình.
func (m *manager) Adapter(name string) adapter.QueueAdapter {
	if name == "" {
		name = m.config.Adapter.Default
	}

	switch name {
	case "redis":
		return m.RedisAdapter()
	case "memory":
		return m.MemoryAdapter()
	default:
		// Mặc định sử dụng memory adapter
		return m.MemoryAdapter()
	}
}

// Client trả về Client.
func (m *manager) Client() Client {
	if m.client == nil {
		if m.config.Adapter.Default == "redis" {
			m.client = NewClientWithUniversalClient(m.RedisClient())
		} else {
			m.client = NewClientWithAdapter(m.Adapter(m.config.Adapter.Default))
		}
	}
	return m.client
}

// Server trả về Server.
func (m *manager) Server() Server {
	if m.server == nil {
		serverOpts := ServerOptions{
			Concurrency:     m.config.Server.Concurrency,
			PollingInterval: m.config.Server.PollingInterval,
			DefaultQueue:    m.config.Server.DefaultQueue,
			StrictPriority:  m.config.Server.StrictPriority,
			Queues:          m.config.Server.Queues,
			ShutdownTimeout: time.Duration(m.config.Server.ShutdownTimeout) * time.Second,
			LogLevel:        m.config.Server.LogLevel,
			RetryLimit:      m.config.Server.RetryLimit,
		}

		if m.config.Adapter.Default == "redis" {
			m.server = NewServer(m.RedisClient(), serverOpts)
		} else {
			m.server = NewServerWithAdapter(m.Adapter(m.config.Adapter.Default), serverOpts)
		}
	}
	return m.server
}
