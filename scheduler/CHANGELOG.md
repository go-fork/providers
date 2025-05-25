# Changelog

## [Unreleased]

## v0.0.3 - 2025-05-25

### Added

- Tích hợp toàn bộ tính năng của thư viện gocron vào DI container của ứng dụng
- Hỗ trợ nhiều loại lịch trình: theo khoảng thời gian, theo thời điểm cụ thể, biểu thức cron
- Hỗ trợ chế độ singleton để tránh chạy song song cùng một task
- Hỗ trợ distributed locking với Redis cho môi trường phân tán
- Hỗ trợ tag để nhóm và quản lý các task
- API fluent cho trải nghiệm lập trình dễ dàng
--
- Khóa Redis được tự động gia hạn sau khi đã chạy được 2/3 thời gian hết hạn
- Việc gia hạn xảy ra trong một goroutine riêng biệt
- Khi job hoàn thành, khóa sẽ được giải phóng

### Added
- Task cancellation API
- Health monitoring for scheduled tasks

## [v0.0.3] - 2025-05-25

### Added
- New Scheduler Provider for managing periodic tasks
- Support for cron expression syntax
- Redis integration for distributed scheduling
- Fluent interface for job configuration
- Lock mechanism for distributed environments

### Changed
- Performance optimization for handling multiple concurrent tasks
- Improved error handling in periodic tasks

### Fixed
- Timezone issues in cron expressions
- Memory leak when many jobs are scheduled and canceled

## [v0.0.2] - 2025-05-22

### Added
- Job scheduling system based on go-co-op/gocron
- Support for multiple scheduling methods:
  - Time interval-based scheduling
  - Specific time-based scheduling
  - Cron expression support
- Single-run and repeated task support
- Tag-based job grouping and management
- Distributed locking with Redis for cluster environments
- Singleton mode to prevent parallel execution of the same job
- Dependency Injection integration through ServiceProvider
- Fluent API for job configuration
