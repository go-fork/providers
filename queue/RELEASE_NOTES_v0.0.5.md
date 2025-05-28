# Queue Provider v0.0.5 Release Notes

**Release Date:** May 29, 2025  
**Version:** v0.0.5  
**Previous Version:** v0.0.4

## Overview

Queue Provider v0.0.5 is a maintenance and documentation-focused release that builds upon the Redis Provider integration introduced in v0.0.4. This release enhances developer experience with improved documentation, better configuration templates, and increased stability while maintaining full backward compatibility.

## 🚀 What's New

### 📚 Enhanced Documentation

- **Improved README.md**: Clearer examples and better organization of Redis Provider integration
- **Enhanced Configuration Samples**: More comprehensive and production-ready configuration templates
- **New Migration Guide**: Added `MIGRATION_v0.0.5.md` for users upgrading from v0.0.4
- **Better Comments**: Improved inline documentation throughout configuration files

### 🔧 Configuration Improvements

- **Updated Sample Configs**: Enhanced `app.sample.yaml` and `production.sample.yaml` with better examples
- **Version Tracking**: All configuration files now clearly indicate the version they're designed for
- **Production Best Practices**: Improved production configuration templates with recommended settings

### 🛠️ Developer Experience

- **Clearer Examples**: Better code examples in documentation showing Redis Provider integration
- **Enhanced Error Handling**: Improved error messages for common configuration issues
- **Better Testing**: Enhanced test coverage for edge cases

## 🔄 Compatibility

### ✅ Full Backward Compatibility

- **No Breaking Changes**: All APIs from v0.0.4 continue to work unchanged
- **Configuration Compatibility**: Existing v0.0.4 configurations work without modification
- **Feature Parity**: All Redis Provider features from v0.0.4 remain available

### 📋 Requirements

- **Go Version**: Go 1.19 or later
- **Dependencies**: Same as v0.0.4
  - Config Provider
  - Redis Provider
  - Scheduler Provider (for delayed tasks)

## 🚀 Migration

### From v0.0.4 to v0.0.5

**Migration is extremely simple:**

```bash
# Update your go.mod
go get github.com/go-fork/providers/queue@v0.0.5

# No code changes needed!
# No configuration changes needed!
```

### From v0.0.3 to v0.0.5

If you're still on v0.0.3, please follow the migration path:
1. First migrate to v0.0.4 using `MIGRATION_v0.0.4.md`
2. Then update to v0.0.5 (no additional changes needed)

## 📖 Key Features (Continued from v0.0.4)

### Redis Provider Integration
- Centralized Redis configuration through Redis Provider
- No Redis connection details in queue configuration
- Clean separation of concerns

### Advanced Redis Features
- **Priority Queues**: Task-level priority using Redis Sorted Sets
- **TTL Support**: Automatic expiration for temporary tasks
- **Batch Operations**: Pipeline support for high-throughput scenarios
- **Queue Monitoring**: Real-time queue statistics and health checks
- **Health Checks**: Built-in Redis connectivity monitoring

### Dual Adapter Support
- **Redis Adapter**: Production-ready with advanced features
- **Memory Adapter**: Perfect for development and testing

## 📝 Configuration Example

```yaml
# Modern v0.0.5 configuration
queue:
  adapter:
    default: "redis"
    redis:
      prefix: "queue:"
      provider_key: "default"  # References Redis Provider

# Redis managed by Redis Provider
redis:
  default:
    host: "localhost"
    port: 6379
    db: 0
    # ... other Redis settings
```

## 🔧 Service Provider Registration

```go
import (
    "github.com/go-fork/di"
    "github.com/go-fork/providers/config"
    "github.com/go-fork/providers/redis"    // Required for Redis adapter
    "github.com/go-fork/providers/scheduler"
    "github.com/go-fork/providers/queue"
)

func main() {
    app := di.New()
    
    app.Register(config.NewServiceProvider())
    app.Register(redis.NewServiceProvider())
    app.Register(scheduler.NewServiceProvider())
    app.Register(queue.NewServiceProvider())
    
    app.Boot()
}
```

## 🧪 Advanced Redis Usage

```go
// Access advanced Redis features
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

## 📊 Version Comparison

| Feature | v0.0.3 | v0.0.4 | v0.0.5 |
|---------|--------|--------|--------|
| Redis Provider Integration | ❌ | ✅ | ✅ |
| Priority Queues | ❌ | ✅ | ✅ |
| TTL Support | ❌ | ✅ | ✅ |
| Batch Operations | ❌ | ✅ | ✅ |
| Queue Monitoring | ❌ | ✅ | ✅ |
| Enhanced Documentation | ❌ | ⚠️ | ✅ |
| Production Templates | ❌ | ⚠️ | ✅ |
| Migration Guides | ❌ | ⚠️ | ✅ |

## 🐛 Bug Fixes

- Improved error handling in edge cases
- Better validation of configuration parameters
- Enhanced stability in high-load scenarios

## 📈 Performance

- No performance regressions from v0.0.4
- All Redis optimizations from v0.0.4 continue to apply
- Better connection pooling recommendations in documentation

## 🔮 Looking Forward

v0.0.5 sets the foundation for future enhancements while maintaining stability:

- **Documentation First**: Better docs lead to better developer experience
- **Configuration Clarity**: Clear, well-documented configuration templates
- **Stability Focus**: Rock-solid foundation for future features

## 📞 Support

- **Documentation**: Updated README.md with comprehensive examples
- **Migration Guides**: Step-by-step migration instructions
- **Configuration Samples**: Production-ready configuration templates
- **Issue Tracking**: GitHub issues for bug reports and feature requests

## 🙏 Acknowledgments

Thanks to all developers who provided feedback on v0.0.4 and helped shape this release with their real-world usage experiences.

---

**Happy Queueing! 🚀**

The Queue Provider Team
