# Changelog

## [Unreleased]

## [v0.0.5] - 2025-05-29

### Fixed
- Fixed resource leak when replacing handlers in AddHandler method
- Ensured old handlers are properly closed when overwritten

## [v0.0.4] - 2025-05-29

### Changed
- Upgraded github.com/go-fork/di dependency from v0.0.4 to v0.0.5
- Implemented new interface methods from di.ServiceProvider: Requires() and Providers()
- Enhanced test coverage for ServiceProvider implementation
- Updated documentation in doc.go to reflect new interface methods

### Added
- Support for OpenTelemetry trace context
- Integration with structured logging standards
- Enhanced ColorFormatter for advanced color formatting
- Log rotation support with file size and time-based triggers
- Context logging with structured metadata

### Fixed
- Fixed concurrent logging handling in file handler
- Fixed memory leak in stack handler
- Improved file handler performance
- Memory optimization for concurrent logging

## [v0.0.3] - 2025-05-25

### Added

- **Đa dạng cấp độ log**: Hỗ trợ các cấp độ từ Debug, Info, Warning, Error đến Fatal.
- **Output đa dạng**: Hỗ trợ ghi log ra console (có màu) và file, dễ dàng mở rộng với custom handlers.
- **Thread-safe**: An toàn khi ghi log từ nhiều goroutines cùng lúc.
- **Hỗ trợ định dạng**: Hỗ trợ chuỗi định dạng kiểu Printf trong các thông điệp log.
- **Xử lý file linh hoạt**: Tự động xoay vòng file log khi đạt kích thước giới hạn.
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container.
- **Cấu trúc mở rộng**: Dễ dàng triển khai handler mới cho các output khác.
- **Truy xuất handler linh hoạt**: Lấy handler đã đăng ký để cấu hình hoặc tùy chỉnh thêm.

## [v0.0.3] - 2025-05-25

### Added
- Enhanced ColorFormatter for advanced color formatting
- Log rotation support with file size and time-based triggers
- Context logging with structured metadata

### Changed
- Improved file handler performance
- Memory optimization for concurrent logging

### Fixed
- Fixed concurrent logging handling in file handler
- Fixed memory leak in stack handler

## [v0.0.2] - 2025-05-22

### Added
- Thread-safe logging manager with multiple log levels (Debug, Info, Warning, Error, Fatal)
- Multiple output handlers:
  - Console handler with color support
  - File handler with configurable path
  - Stack handler for sending logs to multiple handlers
- Support for printf-style formatting with placeholders
- Configurable minimum log levels for each handler
- Automatic file rotation in file handler
- Dependency Injection integration through ServiceProvider
- Extensible handler API for custom log destinations
