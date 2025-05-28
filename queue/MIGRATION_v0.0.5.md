# Migration Guide: Queue Provider v0.0.4 → v0.0.5

This guide helps you migrate from Queue Provider v0.0.4 to v0.0.5, which continues to enhance Redis Provider integration and provides improved documentation and stability.

## Key Changes

### 1. Version Updates

**v0.0.5** is primarily a maintenance and documentation release that builds upon the Redis Provider integration introduced in v0.0.4.

### 2. Documentation Improvements

- Updated README.md with clearer examples
- Enhanced configuration samples
- Improved production configuration templates
- Better error handling documentation

### 3. Configuration Structure (No Changes)

The configuration structure remains the same as v0.0.4:

```yaml
# Current structure (v0.0.4 & v0.0.5)
queue:
  adapter:
    default: "redis"
    redis:
      prefix: "queue:"
      provider_key: "default"  # References Redis Provider

redis:
  default:
    host: "localhost"
    port: 6379
    # ... other Redis settings
```

### 4. Code Changes (No Breaking Changes)

**Service Provider Registration (Same as v0.0.4):**
```go
import (
    "github.com/go-fork/di"
    "github.com/go-fork/providers/config"
    "github.com/go-fork/providers/redis"
    "github.com/go-fork/providers/scheduler"
    "github.com/go-fork/providers/queue"
)

func main() {
    app := di.New()
    
    // Register providers in order
    app.Register(config.NewServiceProvider())
    app.Register(redis.NewServiceProvider())  // Required for Redis adapter
    app.Register(scheduler.NewServiceProvider())
    app.Register(queue.NewServiceProvider())
    
    app.Boot()
}
```

### 5. Enhanced Redis Features (Available since v0.0.4)

All Redis features introduced in v0.0.4 remain available:

```go
// Get Redis queue adapter for advanced features
manager := container.MustMake("queue").(queue.Manager)
redisAdapter := manager.RedisAdapter()

if redisQueue, ok := redisAdapter.(adapter.QueueRedisAdapter); ok {
    ctx := context.Background()
    
    // Priority queues
    err := redisQueue.EnqueueWithPriority(ctx, "tasks", &task, 10)
    
    // TTL support
    err = redisQueue.EnqueueWithTTL(ctx, "temporary", &task, 1*time.Hour)
    
    // Batch operations
    err = redisQueue.EnqueueWithPipeline(ctx, "batch", tasks)
    
    // Monitoring
    info, err := redisQueue.GetQueueInfo(ctx, "tasks")
    
    // Health checks
    err = redisQueue.Ping(ctx)
}
```

## Migration Steps

### Step 1: Version Update

Update your `go.mod` file:

```bash
go get github.com/go-fork/providers/queue@v0.0.5
```

### Step 2: Review Configuration

If you're already on v0.0.4, no configuration changes are needed. Your existing configuration will continue to work.

### Step 3: Update Documentation References

Update any internal documentation that references v0.0.4 to v0.0.5.

### Step 4: Test Your Implementation

```bash
go test ./...
```

## What's New in v0.0.5

### 1. Improved Documentation

- Clearer examples in README.md
- Better production configuration samples
- Enhanced migration guides

### 2. Enhanced Configuration Templates

- Updated `configs/app.sample.yaml` with better comments
- Improved `configs/production.sample.yaml` with production best practices

### 3. Stability Improvements

- Better error handling in edge cases
- Improved test coverage
- Enhanced development experience

## Compatibility

### Backward Compatibility

- ✅ **Full backward compatibility** with v0.0.4
- ✅ No breaking changes in APIs
- ✅ Configuration structure unchanged
- ✅ All existing features continue to work

### Redis Provider Compatibility

- ✅ Compatible with Redis Provider v0.0.4+
- ✅ All Redis features from v0.0.4 continue to work
- ✅ No changes to Redis integration

## Common Issues and Solutions

### Issue 1: No Migration Needed

**Problem:** You're looking for migration steps but are already on v0.0.4.

**Solution:** If you're already using v0.0.4 with Redis Provider integration, simply update to v0.0.5 and you're done!

### Issue 2: Still on v0.0.3

**Problem:** You're migrating from v0.0.3 directly to v0.0.5.

**Solution:** Follow the v0.0.3 → v0.0.4 migration guide first (`MIGRATION_v0.0.4.md`), then update to v0.0.5.

## Testing Migration

```go
package main

import (
    "testing"
    "github.com/go-fork/di"
    "github.com/go-fork/providers/config"
    "github.com/go-fork/providers/redis"
    "github.com/go-fork/providers/queue"
)

func TestQueueProviderV005(t *testing.T) {
    app := di.New()
    
    // Register providers
    app.Register(config.NewServiceProvider())
    app.Register(redis.NewServiceProvider())
    app.Register(queue.NewServiceProvider())
    
    // This should work without issues in v0.0.5
    app.Boot()
    
    // Test queue manager
    manager := app.Container().MustMake("queue").(queue.Manager)
    if manager == nil {
        t.Fatal("Queue manager should be available")
    }
    
    // Test client
    client := app.Container().MustMake("queue.client").(queue.Client)
    if client == nil {
        t.Fatal("Queue client should be available")
    }
    
    t.Log("Queue Provider v0.0.5 migration successful!")
}
```

## Version Comparison

| Feature | v0.0.3 | v0.0.4 | v0.0.5 |
|---------|--------|--------|--------|
| Redis Provider Integration | ❌ | ✅ | ✅ |
| Priority Queues | ❌ | ✅ | ✅ |
| TTL Support | ❌ | ✅ | ✅ |
| Batch Operations | ❌ | ✅ | ✅ |
| Queue Monitoring | ❌ | ✅ | ✅ |
| Enhanced Documentation | ❌ | ⚠️ | ✅ |
| Production Templates | ❌ | ⚠️ | ✅ |
| Stability Improvements | ✅ | ✅ | ✅+ |

## Support

If you encounter any issues during migration:

1. Check the [README.md](README.md) for updated examples
2. Review the configuration samples in `configs/`
3. Refer to the [MIGRATION_v0.0.4.md](MIGRATION_v0.0.4.md) if migrating from v0.0.3
4. Check the [CHANGELOG.md](CHANGELOG.md) for detailed changes

## Conclusion

Migration to v0.0.5 is straightforward for users already on v0.0.4. This release focuses on documentation improvements, enhanced configuration templates, and stability improvements while maintaining full backward compatibility.

The Redis Provider integration remains the same, so all advanced Redis features introduced in v0.0.4 continue to work seamlessly in v0.0.5.
