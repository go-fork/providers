# Go-Fork Providers

[![Go Report Card](https://goreportcard.com/badge/github.com/go-fork/providers)](https://goreportcard.com/report/github.com/go-fork/providers)
[![GoDoc](https://godoc.org/github.com/go-fork/providers?status.svg)](https://godoc.org/github.com/go-fork/providers)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A collection of modular Go service providers for the [PAMM Framework](https://github.com/go-fork/pamm). Each provider is a standalone module that can be used independently or in combination with other providers.

## Available Providers

- [**cache**](cache/): Caching abstraction layer with multiple driver implementations (Redis, MongoDB)
- [**config**](config/): Configuration management system with support for various formats (YAML, JSON, ENV)
- [**http**](http/): HTTP server and client with middleware support
- [**log**](log/): Logging subsystem with multiple handlers and formatters
- [**middleware**](middleware/): Collection of HTTP middleware components
- [**mailer**](mailer/): Email sending abstraction
- [**queue**](queue/): Message queue implementation
- [**scheduler**](scheduler/): Task scheduling system
- [**sms**](sms/): SMS messaging service

## Installation

Each provider is a separate Go module that can be installed independently:

```bash
# Install cache provider
go get github.com/go-fork/providers/cache

# Install config provider
go get github.com/go-fork/providers/config

# Install http provider
go get github.com/go-fork/providers/http

# Install log provider
go get github.com/go-fork/providers/log

# Install middleware
go get github.com/go-fork/providers/middleware/*
```

## Basic Usage

### Cache Provider

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/go-fork/providers/cache"
)

func main() {
    // Create a new cache manager with default configuration
    cacheManager, err := cache.NewManager().Build()
    if err != nil {
        log.Fatalf("Failed to create cache manager: %v", err)
    }

    ctx := context.Background()
    
    // Store a value in cache
    err = cacheManager.Set(ctx, "key", "value", 5*time.Minute)
    if err != nil {
        log.Fatalf("Failed to set cache: %v", err)
    }
    
    // Retrieve a value from cache
    var value string
    exists, err := cacheManager.Get(ctx, "key", &value)
    if err != nil {
        log.Fatalf("Failed to get from cache: %v", err)
    }
    
    if exists {
        log.Printf("Value from cache: %s", value)
    } else {
        log.Println("Value not found in cache")
    }
}
```

### HTTP Provider

```go
package main

import (
    "log"
    "net/http"

    httpProvider "github.com/go-fork/providers/http"
)

func main() {
    // Create a new HTTP application
    app, err := httpProvider.NewApplication()
    if err != nil {
        log.Fatalf("Failed to create HTTP application: %v", err)
    }
    
    // Register a route
    app.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message":"Hello, World!"}`))
    })
    
    // Start the server
    log.Println("Server starting on :8080")
    app.Start(":8080")
}
```

## Features

- **Modular Design**: Each provider is a separate module that can be used independently.
- **Dependency Injection**: Fully compatible with the [go-fork/di](https://github.com/go-fork/di) container.
- **Configurable**: Extensive configuration options for each provider.
- **Extensible**: Easy to extend with custom implementations.
- **Well Tested**: Comprehensive test suite for each provider.
- **Production Ready**: Used in production applications.

## Documentation

- For detailed documentation on each provider, see the README files in their respective directories.
- Examples can be found in the [examples](examples/) directory.
- API documentation is available on [GoDoc](https://godoc.org/github.com/go-fork/providers).
- For release process information, see [RELEASE_PROCESS.md](docs/RELEASE_PROCESS.md).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Laravel's Service Provider pattern
- Built for the [PAMM Framework](https://github.com/go-fork/pamm)
