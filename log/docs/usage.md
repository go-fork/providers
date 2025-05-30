# Log Package - Hướng dẫn Sử dụng

## Cài đặt

```bash
go get go.fork.vn/providers/log@v0.1.0
```

## Sử dụng Cơ bản

### 1. Import Package
```go
import (
    "go.fork.vn/providers/log"
    "go.fork.vn/providers/log/handler"
)
```

### 2. Tạo Manager và Handler
```go
// Tạo manager
manager := log.NewManager()

// Tạo console handler với màu sắc
consoleHandler := handler.NewConsoleHandler(true)
manager.AddHandler("console", consoleHandler)

// Thiết lập mức log tối thiểu
manager.SetMinLevel(handler.INFO)
```

### 3. Ghi Log
```go
manager.Debug("This won't be shown (below INFO level)")
manager.Info("Application started")
manager.Warning("This is a warning: %s", "low disk space")
manager.Error("Database connection failed: %v", err)
manager.Fatal("Critical error occurred")
```

## Các Handler Phổ biến

### Console Handler
```go
// Với màu sắc (development)
consoleHandler := handler.NewConsoleHandler(true)
manager.AddHandler("console", consoleHandler)

// Không màu sắc (production)
consoleHandler := handler.NewConsoleHandler(false)
manager.AddHandler("console", consoleHandler)
```

### File Handler
```go
// File handler với rotation (50MB max)
fileHandler, err := handler.NewFileHandler("/var/log/app.log", 50*1024*1024)
if err != nil {
    log.Fatal("Cannot create file handler:", err)
}
manager.AddHandler("file", fileHandler)
```

### Stack Handler (Multiple Outputs)
```go
stackHandler := handler.NewStackHandler()

// Thêm console
console := handler.NewConsoleHandler(true)
stackHandler.AddHandler(console)

// Thêm file
file, _ := handler.NewFileHandler("/var/log/app.log", 10*1024*1024)
stackHandler.AddHandler(file)

// Đăng ký stack handler
manager.AddHandler("main", stackHandler)
```

## Cấu hình cho Các Môi trường

### Development Environment
```go
func setupDevelopmentLogging() log.Manager {
    manager := log.NewManager()
    manager.SetMinLevel(handler.DEBUG)
    
    // Console với màu sắc
    consoleHandler := handler.NewConsoleHandler(true)
    manager.AddHandler("console", consoleHandler)
    
    return manager
}
```

### Production Environment
```go
func setupProductionLogging() log.Manager {
    manager := log.NewManager()
    manager.SetMinLevel(handler.INFO)
    
    // File handler cho general logs
    generalHandler, _ := handler.NewFileHandler("/var/log/app.log", 100*1024*1024)
    manager.AddHandler("general", generalHandler)
    
    // File handler riêng cho errors
    errorHandler, _ := handler.NewFileHandler("/var/log/errors.log", 50*1024*1024)
    
    // Tạo error-only manager (hoặc filter trong custom handler)
    // Trong thực tế, bạn có thể cần custom handler để filter theo level
    manager.AddHandler("errors", errorHandler)
    
    return manager
}
```

### Testing Environment
```go
func setupTestLogging() log.Manager {
    manager := log.NewManager()
    manager.SetMinLevel(handler.ERROR) // Chỉ log errors trong tests
    
    // Có thể dùng mock handler
    return manager
}
```

## Tích hợp với Dependency Injection

### Sử dụng ServiceProvider
```go
import (
    "go.fork.vn/providers/log"
    "go.fork.vn/di"
)

func main() {
    container := di.NewContainer()
    
    // Đăng ký ServiceProvider
    provider := &log.ServiceProvider{}
    provider.Register(container)
    provider.Boot(container)
    
    // Sử dụng từ container
    var manager log.Manager
    container.Call(func(m log.Manager) {
        manager = m
        manager.Info("Logging system initialized via DI")
    })
}
```

### Custom Configuration với DI
```go
type App struct {
    container *di.Container
    basePath  string
}

func (app *App) Container() *di.Container {
    return app.container
}

func (app *App) BasePath() string {
    return app.basePath
}

func main() {
    app := &App{
        container: di.NewContainer(),
        basePath:  "/app",
    }
    
    provider := &log.ServiceProvider{}
    provider.Register(app)
    provider.Boot(app)
}
```

## Patterns và Best Practices

### 1. Singleton Logger Pattern
```go
var globalLogger log.Manager

func InitLogger() {
    globalLogger = log.NewManager()
    // Setup handlers...
}

func GetLogger() log.Manager {
    return globalLogger
}

// Sử dụng
func someFunction() {
    logger := GetLogger()
    logger.Info("Function called")
}
```

### 2. Contextual Logging
```go
type Service struct {
    logger log.Manager
    name   string
}

func NewService(logger log.Manager, name string) *Service {
    return &Service{
        logger: logger,
        name:   name,
    }
}

func (s *Service) ProcessData(data string) error {
    s.logger.Info("[%s] Processing data: %s", s.name, data)
    
    if err := process(data); err != nil {
        s.logger.Error("[%s] Processing failed: %v", s.name, err)
        return err
    }
    
    s.logger.Info("[%s] Processing completed", s.name)
    return nil
}
```

