# Migration Guide: Queue Provider v0.0.3 → v0.0.4

This guide helps you migrate from Queue Provider v0.0.3 to v0.0.4, which introduces Redis Provider integration and enhanced Redis features.

## Key Changes

### 1. Redis Configuration Migration

**Before (v0.0.3):**
```yaml
queue:
  adapter:
    redis:
      address: "localhost:6379"
      password: "secret"
      db: 0
      tls: true
      prefix: "queue:"
      cluster:
        enabled: true
        addresses:
          - "redis-1:6379"
          - "redis-2:6379"
```

**After (v0.0.4):**
```yaml
# Queue configuration simplified
queue:
  adapter:
    redis:
      prefix: "queue:"
      provider_key: "default"  # References redis.default

# Redis configuration moved to Redis Provider
redis:
  default:
    host: "localhost"
    port: 6379
    password: "secret"
    db: 0
    tls:
      enabled: true
    cluster:
      enabled: true
      hosts:
        - "redis-1:6379"
        - "redis-2:6379"
```

### 2. Service Provider Dependencies

**Before (v0.0.3):**
```go
app.Register(config.NewServiceProvider())
app.Register(scheduler.NewServiceProvider())
app.Register(queue.NewServiceProvider())
```

**After (v0.0.4):**
```go
app.Register(config.NewServiceProvider())
app.Register(redis.NewServiceProvider())  // New requirement
app.Register(scheduler.NewServiceProvider())
app.Register(queue.NewServiceProvider())
```

## Migration Steps

### Step 1: Update Dependencies

Add Redis Provider dependency to your `go.mod`:

```bash
go get github.com/go-fork/providers/redis@latest
```

### Step 2: Update Service Provider Registration

```go
import (
    "github.com/go-fork/providers/config"
    "github.com/go-fork/providers/redis"    // Add this import
    "github.com/go-fork/providers/scheduler"
    "github.com/go-fork/providers/queue"
)

func main() {
    app := di.New()
    
    app.Register(config.NewServiceProvider())
    app.Register(redis.NewServiceProvider())  // Add this line
    app.Register(scheduler.NewServiceProvider())
    app.Register(queue.NewServiceProvider())
    
    app.Boot()
}
```

### Step 3: Migrate Configuration

1. **Extract Redis configuration** from queue config
2. **Move to redis section** in your config file
3. **Add provider_key** reference in queue config

**Migration Script (config.yaml):**

```bash
# Backup your current config
cp config/app.yaml config/app.yaml.backup

# Update config structure manually or use this template:
```

```yaml
# New structure
redis:
  default:
    host: "{{ .OLD_QUEUE_ADAPTER_REDIS_ADDRESS_HOST }}"
    port: {{ .OLD_QUEUE_ADAPTER_REDIS_ADDRESS_PORT }}
    password: "{{ .OLD_QUEUE_ADAPTER_REDIS_PASSWORD }}"
    db: {{ .OLD_QUEUE_ADAPTER_REDIS_DB }}
    cluster:
      enabled: {{ .OLD_QUEUE_ADAPTER_REDIS_CLUSTER_ENABLED }}
      hosts: {{ .OLD_QUEUE_ADAPTER_REDIS_CLUSTER_ADDRESSES }}

queue:
  adapter:
    default: "redis"
    redis:
      prefix: "{{ .OLD_QUEUE_ADAPTER_REDIS_PREFIX }}"
      provider_key: "default"
  # ... rest of queue config remains the same
```

### Step 4: Test Migration

```go
// Test Redis provider integration
func TestMigration(t *testing.T) {
    app := setupApp()
    
    // Test Redis provider is available
    redisManager := app.Container().MustMake("redis")
    assert.NotNil(t, redisManager)
    
    // Test queue manager can access Redis
    queueManager := app.Container().MustMake("queue").(queue.Manager)
    redisClient := queueManager.RedisClient()
    assert.NotNil(t, redisClient)
    
    // Test enhanced Redis features
    redisAdapter := queueManager.RedisAdapter()
    if enhanced, ok := redisAdapter.(adapter.QueueRedisAdapter); ok {
        ctx := context.Background()
        
        // Test health check
        err := enhanced.Ping(ctx)
        assert.NoError(t, err)
        
        // Test queue info
        info, err := enhanced.GetQueueInfo(ctx, "test")
        assert.NoError(t, err)
        assert.NotNil(t, info)
    }
}
```

