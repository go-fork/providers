# Scheduler v0.0.4 Release Notes

**Release Date**: May 28, 2025

## 🆕 Major Features

### Configuration-Driven System
- **Complete Config Support**: Added comprehensive `Config` struct and `RedisLockerOptions` for file-based configuration
- **Auto-Start Control**: New `auto_start` configuration option to control scheduler startup behavior
- **Dual Config System**: 
  - `RedisLockerOptions` with int values for config files (seconds/milliseconds)
  - `RedisLockerOptionsTime` with time.Duration for internal use
  - `ToTimeDuration()` conversion method between the two

### Distributed Locking Enhancement
- **Configuration-Driven Setup**: Automatic Redis locker setup when enabled in config
- **Provider Integration**: Seamless integration with redis provider for distributed locking
- **Auto-Configuration**: Provider automatically sets up distributed locking based on config

### Provider Architecture Improvement
- **Separation of Concerns**: Completely restructured ServiceProvider following best practices
- **Enhanced Manager**: Added `NewSchedulerWithConfig()` function for configuration-driven setup
- **Interface Compliance**: Fixed ServiceProvider to properly implement `di.ServiceProvider` interface

## 🔄 Breaking Changes

- **BREAKING**: `RedisLockerOptions` now uses int values instead of time.Duration for better config file compatibility
- **Provider Method**: Changed from `Provides()` to `Providers()` method for interface compliance

## 🛠️ Technical Improvements

### Code Quality
- **Static Analysis**: Fixed all staticcheck warnings by removing unnecessary type assertions
- **Interface Checks**: Used compile-time interface checks instead of runtime assertions
- **Test Compatibility**: Updated all tests to work with new int-based config values

### Dependencies
- **Enhanced Dependencies**: Added proper require dependencies for config and redis providers
- **Module Updates**: Updated go.mod with proper dependencies and replace directives

## 📝 Configuration Example

```yaml
# config/app.yaml
scheduler:
  auto_start: true
  distributed_lock:
    enabled: true
  options:
    key_prefix: "myapp_scheduler:"
    lock_duration: 60      # seconds
    max_retries: 5
    retry_delay: 200       # milliseconds

redis:
  default:
    addr: "localhost:6379"
    password: ""
    db: 0
```

## 🚀 Usage Example

```go
// Automatic configuration setup
app := di.New()
app.Register(config.NewServiceProvider())
app.Register(redis.NewServiceProvider())  // Required for distributed locking
app.Register(scheduler.NewServiceProvider())

app.Boot() // Scheduler auto-configured with distributed locking

// Use pre-configured scheduler
container := app.Container()
sched := container.Get("scheduler").(scheduler.Manager)

// All jobs automatically use distributed locking if enabled
sched.Every(5).Minutes().Do(func() {
    fmt.Println("This task uses distributed locking automatically")
})
```

## 🎯 Key Benefits

1. **Zero-Configuration Distributed Locking**: Just enable in config and it works
2. **Better Config File Support**: Int values are more natural in YAML/JSON
3. **Improved Developer Experience**: Less boilerplate code needed
4. **Enhanced Reliability**: Better error handling and interface compliance
5. **Future-Proof Architecture**: Clean separation between provider and manager logic

## 📋 Requirements

- Go 1.18 or later
- Redis (optional, only needed when distributed locking is enabled)

## 🔗 Dependencies

- `github.com/go-fork/providers/config` - For configuration management
- `github.com/go-fork/providers/redis` - For distributed locking (optional)
- `github.com/go-co-op/gocron` - Core scheduling functionality
