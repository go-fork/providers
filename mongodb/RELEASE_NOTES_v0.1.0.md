# Release Notes - MongoDB v0.1.0

**Release Date:** May 30, 2025  
**Module:** `go.fork.vn/providers/mongodb`

## 🚀 Major Changes

### Module Migration
- **BREAKING CHANGE**: Migrated module path from `github.com/go-fork/providers/mongodb` to `go.fork.vn/providers/mongodb`
- All import statements must be updated to use the new module path
- This marks the transition to the official go.fork.vn domain

## 📚 Documentation

### New Documentation Structure
- **Added** comprehensive documentation in `docs/` folder
- **Added** `docs/mongodb.md`: Complete technical documentation covering MongoDB manager interface, configuration options, SSL/TLS support, authentication, and performance tuning
- **Added** `docs/usage.md`: Comprehensive usage guide with practical examples for dependency injection, service layer integration, transactions, aggregation pipelines, and testing patterns
- **Updated** `CHANGELOG.md` with v0.1.0 release information

### Documentation Highlights
- Detailed MongoDB Manager interface with all connection management methods
- Complete configuration guide with environment variables support
- SSL/TLS configuration examples for secure connections
- Authentication setup for various MongoDB authentication mechanisms
- Dependency Injection integration patterns with go.fork.vn/di
- Service layer implementation examples with CRUD operations
- Advanced MongoDB operations including transactions and aggregation pipelines
- Comprehensive testing guide with mocks and integration tests
- Performance monitoring and health check implementations
- Connection pool optimization and best practices

## 🔧 Technical Details

### API Compatibility
- **Maintained** full API compatibility with v0.0.x
- **No breaking changes** to the Manager interface
- **No breaking changes** to ServiceProvider implementation
- All existing code will work without modifications (except import paths)

### Core Features (Unchanged)
- MongoDB connection management with automatic reconnection
- Connection pooling with configurable pool size and timeout settings
- SSL/TLS support for secure database connections
- Authentication support for various MongoDB auth mechanisms
- Integration with go.fork.vn/di dependency injection framework
- Type-safe database, collection, and client access methods
- Health checking and connection monitoring capabilities

## 📦 Installation

### New Installation Command
```bash
go get go.fork.vn/providers/mongodb@v0.1.0
```

### Import Statement
```go
import "go.fork.vn/providers/mongodb"
```

## 🔄 Migration Guide

### Update Import Statements
**Old:**
```go
import "github.com/go-fork/providers/mongodb"
```

**New:**
```go
import "go.fork.vn/providers/mongodb"
```

### Update go.mod
```go
module your-app

go 1.23

require (
    go.fork.vn/providers/mongodb v0.1.0
    // other dependencies...
)
```

### Find and Replace
Use your IDE's find and replace functionality:
- Find: `github.com/go-fork/providers/mongodb`
- Replace: `go.fork.vn/providers/mongodb`

## 📋 What's Next

### Upcoming Features (Future Releases)
- Enhanced connection pool monitoring and metrics
- Advanced aggregation pipeline helpers
- MongoDB change streams support
- Schema validation helpers
- Migration tools and utilities

## 🔗 Resources

- **Documentation**: See `docs/mongodb.md` and `docs/usage.md`
- **Examples**: Comprehensive examples available in `docs/usage.md`
- **API Reference**: Detailed in `docs/mongodb.md`
- **Changelog**: See `CHANGELOG.md` for complete history

## 🐛 Known Issues

None reported for this release.

## 🙏 Credits

This release maintains all existing functionality while establishing the foundation for future development under the go.fork.vn domain.

---

**Full Changelog**: https://github.com/go-fork/providers/compare/mongodb/v0.0.1...mongodb/v0.1.0
