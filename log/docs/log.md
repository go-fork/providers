# Log Package - Tài liệu Kỹ thuật

## Tổng quan

Package `go.fork.vn/providers/log` cung cấp một hệ thống logging linh hoạt và mạnh mẽ cho ứng dụng Go, với khả năng tích hợp dependency injection và hỗ trợ nhiều output handlers.

## Kiến trúc

### Core Components

#### 1. Manager Interface
```go
type Manager interface {
    Debug(message string, args ...interface{}) error
    Info(message string, args ...interface{}) error
    Warning(message string, args ...interface{}) error
    Error(message string, args ...interface{}) error
    Fatal(message string, args ...interface{}) error
    AddHandler(name string, handler handler.Handler) error
    RemoveHandler(name string) error
    GetHandler(name string) (handler.Handler, bool)
    SetMinLevel(level handler.Level)
    Close() error
}
```

#### 2. Handler System
Package hỗ trợ 3 loại handler chính:

**ConsoleHandler**: Xuất log ra console với hỗ trợ màu sắc
```go
handler := handler.NewConsoleHandler(true) // enable colors
```

**FileHandler**: Ghi log vào file với tính năng rotation
```go
handler, err := handler.NewFileHandler("/path/to/log.txt", 10*1024*1024) // 10MB max size
```

**StackHandler**: Gửi log đến nhiều handlers cùng lúc
```go
stackHandler := handler.NewStackHandler()
stackHandler.AddHandler(consoleHandler)
stackHandler.AddHandler(fileHandler)
```

### 3. ServiceProvider
Tích hợp với DI container thông qua interface `di.ServiceProvider`:

```go
type ServiceProvider struct{}

func (p *ServiceProvider) Register(app interface{}) error
func (p *ServiceProvider) Boot(app interface{}) error
func (p *ServiceProvider) Requires() []string
func (p *ServiceProvider) Providers() []string
```

## Log Levels

Package định nghĩa 5 mức độ log:

```go
const (
    DEBUG Level = iota
    INFO
    WARNING
    ERROR
    FATAL
)
```

## Thread Safety

- **Manager**: Thread-safe, có thể được sử dụng từ nhiều goroutines
- **Handlers**: Tất cả handlers đều thread-safe
- **File rotation**: An toàn với concurrent writes

## Performance Features

### 1. Lazy Formatting
Chỉ format message khi thực sự cần thiết (handler level >= min level)

### 2. Efficient File Rotation
- Size-based rotation
- Automatic backup file management
- Minimal lock contention

### 3. Memory Management
- Proper resource cleanup
- Handler replacement without leaks
- Automatic closure on manager shutdown

## Error Handling

### Handler Errors
```go
// Lỗi từ individual handlers không làm crash toàn bộ logging system
manager.Info("message") // Tiếp tục hoạt động ngay cả khi một handler lỗi
```

### File Operations
```go
// Graceful handling của file permissions, disk space issues
fileHandler, err := handler.NewFileHandler("/readonly/path", 1024)
if err != nil {
    // Xử lý lỗi không thể tạo file
}
```

## Configuration Patterns

### 1. Basic Setup
```go
manager := log.NewManager()
manager.SetMinLevel(handler.INFO)

consoleHandler := handler.NewConsoleHandler(true)
manager.AddHandler("console", consoleHandler)
```

### 2. Production Setup
```go
manager := log.NewManager()
manager.SetMinLevel(handler.WARNING)

// File handler cho production logs
fileHandler, _ := handler.NewFileHandler("/var/log/app.log", 50*1024*1024)
manager.AddHandler("file", fileHandler)

// Console handler cho development
if os.Getenv("ENV") == "development" {
    consoleHandler := handler.NewConsoleHandler(true)
    manager.AddHandler("console", consoleHandler)
}
```

### 3. Multi-destination Logging
```go
stackHandler := handler.NewStackHandler()

// Add multiple destinations
errorFile, _ := handler.NewFileHandler("/var/log/errors.log", 10*1024*1024)
generalFile, _ := handler.NewFileHandler("/var/log/general.log", 100*1024*1024)
console := handler.NewConsoleHandler(true)

stackHandler.AddHandler(errorFile)
stackHandler.AddHandler(generalFile) 
stackHandler.AddHandler(console)

manager.AddHandler("stack", stackHandler)
```

## Advanced Features

### 1. Custom Handlers
Implement interface `handler.Handler`:
```go
type CustomHandler struct{}

func (h *CustomHandler) Log(level handler.Level, message string, args ...interface{}) error {
    // Custom implementation
    return nil
}

func (h *CustomHandler) Close() error {
    // Cleanup resources
    return nil
}
```

### 2. Dynamic Configuration
```go
// Runtime handler management
manager.RemoveHandler("console")
newHandler := handler.NewConsoleHandler(false) // disable colors
manager.AddHandler("console", newHandler)
```

### 3. Conditional Logging
```go
if manager.GetMinLevel() <= handler.DEBUG {
    expensiveDebugInfo := generateDebugInfo()
    manager.Debug("Debug info: %s", expensiveDebugInfo)
}
```

## Testing Support

### Mock Handlers
```go
type MockHandler struct {
    LogCalled   bool
    CloseCalled bool
    Messages    []string
}

func (m *MockHandler) Log(level handler.Level, message string, args ...interface{}) error {
    m.LogCalled = true
    m.Messages = append(m.Messages, fmt.Sprintf(message, args...))
    return nil
}
```

### Test Utilities
```go
func TestMyFunction(t *testing.T) {
    manager := log.NewManager()
    mockHandler := &MockHandler{}
    manager.AddHandler("mock", mockHandler)
    
    // Test your function
    myFunction(manager)
    
    // Verify logging behavior
    assert.True(t, mockHandler.LogCalled)
    assert.Contains(t, mockHandler.Messages[0], "expected message")
}
```

## Best Practices

### 1. Resource Management
```go
defer manager.Close() // Đảm bảo cleanup resources
```

### 2. Error Context
```go
manager.Error("Database connection failed: %v", err)
// Thay vì chỉ: manager.Error("Database error")
```

### 3. Structured Logging
```go
manager.Info("User action: user_id=%d, action=%s, timestamp=%s", 
    userID, action, time.Now().Format(time.RFC3339))
```

### 4. Performance Considerations
```go
// Tránh expensive operations trong log messages khi không cần thiết
if manager.GetMinLevel() <= handler.DEBUG {
    manager.Debug("Complex state: %+v", buildComplexState())
}
```

## Integration Examples

### Với HTTP Server
```go
func loggingMiddleware(manager log.Manager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            manager.Info("Request started: %s %s", r.Method, r.URL.Path)
            
            next.ServeHTTP(w, r)
            
            manager.Info("Request completed: %s %s (%v)", 
                r.Method, r.URL.Path, time.Since(start))
        })
    }
}
```

### Với Background Workers
```go
func worker(manager log.Manager) {
    for {
        select {
        case task := <-taskChannel:
            manager.Debug("Processing task: %+v", task)
            if err := processTask(task); err != nil {
                manager.Error("Task failed: %v", err)
            } else {
                manager.Info("Task completed successfully")
            }
        }
    }
}
```
