# Release Notes v0.1.0

## Log Package v0.1.0 - First Stable Release

**Release Date**: May 30, 2025

### 🚀 Major Changes

#### Module Path Migration
- **BREAKING CHANGE**: Migrated from `github.com/go-fork/providers/log` to `go.fork.vn/providers/log`
- This marks the first stable release under the new domain structure
- All internal import paths have been updated to use the new module location

### 📈 What's Included

#### Core Features
- **Thread-safe logging manager** with multiple log levels (Debug, Info, Warning, Error, Fatal)
- **Multiple output handlers**:
  - ConsoleHandler with color support and formatting
  - FileHandler with automatic rotation based on file size
  - StackHandler for sending logs to multiple destinations simultaneously
- **Printf-style formatting** with placeholder support
- **Configurable minimum log levels** for each handler
- **Dependency Injection integration** through ServiceProvider interface
- **Extensible handler API** for custom log destinations

#### Advanced Capabilities
- **Resource management**: Automatic cleanup and proper handler closure
- **Performance optimized**: Efficient concurrent logging with minimal lock contention
- **Error resilience**: Individual handler failures don't crash the logging system
- **Flexible configuration**: Runtime handler management and dynamic reconfiguration

### 🔧 Migration Guide

#### For Existing Users (v0.0.x → v0.1.0)

1. **Update Import Statements**
   ```go
   // OLD
   import "github.com/go-fork/providers/log"
   import "github.com/go-fork/providers/log/handler"
   
   // NEW
   import "go.fork.vn/providers/log"
   import "go.fork.vn/providers/log/handler"
   ```

2. **Update go.mod**
   ```bash
   # Remove old dependency
   go mod edit -droprequire github.com/go-fork/providers/log
   
   # Add new dependency
   go get go.fork.vn/providers/log@v0.1.0
   
   # Clean up
   go mod tidy
   ```

3. **No Code Changes Required**
   - All APIs remain the same
   - No breaking changes to interfaces or method signatures
   - Configuration patterns unchanged

### 📚 Documentation

#### New Vietnamese Documentation
- **Technical Documentation**: `docs/log.md` - Comprehensive technical reference
- **Usage Guide**: `docs/usage.md` - Practical examples and best practices
- **API Reference**: Complete interface documentation with examples

#### Key Documentation Sections
- Architecture overview and component design
- Performance features and optimization tips
- Integration patterns with HTTP servers and background workers
- Testing utilities and mock handlers
- Troubleshooting common issues

### 🔍 Quality Assurance

#### Testing Coverage
- **100% test pass rate** - All 47 tests passing
- **Comprehensive test suite** covering:
  - Core logging functionality
  - Handler implementations (Console, File, Stack)
  - ServiceProvider integration
  - Error handling and edge cases
  - Concurrent access patterns

#### Code Quality
- **Static analysis**: Clean staticcheck results
- **Go vet**: No issues reported
- **Performance**: Optimized for high-throughput logging scenarios

### 🏗️ Technical Details

#### Dependencies
- `go.fork.vn/di v0.1.0` - Dependency injection framework
- **Go version**: 1.23.9+
- **External dependencies**: None (beyond DI framework)

#### Compatibility
- **Backward compatible** with all v0.0.x usage patterns
- **Go modules**: Full support with proper versioning
- **Cross-platform**: Tested on Linux, macOS, and Windows

### 💡 Usage Examples

#### Basic Setup
```go
manager := log.NewManager()
manager.SetMinLevel(handler.INFO)

consoleHandler := handler.NewConsoleHandler(true)
manager.AddHandler("console", consoleHandler)

manager.Info("Application started with log v0.1.0")
```

#### Production Configuration
```go
// File handler with 50MB rotation
fileHandler, _ := handler.NewFileHandler("/var/log/app.log", 50*1024*1024)
manager.AddHandler("file", fileHandler)

// Multiple destinations
stackHandler := handler.NewStackHandler()
stackHandler.AddHandler(fileHandler)
stackHandler.AddHandler(consoleHandler)
manager.AddHandler("main", stackHandler)
```

#### Dependency Injection
```go
provider := &log.ServiceProvider{}
provider.Register(app)
provider.Boot(app)

// Use from container
container.Call(func(manager log.Manager) {
    manager.Info("DI integration working perfectly")
})
```

### 🔗 Related Packages

This release coordinates with:
- `go.fork.vn/di v0.1.0` - Dependency injection framework
- Works seamlessly with other `go.fork.vn/providers/*` packages

### 🚨 Important Notes

#### Breaking Changes
- **Module path change** is the only breaking change
- **No API changes** - all existing code works after import updates
- **Version jump** from v0.0.5 to v0.1.0 indicates stability milestone

#### Recommendations
- **Update immediately** for new projects
- **Plan migration** for existing projects during next maintenance window
- **Test thoroughly** after migration (though no issues expected)

### 📞 Support

For issues, questions, or feedback:
- **GitHub Issues**: [go-fork/providers](https://github.com/go-fork/providers/issues)
- **Documentation**: Check `docs/` directory for comprehensive guides
- **Examples**: See `docs/usage.md` for practical implementation patterns

---

**Note**: This is a production-ready release suitable for use in mission-critical applications. The logging system has been thoroughly tested and is actively used in production environments.
