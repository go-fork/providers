package handler

// StackHandler triển khai một handler log tổng hợp chuyển tiếp các bản ghi log đến nhiều handlers.
//
// Tính năng:
//   - Quản lý nhiều handlers như một đơn vị
//   - Chuyển tiếp tuần tự đến tất cả các handlers con
//   - Xử lý lỗi tập trung
//   - Thêm handler động
type StackHandler struct {
	handlers []Handler // Slice chứa các handlers con
}

// NewStackHandler tạo một stack handler mới với các handlers con được chỉ định.
//
// Tham số:
//   - handlers: ...Handler - không hoặc nhiều handlers để bao gồm trong stack
//
// Trả về:
//   - *StackHandler: một stack handler đã được cấu hình
//
// Ví dụ:
//
//	// Tạo một stack với console và file handlers
//	consoleHandler := handler.NewConsoleHandler(true)
//	fileHandler, _ := handler.NewFileHandler("app.log", 10*1024*1024)
//	stackHandler := handler.NewStackHandler(consoleHandler, fileHandler)
func NewStackHandler(handlers ...Handler) *StackHandler {
	return &StackHandler{
		handlers: handlers,
	}
}

// Log chuyển tiếp một log entry đến tất cả các handlers trong stack.
//
// Phương thức này gọi phương thức Log của mỗi handler con theo thứ tự.
// Nếu bất kỳ handler nào trả về lỗi, lỗi đầu tiên sẽ được trả về,
// nhưng tất cả các handlers sẽ vẫn được gọi.
//
// Tham số:
//   - level: Level - cấp độ nghiêm trọng của log entry
//   - message: string - thông điệp log
//   - args: ...interface{} - các tham số định dạng tùy chọn
//
// Trả về:
//   - error: lỗi đầu tiên gặp phải, hoặc nil nếu tất cả handlers thành công
func (a *StackHandler) Log(level Level, message string, args ...interface{}) error {
	var firstErr error
	for _, handler := range a.handlers {
		if err := handler.Log(level, message, args...); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Close đóng đúng cách tất cả các handlers trong stack.
//
// Phương thức này gọi phương thức Close của mỗi handler con theo thứ tự.
// Nếu bất kỳ handler nào trả về lỗi, lỗi đầu tiên sẽ được trả về,
// nhưng tất cả các handlers sẽ vẫn được đóng.
//
// Trả về:
//   - error: lỗi đầu tiên gặp phải, hoặc nil nếu tất cả handlers đóng thành công
func (a *StackHandler) Close() error {
	var firstErr error
	for _, handler := range a.handlers {
		if err := handler.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// AddHandler thêm một handler mới vào stack.
//
// Tham số:
//   - handler: Handler - handler để thêm vào stack
//
// Ví dụ:
//
//	// Thêm một network handler mới vào một stack đã tồn tại
//	networkHandler := NewNetworkHandler("logs.example.com:514")
//	stackHandler.AddHandler(networkHandler)
func (a *StackHandler) AddHandler(handler Handler) {
	a.handlers = append(a.handlers, handler)
}
