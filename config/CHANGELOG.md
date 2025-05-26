# Changelog

## [Unreleased]

### Changed
- Upgraded github.com/go-fork/di dependency from v0.0.4 to v0.0.5
- Implemented new interface methods from di.ServiceProvider: Requires() and Providers()
- Enhanced test coverage for ServiceProvider implementation

## v0.0.4 - 2025-05-25

* See GitHub release notes

### Added
- Support for secure credential storage
- Configuration validation framework

## [v0.0.3] - 2025-05-25

### Added
- New utility APIs for configuration access
- MockManager for unit testing
- Support for dynamic environment variables

### Changed
- Updated spf13/viper dependency to latest version
- Improved performance for reading large configuration files

### Fixed
- Fixed issue when reading complex YAML configuration files
- Fixed environment variable handling with underscores

## [v0.0.2] - 2025-05-22

### Added
- Viper-based configuration management
- Support for multiple formats (YAML, JSON, TOML)
- Comprehensive API for type-safe configuration access:
  - String, Int, Bool, Float value retrieval
  - Duration, Time value retrieval
  - Slice, Map, and complex structure support
- Environment variable integration
- File-based configuration with search paths
- Automatic reload when configuration files change
- Default values support
- Configuration mounting and merging
- ServiceProvider for dependency injection integration
- Comprehensive error handling
