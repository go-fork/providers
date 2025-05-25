# Changelog

## [Unreleased]

## [v0.0.4] - 2025-05-26

### Đã thêm
- **CI/CD**: Tối ưu hóa GitHub Actions workflows cho monorepo
- **CI/CD**: Thêm workflow module-release cho việc release từng module riêng lẻ
- **docs**: Cập nhật tài liệu hướng dẫn cài đặt module và quy trình release
- **build**: Cập nhật .goreleaser.yml để hỗ trợ monorepo structure

### Đã thay đổi
- **CI/CD**: Loại bỏ các workflow không cần thiết (benchmark, examples, metrics-report)
- **CI/CD**: Cập nhật CI workflow để test từng module riêng biệt
- **build**: Cải thiện cấu hình GoReleaser cho các module hiện có
- **docs**: Cập nhật README với cấu trúc module mới

### Đã sửa
- **modules**: Sửa lỗi versioning trong go.mod files
- **CI/CD**: Sửa lỗi golangci-lint installation trong workflows

## [v0.0.3] - 2025-05-25

### Đã thêm
- **mailer**: Cập nhật từ `github.com/go-gomail/gomail` sang `gopkg.in/gomail.v2` với API tương thích
- **scheduler**: Thêm mới Provider Scheduler để quản lý các tác vụ định kỳ
- **scheduler**: Nâng cao chức năng scheduler với validation và các phương thức fluent interface
- **queue**: Hoàn thiện chức năng worker với tích hợp scheduler
- **tests**: Thêm MockManager cho unit testing cấu hình
- **tests**: Nâng cao queue tests với các kịch bản và xử lý lỗi bổ sung

### Đã thay đổi
- **dependencies**: Nâng cấp di v0.0.2 lên v0.0.3
- **dependencies**: Nâng cấp fsnotify v1.8.0 lên v1.9.0
- **dependencies**: Nâng cấp các thư viện phụ thuộc, bao gồm:
  - golang.org/x/sys v0.29.0 -> v0.33.0
  - golang.org/x/text v0.21.0 -> v0.25.0
  - github.com/pelletier/go-toml/v2 v2.2.3 -> v2.2.4
  - github.com/sagikazarmark/locafero v0.7.0 -> v0.9.0
  - github.com/spf13/afero v1.12.0 -> v1.14.0
  - github.com/spf13/cast v1.7.1 -> v1.8.0
  - go.uber.org/multierr v1.9.0 -> v1.11.0
  - go.uber.org/atomic v1.9.0 -> v1.11.0

### Đã xóa
- **sms**: Xóa module go.mod cho SMS, chuẩn bị cho việc cải tiến module này

## v0.0.2 (2025-05-22)

### Cache Package

- **Core Features**
  - Thread-safe cache manager with support for multiple simultaneous drivers
  - TTL (Time To Live) support for cached items
  - "Remember" pattern for lazy computation and caching
  - Batch operations support for efficient bulk data handling
  - Direct serialization and deserialization for Go structs
  - Comprehensive error handling

- **Drivers**
  - **Memory Driver**: In-memory caching with automatic cleanup of expired items
  - **File Driver**: File-based persistent caching with atomic write operations
  - **Redis Driver**: Full Redis integration (v9+) with connection pooling
  - **MongoDB Driver**: MongoDB integration for distributed cache storage

- **Extensibility**
  - Interface-based driver design for easy custom implementations
  - Service Provider integration with DI containers

### Config Package

- **Core Features**
  - Multi-source configuration loading (YAML, JSON, ENV)
  - Dot notation access for hierarchical config values ("a.b.c")
  - Type-safe accessors (GetString, GetInt, GetBool, GetStringMap, GetStringSlice)
  - Default value support when keys don't exist
  - Direct struct mapping with automatic conversion
  - Thread-safe for concurrent access and updates

- **Formatters**
  - **YAML Formatter**: Support for YAML configuration files
  - **JSON Formatter**: Support for JSON configuration files
  - **ENV Formatter**: Environment variable integration with prefix filtering

- **Utilities**
  - Helper functions for nested map flattening and expansion
  - Service Provider integration with DI containers

### Log Package

- **Core Features**
  - Multiple severity levels (Debug, Info, Warning, Error, Fatal)
  - Concurrent handler support for multiple outputs
  - Thread-safe logging operations
  - Formatted log message support
  - Minimum log level filtering
  - Centralized log management

- **Handlers**
  - **Console Handler**: Terminal output with ANSI color support
  - **File Handler**: File-based logging with automatic rotation
  - **Stack Handler**: Automatic stack trace capture on errors

- **Extensibility**
  - Custom handler support via Handler interface
  - Service Provider integration with DI containers

