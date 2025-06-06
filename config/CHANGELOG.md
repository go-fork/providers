# Changelog

## [Unreleased]

## [v0.1.0] - 2025-05-30

### Changed
- **BREAKING**: Migrated package module from `github.com/go-fork/providers/config` to `go.fork.vn/providers/config`
- Updated all import paths and references to use new module address
- This is a major version release establishing the new stable API under go.fork.vn domain

### Technical Details
- Module path changed to go.fork.vn/providers/config
- All documentation and examples updated to reflect new import path
- Maintains full API compatibility with v0.0.6

## [v0.0.6] - 2025-05-26

### Changed
- **BREAKING**: Replaced manual MockManager implementation with mockery-generated version
- Enhanced MockManager with full testify/mock integration and expecter pattern support
- Updated mocks/doc.go with comprehensive documentation and usage examples

### Added
- Added .mockery.yaml configuration file for consistent mock generation
- Support for both Expecter pattern (`EXPECT()` method) and Traditional mock pattern (`On()` method)
- Automatic panic when mock methods are called without proper expectations
- Better type safety and IDE support for mocks

### Technical Details
- MockManager now provides 1,748 lines of auto-generated mock code
- Full compatibility with testify/mock framework
- Improved test reliability with `AssertExpectations()` validation
- Easy mock regeneration with `mockery --name Manager` command

## [v0.0.5] - 2025-05-26

### Changed
- Upgraded go.fork.vn/di dependency from v0.0.4 to v0.0.5
- Implemented new interface methods from di.ServiceProvider: Requires() and Providers()
- Enhanced test coverage for ServiceProvider implementation
- Updated documentation in doc.go to reflect new interface methods

### Fixed
- Fixed staticcheck warning S1040 by removing redundant type assertion in tests

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
