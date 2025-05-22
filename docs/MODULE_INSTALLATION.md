# Module Installation Guide

## Using Go Modules

Each package in this repository is a separate Go module, and can be installed independently.

### Installing a specific version

To install a specific version of a module, use the standard `go get` command with the proper version:

```bash
# Install the config module at version v0.1.0
go get github.com/go-fork/providers/config@v0.1.0

# Install the cache module at version v0.2.0
go get github.com/go-fork/providers/cache@v0.2.0

# Install the log module at version v0.3.0
go get github.com/go-fork/providers/log@v0.3.0
```

### How versioning works in this repository

In our monorepo structure, each module has its own versioning. The Git tags for the modules follow this pattern:

```
MODULE_PATH-vX.Y.Z
```

For example:
- `config-v0.1.0` - For the config module
- `cache-v0.2.0` - For the cache module
- `log-v0.3.0` - For the log module

When using `go get`, you still use the standard version format:

```bash
# CORRECT
go get github.com/go-fork/providers/config@v0.1.0
```

This is because Go's module system automatically resolves the proper tag based on the module path.
