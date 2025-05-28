# Go-Fork Log Provider

Gói `log` cung cấp hệ thống logging hiện đại, linh hoạt và mở rộng cho ứng dụng Go. Gói này được thiết kế để thread-safe và dễ dàng tích hợp vào các ứng dụng Go hiện đại với khả năng xử lý đa dạng các loại output.

## Giới thiệu

Logging là một phần thiết yếu trong phát triển và vận hành ứng dụng. Gói `log` cung cấp một hệ thống logging đơn giản nhưng mạnh mẽ với khả năng phân loại theo mức độ nghiêm trọng và hỗ trợ nhiều handlers khác nhau. Được thiết kế thread-safe, hệ thống này đảm bảo ghi log an toàn trong môi trường concurrent.

## Tính năng nổi bật

- **Đa dạng cấp độ log**: Hỗ trợ các cấp độ từ Debug, Info, Warning, Error đến Fatal.
- **Output đa dạng**: Hỗ trợ ghi log ra console (có màu) và file, dễ dàng mở rộng với custom handlers.
- **Thread-safe**: An toàn khi ghi log từ nhiều goroutines cùng lúc.
- **Hỗ trợ định dạng**: Hỗ trợ chuỗi định dạng kiểu Printf trong các thông điệp log.
- **Xử lý file linh hoạt**: Tự động xoay vòng file log khi đạt kích thước giới hạn.
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container, tương thích đầy đủ với di v0.0.5.
- **Cấu trúc mở rộng**: Dễ dàng triển khai handler mới cho các output khác.
- **Truy xuất handler linh hoạt**: Lấy handler đã đăng ký để cấu hình hoặc tùy chỉnh thêm.
- **Tuân thủ interface ServiceProvider**: Triển khai đầy đủ các methods Requires và Providers cho di v0.0.5.

## Cấu trúc package

```
log/
  ├── doc.go                 # Tài liệu tổng quan về package
  ├── manager.go             # Định nghĩa interface Manager và DefaultManager
  ├── provider.go            # ServiceProvider tích hợp với DI container
  └── handler/
      ├── handler.go         # Định nghĩa interface Handler và cấp độ log
      ├── console.go         # Console handler (hỗ trợ màu sắc)
      ├── file.go            # File handler (hỗ trợ xoay vòng)
      └── stack.go           # Stack handler (ghi log cho nhiều handler)
```

## Cách hoạt động

### Đăng ký Service Provider

Service Provider cho phép tích hợp dễ dàng gói `log` vào ứng dụng sử dụng DI container:

```go
// Trong file bootstrap của ứng dụng
import "github.com/go-fork/providers/log"

func bootstrap(app interface{}) {
    // Đăng ký log provider
    logProvider := log.NewServiceProvider()
    logProvider.Register(app)
    
    // Boot các providers sau khi tất cả đã đăng ký
    logProvider.Boot(app)
}
```

ServiceProvider sẽ tự động:
1. Tạo một log manager mới
2. Cấu hình console handler với màu sắc
3. Cấu hình file handler trong thư mục storage/logs
4. Đăng ký manager vào container với key "log"

### Sử dụng trực tiếp

Bạn có thể tạo và sử dụng log manager mà không cần thông qua DI container:

```go
// Tạo manager mới
manager := log.NewManager()

// Thêm console handler có màu
consoleHandler := handler.NewConsoleHandler(true)
manager.AddHandler("console", consoleHandler)

// Thêm file handler với kích thước tối đa 10MB
fileHandler, err := handler.NewFileHandler("app.log", 10*1024*1024)
if err == nil {
    manager.AddHandler("file", fileHandler)
}

// Bắt đầu ghi log
manager.Debug("Khởi động ứng dụng")
manager.Info("Cấu hình đã được nạp: %v", config)
manager.Warning("Tài nguyên cao: %d%%", usagePercent)
manager.Error("Lỗi kết nối: %v", err)
manager.Fatal("Không thể khởi tạo database")
```

### Sử dụng các Handlers

#### Console Handler

Handler này ghi log ra terminal với hỗ trợ màu sắc:

