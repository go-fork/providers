# Go-Fork Providers

[![Go Report Card](https://goreportcard.com/badge/github.com/go-fork/providers)](https://goreportcard.com/report/github.com/go-fork/providers)
[![GoDoc](https://god## Tương thích giữa các phiên bản

| Provider | v0.0.2 | v0.0.3 | di v0.0.2 | di v0.0.3 |
|----------|--------|--------|-----------|-----------|
| cache    | ✅     | -      | ✅        | ✅        |
| config   | ✅     | ✅     | ✅        | ✅        |
| log      | ✅     | ✅     | ✅        | ✅        |
| mailer   | ✅     | ✅     | ✅        | ✅        |
| queue    | ✅     | ✅     | ✅        | ✅        |
| scheduler| -      | ✅     | ❌        | ✅        |
| sms      | ✅     | -      | ✅        | ❌        |b.com/go-fork/providers?status.svg)](https://godoc.org/github.com/go-fork/providers)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A collection of modular Go service providers for the [Fork Framework](https://github.com/go-fork). Each provider là một module độc lập có thể sử dụng riêng lẻ hoặc kết hợp với các provider khác. Lưu ý: HTTP và Middleware providers đã được chuyển sang repository riêng.

## Available Providers

Phiên bản hiện tại: **v0.0.3** (Phát hành ngày: 25/05/2025)

| Provider | Phiên bản | Mô tả | Tính năng nổi bật |
|----------|-----------|-------|-------------------|
| [**cache**](cache/) | v0.0.2 | Hệ thống cache đa lớp | Redis, MongoDB, Memory, File drivers |
| [**config**](config/) | v0.0.3 | Quản lý cấu hình ứng dụng | YAML, JSON, ENV, dot notation |
| [**log**](log/) | v0.0.3 | Hệ thống logging | Console, File, Stack handlers |
| [**mailer**](mailer/) | v0.0.3 | Gửi email | SMTP, Template, Attachments |
| [**queue**](queue/) | v0.0.3 | Hàng đợi xử lý tác vụ | Memory, Redis adapters |
| [**scheduler**](scheduler/) | v0.0.3 | Lập lịch tác vụ | Cron expressions, Redis locking |
| [**sms**](sms/) | v0.0.2 | Dịch vụ gửi SMS | Multiple provider support |

## Các cập nhật trong phiên bản v0.0.3

### Tính năng mới
- **scheduler**: Hệ thống lập lịch tác vụ mới với khả năng tích hợp Redis
- **queue**: Hoàn thiện chức năng worker với tích hợp scheduler
- **mailer**: Cập nhật từ `github.com/go-gomail/gomail` sang `gopkg.in/gomail.v2`
- **config**: Bổ sung các API tiện ích và MockManager cho unit testing

### Cải tiến
- Nâng cấp toàn bộ các dependency lõi lên phiên bản mới nhất
- Cải thiện performance và bảo mật
- Tăng cường độ bao phủ của test suite

### Lưu ý
- **http** và **middleware** đã được chuyển sang repository riêng

### Phiên bản trước đó
Xem chi tiết về các phiên bản trước tại [CHANGELOG.md](CHANGELOG.md)

## Installation

Mỗi provider là một Go module độc lập và có thể được cài đặt riêng biệt:

```bash
# Cài đặt cache provider
go get github.com/go-fork/providers/cache@v0.0.2

# Cài đặt config provider
go get github.com/go-fork/providers/config@v0.0.3

# Cài đặt mailer provider
go get github.com/go-fork/providers/mailer@v0.0.3

# Cài đặt queue provider
go get github.com/go-fork/providers/queue@v0.0.3

# Cài đặt scheduler provider
go get github.com/go-fork/providers/scheduler@v0.0.3

# Cài đặt log provider
go get github.com/go-fork/providers/log@v0.0.3

# Cài đặt sms provider
go get github.com/go-fork/providers/sms@v0.0.2
```

### Cài đặt tất cả providers

```bash
# Cài đặt tất cả các providers với phiên bản mới nhất
go get github.com/go-fork/providers/...@latest
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

```
```

## Tính năng chung

- **Thiết kế module hóa**: Mỗi provider là một module riêng biệt, có thể sử dụng độc lập.
- **Dependency Injection**: Tương thích hoàn toàn với [go-fork/di](https://github.com/go-fork/di) v0.0.3.
- **Khả năng cấu hình cao**: Mỗi provider có nhiều tùy chọn cấu hình mở rộng.
- **Dễ dàng mở rộng**: Các interface tiện lợi để triển khai tùy chỉnh.
- **Kiểm thử kỹ lưỡng**: Bộ test suite toàn diện cho mỗi provider.
- **Sẵn sàng cho production**: Đã được sử dụng trong các ứng dụng thực tế.
- **Tài liệu đầy đủ**: API docs, README và ví dụ chi tiết.

## Roadmap & Phát triển

| Provider | Phiên bản tiếp theo | Kế hoạch phát triển |
|----------|---------------------|---------------------|
| cache    | v0.0.3              | Thêm Memcached driver, cải thiện hiệu suất |
| config   | v0.0.4              | Hỗ trợ HCL format, cải thiện remote config |
| mailer   | v0.0.4              | Thêm queue integration, HTML template cải tiến |
| queue    | v0.0.4              | Kafka adapter, DLQ (Dead Letter Queue) |
| scheduler| v0.0.4              | Cải thiện distributed scheduling |
| sms      | v0.0.3              | Tích hợp nhiều nhà cung cấp SMS mới |

## Chi tiết về các Provider

### Cache Provider (v0.0.2)
Hệ thống cache với multiple driver:
- **Memory**: Bộ nhớ tạm trong ứng dụng với auto cleanup
- **File**: Lưu trữ cache vào hệ thống file
- **Redis**: Tích hợp đầy đủ với Redis v9+
- **MongoDB**: Lưu trữ cache phân tán với MongoDB

### Config Provider (v0.0.3)
Hệ thống quản lý cấu hình:
- Hỗ trợ YAML, JSON, ENV
- Truy cập dot notation (app.database.host)
- Type-safe accessors (GetString, GetInt...)
- Binding với struct tự động
- Watcher cho file cấu hình

### Log Provider (v0.0.3)
Hệ thống logging:
- Multiple severity levels (Debug, Info, Warning, Error, Fatal)
- Console handler với ANSI color support
- File handler với rotation tự động
- Stack trace capture tự động khi error

### Mailer Provider (v0.0.3)
Giải pháp gửi email:
- SMTP support với cấu hình linh hoạt
- Template rendering (HTML, Text)
- File đính kèm và inline images
- Queue integration để gửi bất đồng bộ

### Queue Provider (v0.0.3)
Hệ thống xử lý hàng đợi:
- Memory và Redis adapters
- Delayed/scheduled jobs
- Retry và error handling
- Priority queues
- Worker pool management

### Scheduler Provider (v0.0.3)
Hệ thống lập lịch tác vụ:
- Cron expressions
- Distributed locking với Redis
- One-time và recurring tasks
- Callbacks và error handling

### SMS Provider (v0.0.2)
Dịch vụ gửi tin nhắn SMS:
- Support nhiều provider SMS
- Template support
- Queue integration
- Retry mechanism

## Tài liệu

- Tài liệu chi tiết về mỗi provider có trong thư mục tương ứng (README.md).
- Ví dụ có thể tìm thấy trong thư mục [examples](examples/).
- Tài liệu API có sẵn trên [GoDoc](https://godoc.org/github.com/go-fork/providers).
- Thông tin quy trình phát hành, xem [RELEASE_PROCESS.md](docs/RELEASE_PROCESS.md).
- Lịch sử thay đổi chi tiết tại [CHANGELOG.md](CHANGELOG.md).

## Tương thích giữa các phiên bản

| Provider | v0.0.2 | v0.0.3 | di v0.0.2 | di v0.0.3 |
|----------|--------|--------|-----------|-----------|
| cache    | ✅     | -      | ✅        | ✅        |
| config   | ✅     | ✅     | ✅        | ✅        |
| http     | ✅     | -      | ✅        | ❌        |
| log      | ✅     | ✅     | ✅        | ✅        |
| mailer   | ✅     | ✅     | ✅        | ✅        |
| queue    | ✅     | ✅     | ✅        | ✅        |
| scheduler| -      | ✅     | ❌        | ✅        |
| sms      | ❌     | ❌      | ❌        | ❌        |

✅ = Tương thích đầy đủ
❌ = Không tương thích
\- = Không có phiên bản

## Đóng góp

Mọi đóng góp đều được chào đón! Vui lòng gửi Pull Request.

1. Fork repository
2. Tạo branch tính năng của bạn (`git checkout -b feature/amazing-feature`)
3. Commit các thay đổi (`git commit -m 'Add some amazing feature'`)
4. Push đến branch (`git push origin feature/amazing-feature`)
5. Mở Pull Request

## Giấy phép

Dự án này được cấp phép theo giấy phép MIT - xem file [LICENSE](LICENSE) để biết chi tiết.

## Ghi công

- Lấy cảm hứng từ mẫu Service Provider
- Được xây dựng cho [Fork Framework](https://github.com/go-fork)
- Phát triển bởi cộng đồng Go-Fork