## New Features Available After Migration

### 1. Enhanced Redis Operations

```go
manager := app.Container().MustMake("queue").(queue.Manager)
if redisAdapter, ok := manager.RedisAdapter().(adapter.QueueRedisAdapter); ok {
    ctx := context.Background()
    
    // Priority queues
    err := redisAdapter.EnqueueWithPriority(ctx, "tasks", &task, 10)
    
    // TTL support
    err = redisAdapter.EnqueueWithTTL(ctx, "temporary", &task, 1*time.Hour)
    
    // Batch operations
    err = redisAdapter.EnqueueWithPipeline(ctx, "batch", tasks)
    
    // Multi-dequeue
    results, err := redisAdapter.MultiDequeue(ctx, "queue", 5)
    
    // Queue monitoring
    info, err := redisAdapter.GetQueueInfo(ctx, "queue")
    
    // Health checks
    err = redisAdapter.Ping(ctx)
}
```

### 2. Centralized Redis Management

```go
// Access Redis manager directly
redisManager := app.Container().MustMake("redis").(redis.Manager)
client := redisManager.Client("default")

// Multiple Redis instances
cacheClient := redisManager.Client("cache")
sessionClient := redisManager.Client("session")
```

### 3. Advanced Configuration

```yaml
redis:
  default:
    # Connection pooling
    pool_size: 100
    min_idle_conns: 20
    max_conn_age: 1800
    
    # Timeout configurations  
    dial_timeout: 10
    read_timeout: 5
    write_timeout: 5
    
    # Retry configurations
    max_retries: 5
    min_retry_backoff: 100
    max_retry_backoff: 2000
    
    # TLS support
    tls:
      enabled: true
      cert_file: "/path/to/cert.pem"
      key_file: "/path/to/key.pem"
      ca_file: "/path/to/ca.pem"
```

## Troubleshooting

### Common Issues

1. **"redis provider not found" error**
   - Ensure you've added `redis.NewServiceProvider()` before `queue.NewServiceProvider()`

2. **"provider_key 'default' not found" error**
   - Ensure your redis config has a `default` section
   - Check that `queue.adapter.redis.provider_key` matches your redis key

3. **Connection errors after migration**
   - Verify Redis configuration syntax (host/port vs address)
   - Check that all Redis connection settings are properly migrated

4. **Missing advanced features**
   - Type assert to `adapter.QueueRedisAdapter` to access enhanced features
   - Ensure you're using Redis adapter, not memory adapter

### Validation Checklist

- [ ] Redis Provider registered before Queue Provider
- [ ] Redis configuration moved to `redis` section
- [ ] Queue config references correct `provider_key`
- [ ] All Redis connection settings migrated
- [ ] Tests pass with new configuration
- [ ] Enhanced Redis features accessible
- [ ] Production environment updated
- [ ] Monitoring/alerting updated for new structure

## Rollback Plan

If you need to rollback to v0.0.3:

1. **Restore configuration backup:**
   ```bash
   cp config/app.yaml.backup config/app.yaml
   ```

2. **Downgrade packages:**
   ```bash
   go get github.com/go-fork/providers/queue@v0.0.3
   ```

3. **Remove Redis Provider:**
   ```go
   // Remove this line
   // app.Register(redis.NewServiceProvider())
   ```

4. **Test rollback:**
   ```bash
   go test ./...
   ```

## Support

If you encounter issues during migration:

1. Check the [Queue Provider Documentation](README.md)
2. Review [Redis Provider Documentation](../redis/README.md)
3. Check existing issues on GitHub
4. Create a new issue with migration details

## Benefits After Migration

- ✅ **Centralized Redis Management** - No duplicate Redis configurations
- ✅ **Enhanced Redis Features** - Priority queues, TTL, batch operations
- ✅ **Better Performance** - Optimized Redis operations and connection pooling
- ✅ **Improved Monitoring** - Queue info, health checks, and metrics
- ✅ **Production Ready** - Advanced Redis configurations and cluster support
- ✅ **Future Proof** - Foundation for upcoming features and improvements