```go
// Tạo console handler với màu sắc
consoleHandler := handler.NewConsoleHandler(true)

// Không sử dụng màu
plainConsole := handler.NewConsoleHandler(false)

// Thiết lập cấp độ log tối thiểu (bỏ qua Debug)
consoleHandler.SetMinLevel(handler.InfoLevel)
```

#### File Handler

Handler này ghi log ra file với hỗ trợ xoay vòng:

```go
// Tạo file handler với kích thước tối đa 5MB
fileHandler, err := handler.NewFileHandler("/path/to/app.log", 5*1024*1024)
if err != nil {
    // Xử lý lỗi
}

// Thiết lập cấp độ log tối thiểu (bỏ qua Debug và Info)
fileHandler.SetMinLevel(handler.WarningLevel)
```

#### Stack Handler

Handler này chuyển tiếp log đến nhiều handlers khác:

```go
// Tạo các handlers riêng lẻ
consoleHandler := handler.NewConsoleHandler(true)
fileHandler, _ := handler.NewFileHandler("app.log", 10*1024*1024)

// Tạo stack handler
stackHandler := handler.NewStackHandler()
stackHandler.PushHandler(consoleHandler)
stackHandler.PushHandler(fileHandler)

// Bây giờ chỉ cần thêm stack handler vào manager
manager.AddHandler("combined", stackHandler)
```

### Truy xuất và tùy chỉnh Handler

Bạn có thể truy xuất các handler đã đăng ký để thực hiện cấu hình bổ sung hoặc kiểm tra trạng thái:

```go
// Thêm handler vào manager
consoleHandler := handler.NewConsoleHandler(true)
manager.AddHandler("console", consoleHandler)

// Sau đó, truy xuất handler để cấu hình thêm
if handlerObj := manager.GetHandler("console"); handlerObj != nil {
    // Chuyển đổi kiểu để truy cập các phương thức đặc thù cho console handler
    if typedHandler, ok := handlerObj.(*handler.ConsoleHandler); ok {
        // Thực hiện cấu hình bổ sung...
    }
}

// Kiểm tra handler có tồn tại không trước khi xóa
if manager.GetHandler("old-handler") != nil {
    manager.RemoveHandler("old-handler")
}
```

### Lọc theo cấp độ

Mỗi handler có thể được cấu hình để chỉ xử lý các log từ cấp độ nghiêm trọng nhất định trở lên:

```go
// Tạo handler
consoleHandler := handler.NewConsoleHandler(true)

// Chỉ ghi log từ Warning trở lên (Warning, Error, Fatal)
// Bỏ qua Debug và Info
consoleHandler.SetMinLevel(handler.WarningLevel)
```

### Đóng handlers đúng cách

Luôn đóng manager khi không sử dụng để đảm bảo tất cả các handlers được dọn dẹp đúng cách:

```go
manager := log.NewManager()
// ... cấu hình và sử dụng manager ...

// Đóng manager khi kết thúc
defer manager.Close()
```

### Tạo Custom Handler

Bạn có thể triển khai handler của riêng mình bằng cách tuân thủ interface Handler:

```go
type MyCustomHandler struct {
    minLevel handler.Level
}

func (h *MyCustomHandler) Handle(level handler.Level, message string) {
    if level < h.minLevel {
        return
    }
    // Triển khai xử lý log theo cách của bạn
}

func (h *MyCustomHandler) SetMinLevel(level handler.Level) {
    h.minLevel = level
}

func (h *MyCustomHandler) Close() error {
    // Dọn dẹp tài nguyên nếu cần
    return nil
}
```

### Lấy Handler Đã Đăng Ký

Để lấy handler đã đăng ký và thực hiện cấu hình hoặc tùy chỉnh thêm:

```go
// Giả sử bạn đã đăng ký một handler với tên "file"
fileHandler := manager.GetHandler("file")

// Thực hiện cấu hình cho handler
if fileHandler != nil {
    fileHandler.SetMinLevel(handler.ErrorLevel)
}
```

---

Để biết thêm thông tin chi tiết và API reference, vui lòng xem tài liệu trong file `doc.go` hoặc chạy lệnh `go doc github.com/go-fork/providers/log`.
