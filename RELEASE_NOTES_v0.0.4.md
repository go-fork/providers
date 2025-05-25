# Release Notes v0.0.4

## Go-Fork Providers v0.0.4

**Release Date**: 2025-05-26

## 🆕 New Features

### CI/CD Improvements
- **GitHub Actions Optimization**: Completely redesigned workflows for monorepo structure
- **Module Release Workflow**: Added dedicated workflow for individual module releases (e.g., `cache/v1.0.0`)
- **Automated Testing**: Enhanced CI to test each module separately with proper isolation

### Documentation Enhancement
- **Module Installation Guide**: Comprehensive documentation for installing individual modules
- **Release Process**: Detailed documentation of the release workflow for maintainers
- **README Updates**: Improved main README with proper module structure information

### Build System
- **GoReleaser Configuration**: Updated for monorepo support with proper module handling
- **Archive Generation**: Improved packaging of individual modules and documentation

## 🔄 Changes

### Workflow Optimization
- **Removed Unnecessary Workflows**: Eliminated benchmark, examples, metrics-report, and update-readme workflows
- **Streamlined CI**: Updated continuous integration to focus on essential testing and validation
- **Module-Specific Testing**: Each Go module is now tested independently with its own dependencies

### Configuration Updates
- **GoReleaser**: Refined configuration to exclude non-existent modules
- **Lint Integration**: Added proper golangci-lint installation for code quality checks

## 🐛 Bug Fixes

### Module Dependencies
- **Version Consistency**: Fixed versioning inconsistencies across go.mod files
- **Dependency Updates**: Resolved dependency conflicts between modules

### CI/CD Fixes
- **Workflow Validation**: Fixed issues with GitHub Actions workflow syntax
- **Build Process**: Resolved GoReleaser configuration validation errors

## 📦 Modules Included

This release includes the following stable modules:

- **cache** (`github.com/go-fork/providers/cache`) - Caching abstraction layer
- **config** (`github.com/go-fork/providers/config`) - Configuration management
- **log** (`github.com/go-fork/providers/log`) - Logging utilities
- **mailer** (`github.com/go-fork/providers/mailer`) - Email sending capabilities
- **queue** (`github.com/go-fork/providers/queue`) - Queue management system
- **scheduler** (`github.com/go-fork/providers/scheduler`) - Task scheduling
- **sms** (`github.com/go-fork/providers/sms`) - SMS messaging

## 🚀 Installation

### Full Package
```bash
go get github.com/go-fork/providers@v0.0.4
```

### Individual Modules
```bash
go get github.com/go-fork/providers/cache@v0.0.4
go get github.com/go-fork/providers/config@v0.0.4
go get github.com/go-fork/providers/log@v0.0.4
go get github.com/go-fork/providers/mailer@v0.0.4
go get github.com/go-fork/providers/queue@v0.0.4
go get github.com/go-fork/providers/scheduler@v0.0.4
go get github.com/go-fork/providers/sms@v0.0.4
```

## 📋 Requirements

- **Go Version**: 1.23.9 or later
- **Dependencies**: See individual module go.mod files for specific requirements

## 🔗 Links

- **Repository**: https://github.com/go-fork/providers
- **Documentation**: https://github.com/go-fork/providers/docs
- **Issues**: https://github.com/go-fork/providers/issues

## 👥 Contributors

Thanks to all contributors who made this release possible!

---

For detailed technical changes, see [CHANGELOG.md](./CHANGELOG.md)
