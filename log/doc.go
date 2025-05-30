// Package log cung cấp hệ thống logging linh hoạt, dễ mở rộng và thread-safe cho ứng dụng Go.
//
// # Tổng quan
//
// Package này triển khai hệ thống logging với nhiều cấp độ nghiêm trọng, các output
// handler khác nhau (console, file, v.v.), và interface quản lý tập trung. Nó được thiết kế
// để thread-safe và quản lý tài nguyên hiệu quả trong các ứng dụng concurrent.
//
// # Tính năng
//
//   - Lọc theo cấp độ log (Debug, Info, Warning, Error, Fatal)
//   - Nhiều output handler hoạt động đồng thời
//   - Hoạt động thread-safe
//   - Hỗ trợ chuỗi định dạng
//   - Khả năng mở rộng handler tùy chỉnh
//   - Output console có màu
//   - Tự động xoay vòng file log
//   - Hỗ trợ dependency injection
//
// # Sử dụng cơ bản
//
//	// Tạo một log manager mới
//	manager := log.NewManager()
//
//	// Thêm console handler có màu
//	consoleHandler := handler.NewConsoleHandler(true)
//	manager.AddHandler("console", consoleHandler)
//
//	// Thêm file handler với kích thước tối đa 10MB
//	fileHandler, err := handler.NewFileHandler("app.log", 10*1024*1024)
//	if err == nil {
//	    manager.AddHandler("file", fileHandler)
//	}
//
//	// Ghi log ở các cấp độ khác nhau
//	manager.Debug("Thông tin debug")
//	manager.Info("Ứng dụng đã khởi động với cấu hình: %v", config)
//	manager.Warning("Sử dụng tài nguyên cao: %d%%", usagePercent)
//	manager.Error("Xử lý yêu cầu thất bại: %v", err)
//	manager.Fatal("Mất kết nối database")
//
//	// Thiết lập cấp độ log tối thiểu
//	manager.SetMinLevel(handler.InfoLevel) // Bỏ qua Debug logs
//
//	// Đóng tất cả các handler đúng cách khi tắt
//	defer manager.Close()
//
// # Sử dụng nâng cao
//
// Đối với các yêu cầu logging phức tạp hơn, package hỗ trợ handler tùy chỉnh thông qua
// interface handler.Handler, cho phép các hành vi logging cụ thể như
// network logging, lưu trữ database, hoặc tích hợp với hệ thống monitoring bên ngoài.
//
// Truy xuất handler đã đăng ký để cấu hình hoặc kiểm tra:
//
//	// Lấy handler theo tên để cấu hình
//	if consoleHandler := manager.GetHandler("console"); consoleHandler != nil {
//	    // Cấu hình thêm cho handler
//	    if typedHandler, ok := consoleHandler.(*handler.ConsoleHandler); ok {
//	        // Thực hiện cấu hình đặc thù cho ConsoleHandler
//	    }
//	}
//
// # Xem thêm
//
// Interface Manager và triển khai DefaultManager cho các thao tác logging chính.
// Package handler để tùy chỉnh xử lý output.
//
// # Tương thích DI
//
// Module này tương thích đầy đủ với go.fork.vn/di từ phiên bản v0.0.5 trở lên,
// cài đặt đầy đủ interface ServiceProvider với các phương thức Register, Boot, Requires và Providers.
package log
