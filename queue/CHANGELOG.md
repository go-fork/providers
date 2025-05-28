# Changelog

## [v0.0.5] - 2025-05-29

### Added
- **Enhanced Documentation**: Comprehensive Redis Provider integration guide with production examples
- **Migration Guide**: Complete v0.0.5 migration guide (`MIGRATION_v0.0.5.md`)
- **Release Notes**: Detailed release notes (`RELEASE_NOTES_v0.0.5.md`) with feature comparison
- **Production Templates**: Advanced configuration templates for production environments
- **Configuration Samples**: Enhanced app.sample.yaml and production.sample.yaml with better structure

### Enhanced
- **README.md**: Improved Redis Provider integration examples with clearer code samples
- **doc.go**: Comprehensive package documentation with all v0.0.5 features
- **Configuration Structure**: Better organized config templates with detailed comments
- **Best Practices**: Enhanced documentation for production deployment patterns

### Improved
- **Error Handling**: Better documentation for Redis connection error scenarios
- **Monitoring**: Enhanced examples for queue monitoring and debugging
- **Performance**: Documentation for optimal Redis Provider configuration
- **Testing**: Improved examples for testing Redis queue functionality

### Compatibility
- **Full Backward Compatibility**: No breaking changes from v0.0.4
- **Redis Provider Features**: All enhanced Redis features continue to work
- **Configuration**: Existing v0.0.4 configurations remain valid

## [v0.0.4] - 2025-05-29

### Added
- **Redis Provider Integration**: Complete integration with Redis Provider for centralized Redis configuration
- **Enhanced Redis Features**: 
  - Priority queues with Redis Sorted Sets (`EnqueueWithPriority`, `DequeueFromPriority`)
  - TTL support for temporary tasks (`EnqueueWithTTL`)
  - Batch operations with Redis pipelines (`EnqueueWithPipeline`, `MultiDequeue`)
  - Advanced queue monitoring (`GetQueueInfo`)
  - Health checks (`Ping`)
  - Development utilities (`FlushQueues`)
- **QueueRedisAdapter Interface**: New interface extending QueueAdapter with Redis-specific methods
- **Improved Architecture**: Clean separation between queue logic and Redis connection management
- **Comprehensive Documentation**: Updated README.md, doc.go with Redis Provider integration examples
- **Migration Guide**: Complete migration guide from v0.0.3 to v0.0.4
- **Production Configuration**: Sample production configurations with advanced Redis settings

### Changed
- **BREAKING**: Redis configuration moved from queue config to Redis Provider
- **BREAKING**: Queue Provider now requires Redis Provider in dependencies
- **Configuration Structure**: Simplified queue config, Redis details managed by Redis Provider
- **Constructor Patterns**: Added `NewRedisQueueWithProvider()` for Redis Provider integration
- **Test Coverage**: Updated all tests to work with new Redis Provider integration

### Removed
- **BREAKING**: Direct Redis connection fields from queue configuration
- **BREAKING**: `RedisClusterConfig` struct (moved to Redis Provider)
- **BREAKING**: Redis connection management from queue manager

### Migration
- See `MIGRATION_v0.0.4.md` for detailed migration instructions
- Update service provider registration to include Redis Provider
- Move Redis configuration from `queue.adapter.redis` to `redis` section
- Add `provider_key` reference in queue Redis configuration
- Dead letter queue support
- Message filtering capabilities

## [v0.0.3] - 2025-05-25

### Added
- Complete worker implementation with scheduler integration
- Support for delayed tasks
- New API for batch processing

### Changed
- Improved Redis adapter performance
- Optimized message handling in memory adapter

### Fixed
- Fixed error handling when Redis is unavailable
- Fixed task state management in distributed environments

## [v0.0.2] - 2025-05-22

### Added
- Queue management system with multiple adapter support:
  - Memory adapter for development environments
  - Redis adapter for production environments
- Asynchronous message processing
- Simple client API for enqueueing tasks:
  - Immediate task execution
  - Scheduled tasks (delayed by time interval)
  - Time-specific scheduled tasks
- Worker model with configurable retry logic and backoff strategies
- Server component for processing queue tasks
- Task payload serialization/deserialization
- Dependency Injection integration through ServiceProvider
- Task status tracking and management