### 3. Structured Logging Pattern
```go
type LogEntry struct {
    Service   string
    Operation string
    UserID    int
    Duration  time.Duration
    Error     error
}

func (entry LogEntry) Log(manager log.Manager) {
    if entry.Error != nil {
        manager.Error("Operation failed: service=%s, op=%s, user=%d, duration=%v, error=%v",
            entry.Service, entry.Operation, entry.UserID, entry.Duration, entry.Error)
    } else {
        manager.Info("Operation success: service=%s, op=%s, user=%d, duration=%v",
            entry.Service, entry.Operation, entry.UserID, entry.Duration)
    }
}

// Sử dụng
func processUserData(manager log.Manager, userID int, data string) error {
    start := time.Now()
    
    err := doProcessing(data)
    
    LogEntry{
        Service:   "UserService",
        Operation: "ProcessData",
        UserID:    userID,
        Duration:  time.Since(start),
        Error:     err,
    }.Log(manager)
    
    return err
}
```

### 4. HTTP Request Logging
```go
func LoggingMiddleware(manager log.Manager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Log request
            manager.Info("HTTP Request: method=%s, path=%s, remote=%s", 
                r.Method, r.URL.Path, r.RemoteAddr)
            
            // Wrap ResponseWriter to capture status
            wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
            
            next.ServeHTTP(wrapped, r)
            
            // Log response
            duration := time.Since(start)
            manager.Info("HTTP Response: method=%s, path=%s, status=%d, duration=%v",
                r.Method, r.URL.Path, wrapped.statusCode, duration)
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### 5. Background Worker Logging
```go
func Worker(manager log.Manager, jobQueue <-chan Job) {
    for job := range jobQueue {
        manager.Debug("Worker received job: id=%s, type=%s", job.ID, job.Type)
        
        if err := processJob(job); err != nil {
            manager.Error("Job failed: id=%s, type=%s, error=%v", job.ID, job.Type, err)
            // Có thể gửi job vào retry queue
        } else {
            manager.Info("Job completed: id=%s, type=%s", job.ID, job.Type)
        }
    }
}
```

## Performance Tips

### 1. Conditional Expensive Logging
```go
// Tránh tính toán expensive khi không cần thiết
if manager.GetMinLevel() <= handler.DEBUG {
    expensiveData := generateExpensiveDebugInfo()
    manager.Debug("Debug info: %+v", expensiveData)
}
```

### 2. Batch Logging cho High Volume
```go
type BatchLogger struct {
    manager log.Manager
    buffer  []string
    mutex   sync.Mutex
    ticker  *time.Ticker
}

func NewBatchLogger(manager log.Manager) *BatchLogger {
    bl := &BatchLogger{
        manager: manager,
        buffer:  make([]string, 0, 100),
        ticker:  time.NewTicker(5 * time.Second),
    }
    
    go bl.flush()
    return bl
}

func (bl *BatchLogger) Log(message string) {
    bl.mutex.Lock()
    bl.buffer = append(bl.buffer, message)
    if len(bl.buffer) >= 100 {
        bl.flushBuffer()
    }
    bl.mutex.Unlock()
}

func (bl *BatchLogger) flush() {
    for range bl.ticker.C {
        bl.mutex.Lock()
        bl.flushBuffer()
        bl.mutex.Unlock()
    }
}

func (bl *BatchLogger) flushBuffer() {
    if len(bl.buffer) == 0 {
        return
    }
    
    batch := strings.Join(bl.buffer, "\n")
    bl.manager.Info("Batch log:\n%s", batch)
    bl.buffer = bl.buffer[:0]
}
```

## Troubleshooting

### 1. Logs không xuất hiện
```go
// Kiểm tra log level
manager.SetMinLevel(handler.DEBUG) // Temporarily lower level

// Kiểm tra handlers
if handler, exists := manager.GetHandler("console"); exists {
    // Handler exists
} else {
    // Handler not found
}
```

### 2. File permission errors
```go
fileHandler, err := handler.NewFileHandler("/var/log/app.log", 1024*1024)
if err != nil {
    // Fallback to temp directory
    tempFile := filepath.Join(os.TempDir(), "app.log")
    fileHandler, err = handler.NewFileHandler(tempFile, 1024*1024)
}
```

### 3. Memory leaks
```go
// Đảm bảo cleanup
defer manager.Close()

// Khi thay thế handlers
manager.RemoveHandler("old") // Tự động close old handler
manager.AddHandler("new", newHandler)
```

### 4. Concurrent access issues
```go
// Log manager là thread-safe, nhưng đảm bảo proper shutdown
var manager log.Manager
var once sync.Once

func GetManager() log.Manager {
    once.Do(func() {
        manager = log.NewManager()
        // Setup handlers...
    })
    return manager
}
```
