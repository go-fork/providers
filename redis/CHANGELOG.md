# CHANGELOG.md - Redis Provider

All notable changes to the Redis provider will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-05-30

### Changed
- **BREAKING CHANGE**: Migrated module path from `github.com/go-fork/providers/redis` to `go.fork.vn/providers/redis`
- Updated all import statements to use the new module path
- Updated dependencies to use new go.fork.vn domain

### Added
- Comprehensive documentation in `docs/` folder
- Added `docs/redis.md`: Complete technical documentation covering architecture, interfaces, and implementation details
- Added `docs/usage.md`: Practical usage guide with examples and best practices
- Migration guide in release notes

### Technical Details
- All API compatibility maintained with v0.0.1
- No breaking changes to the Manager interface or ServiceProvider implementation
- Full backward compatibility (except import paths)

## [0.0.1] - 2025-05-26

### Added
- Initial release
- Support for standard Redis client configuration
- Support for Redis Universal client (for Cluster, Sentinel, and standalone modes)
- Service provider integration with DI container
- Configuration support through config provider
- Improved test coverage (>75%)
- Mockery integration for easier testing
- Fixed error message capitalization
- Added nil checks for better error handling
