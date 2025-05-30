# Release Notes - Config v0.1.0

**Release Date:** May 30, 2025  
**Module:** `go.fork.vn/providers/config`

## 🚀 Major Changes

### Module Migration
- **BREAKING CHANGE**: Migrated module path from `github.com/go-fork/providers/config` to `go.fork.vn/providers/config`
- All import statements must be updated to use the new module path
- This marks the transition to the official go.fork.vn domain

## 📚 Documentation

### New Documentation Structure
- **Added** comprehensive documentation in `docs/` folder
- **Added** `docs/config.md`: Complete technical documentation covering architecture, interfaces, and implementation details
- **Added** `docs/usage.md`: Practical usage guide with examples and best practices
- **Updated** `CHANGELOG.md` with v0.1.0 release information

### Documentation Highlights
- Detailed API reference with all methods of the Manager interface
- Type-safe operations examples and best practices
- Environment variable configuration patterns
- Configuration watching and hot-reload examples
- Struct unmarshaling with mapstructure tags
- Dependency Injection integration guide
- Error handling and validation patterns

## 🔧 Technical Details

### API Compatibility
- **Maintained** full API compatibility with v0.0.6
- **No breaking changes** to the Manager interface
- **No breaking changes** to ServiceProvider implementation
- All existing code will work without modifications (except import paths)

### Core Features (Unchanged)
- Type-safe configuration value retrieval with (value, exists) return pattern
- Automatic environment variable support
- Configuration file watching and hot-reload
- Struct unmarshaling with mapstructure
- Integration with go.fork.vn/di dependency injection
- Support for multiple configuration sources and formats

## 📦 Installation

### New Installation Command
```bash
go get go.fork.vn/providers/config@v0.1.0
```

### Import Statement
```go
import "go.fork.vn/providers/config"
```

## 🔄 Migration Guide

### Update Import Statements
**Old:**
```go
import "github.com/go-fork/providers/config"
```

**New:**
```go
import "go.fork.vn/providers/config"
```

### Update go.mod
```go
module your-app

go 1.23

require (
    go.fork.vn/providers/config v0.1.0
    // other dependencies...
)
```

### Find and Replace
Use your IDE's find and replace functionality:
- Find: `github.com/go-fork/providers/config`
- Replace: `go.fork.vn/providers/config`

## 📋 What's Next

### Upcoming Features (Future Releases)
- Enhanced validation framework
- Configuration schema support
- More configuration source adapters
- Performance optimizations

## 🔗 Resources

- **Documentation**: See `docs/config.md` and `docs/usage.md`
- **Examples**: Available in `docs/usage.md`
- **API Reference**: Detailed in `docs/config.md`
- **Changelog**: See `CHANGELOG.md` for complete history

## 🐛 Known Issues

None reported for this release.

## 🙏 Credits

This release maintains all existing functionality while establishing the foundation for future development under the go.fork.vn domain.

---

**Full Changelog**: https://github.com/go-fork/providers/compare/config/v0.0.6...config/v0.1.0
