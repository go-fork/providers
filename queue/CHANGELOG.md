# Changelog

## [Unreleased]

## v0.0.3 - 2025-05-25

### Added

- Triển khai đơn giản, dễ bảo trì và mở rộng
- Hỗ trợ Redis làm message broker và bộ nhớ trong cho môi trường phát triển
- Hỗ trợ đặt lịch công việc (ngay lập tức, sau một khoảng thời gian, vào một thời điểm)
- Tự động thử lại tác vụ thất bại với chiến lược backoff
- Tích hợp dễ dàng với DI container của ứng dụng
- API đơn giản và tiện lợi cho người sử dụng

### Added
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
