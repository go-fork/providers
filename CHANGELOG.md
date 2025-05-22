# Changelog

## v0.0.2 (2025-05-22)
## v0.0.2 - 2025-05-22nn* See GitHub release notesnn
## v0.0.2 - 2025-05-22nn* See GitHub release notesnn

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

