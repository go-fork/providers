# Changelog

## [Unreleased]

## v0.0.4 - 2025-05-25

### Added

- **Đa dạng driver**: Hỗ trợ bộ nhớ (Memory), tệp tin (File), Redis, MongoDB và khả năng mở rộng driver tùy chỉnh.
- **TTL (Time To Live)**: Tự động quản lý thời gian sống cho các mục trong cache.
- **Remember pattern**: Hỗ trợ tính toán lười biếng và lưu trữ kết quả trong cache.
- **Batch operations**: Thao tác hàng loạt để tối ưu hiệu suất.
- **Serialization**: Tự động chuyển đổi giữa cấu trúc dữ liệu Go và định dạng lưu trữ.
- **Thread-safe**: An toàn khi truy xuất và cập nhật đồng thời.
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container.
- **Extensible**: Dễ dàng mở rộng với driver tùy chỉnh thông qua interface Driver.

## v0.0.3 - 2025-05-25

### Added

- **Đa dạng driver**: Hỗ trợ bộ nhớ (Memory), tệp tin (File), Redis, MongoDB và khả năng mở rộng driver tùy chỉnh.
- **TTL (Time To Live)**: Tự động quản lý thời gian sống cho các mục trong cache.
- **Remember pattern**: Hỗ trợ tính toán lười biếng và lưu trữ kết quả trong cache.
- **Batch operations**: Thao tác hàng loạt để tối ưu hiệu suất.
- **Serialization**: Tự động chuyển đổi giữa cấu trúc dữ liệu Go và định dạng lưu trữ.
- **Thread-safe**: An toàn khi truy xuất và cập nhật đồng thời.
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container.
- **Extensible**: Dễ dàng mở rộng với driver tùy chỉnh thông qua interface Driver.

## v0.0.3 - 2025-05-25

### Added

- **Đa dạng driver**: Hỗ trợ bộ nhớ (Memory), tệp tin (File), Redis, MongoDB và khả năng mở rộng driver tùy chỉnh.
- **TTL (Time To Live)**: Tự động quản lý thời gian sống cho các mục trong cache.
- **Remember pattern**: Hỗ trợ tính toán lười biếng và lưu trữ kết quả trong cache.
- **Batch operations**: Thao tác hàng loạt để tối ưu hiệu suất.
- **Serialization**: Tự động chuyển đổi giữa cấu trúc dữ liệu Go và định dạng lưu trữ.
- **Thread-safe**: An toàn khi truy xuất và cập nhật đồng thời.
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container.
- **Extensible**: Dễ dàng mở rộng với driver tùy chỉnh thông qua interface Driver.

## v0.0.3 - 2025-05-25

### Added

- **Đa dạng driver**: Hỗ trợ bộ nhớ (Memory), tệp tin (File), Redis, MongoDB và khả năng mở rộng driver tùy chỉnh.
- **TTL (Time To Live)**: Tự động quản lý thời gian sống cho các mục trong cache.
- **Remember pattern**: Hỗ trợ tính toán lười biếng và lưu trữ kết quả trong cache.
- **Batch operations**: Thao tác hàng loạt để tối ưu hiệu suất.
- **Serialization**: Tự động chuyển đổi giữa cấu trúc dữ liệu Go và định dạng lưu trữ.
- **Thread-safe**: An toàn khi truy xuất và cập nhật đồng thời.
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container.
- **Extensible**: Dễ dàng mở rộng với driver tùy chỉnh thông qua interface Driver.

## v0.0.3 - 2025-05-25

### Added

- **Đa dạng driver**: Hỗ trợ bộ nhớ (Memory), tệp tin (File), Redis, MongoDB và khả năng mở rộng driver tùy chỉnh.
- **TTL (Time To Live)**: Tự động quản lý thời gian sống cho các mục trong cache.
- **Remember pattern**: Hỗ trợ tính toán lười biếng và lưu trữ kết quả trong cache.
- **Batch operations**: Thao tác hàng loạt để tối ưu hiệu suất.
- **Serialization**: Tự động chuyển đổi giữa cấu trúc dữ liệu Go và định dạng lưu trữ.
- **Thread-safe**: An toàn khi truy xuất và cập nhật đồng thời.
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container.
- **Extensible**: Dễ dàng mở rộng với driver tùy chỉnh thông qua interface Driver.

## [v0.0.3] - 2025-05-25

### Added
- **Đa dạng driver**: Hỗ trợ bộ nhớ (Memory), tệp tin (File), Redis, MongoDB và khả năng mở rộng driver tùy chỉnh.
- **TTL (Time To Live)**: Tự động quản lý thời gian sống cho các mục trong cache.
- **Remember pattern**: Hỗ trợ tính toán lười biếng và lưu trữ kết quả trong cache.
- **Batch operations**: Thao tác hàng loạt để tối ưu hiệu suất.
- **Serialization**: Tự động chuyển đổi giữa cấu trúc dữ liệu Go và định dạng lưu trữ.
- **Thread-safe**: An toàn khi truy xuất và cập nhật đồng thời.
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container.
- **Extensible**: Dễ dàng mở rộng với driver tùy chỉnh thông qua interface Driver.
- Support for graceful error handling during cache invalidation
- Integration with telemetry and monitoring systems

## [v0.0.2] - 2023-05-22

### Added
- Thread-safe cache manager with support for multiple concurrent drivers
- Multiple cache drivers supported:
  - Memory driver for in-RAM cache storage
  - File driver for filesystem-based cache
  - Redis driver for distributed caching
  - MongoDB driver for document-based cache storage
- TTL (Time To Live) support for cache entries with automatic expiration
- "Remember" pattern for lazy computation and caching of results
- Batch operations for efficient handling of large volumes of data
- Direct serialization and deserialization for Go structs
- Comprehensive error handling
- Dependency Injection integration through ServiceProvider
- Extensible API for custom cache drivers through Driver interface

### Changed
- Improved file driver performance
- Optimized Redis connection handling

### Fixed
- Fixed TTL handling in memory driver
- Fixed cache deletion in MongoDB driver
