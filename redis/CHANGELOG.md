# CHANGELOG.md - Redis Provider

All notable changes to the Redis provider will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
