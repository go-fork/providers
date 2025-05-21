package handler

import (
	"fmt"
	"os"
	"time"
)

// ConsoleHandler triển khai một log handler ghi ra console (stdout/stderr).
//
// Tính năng:
//   - Output có mã màu dựa trên cấp độ log
//   - Tự động định tuyến errors ra stderr
//   - Định dạng timestamp chuẩn
//   - Tùy chọn zero-configuration
type ConsoleHandler struct {
	colored bool // Có sử dụng mã màu ANSI hay không
}

// NewConsoleHandler tạo một console handler mới.
//
// Tham số:
//   - colored: bool - có sử dụng mã màu ANSI trong output hay không
//
// Trả về:
//   - *ConsoleHandler: một console handler đã được cấu hình
//
// Ví dụ:
//
//	// Tạo một console handler có màu
//	handler := handler.NewConsoleHandler(true)
//
//	// Tạo một console handler plain-text
//	handler := handler.NewConsoleHandler(false)
func NewConsoleHandler(colored bool) *ConsoleHandler {
	return &ConsoleHandler{
		colored: colored,
	}
}

// Log ghi một log entry ra console.
//
// Method này định dạng log entry với timestamp và chỉ báo cấp độ,
// áp dụng màu nếu được bật, và ghi ra stdout hoặc stderr tùy thuộc
// vào mức độ nghiêm trọng.
//
// Tham số:
//   - level: Level - cấp độ nghiêm trọng của log entry
//   - message: string - thông điệp log
//   - args: ...interface{} - tham số định dạng tùy chọn (hiện không sử dụng)
//
// Trả về:
//   - error: một lỗi nếu ghi ra console thất bại
func (a *ConsoleHandler) Log(level Level, message string, args ...interface{}) error {
	// Định dạng với timestamp và cấp độ
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	formattedMessage := fmt.Sprintf("%s [%s] %s\n", timestamp, level.String(), message)

	// Áp dụng mã màu ANSI nếu được bật
	if a.colored {
		formattedMessage = a.colorize(level, formattedMessage)
	}

	// Ghi ra stderr cho log Error và Fatal
	if level >= ErrorLevel {
		_, err := fmt.Fprint(os.Stderr, formattedMessage)
		return err
	}

	// Ghi ra stdout cho các cấp độ khác
	_, err := fmt.Fprint(os.Stdout, formattedMessage)
	return err
}

// Close giải phóng tài nguyên được sử dụng bởi console handler.
//
// Đối với console handler, đây là thao tác rỗng vì I/O console không yêu cầu
// dọn dẹp rõ ràng.
//
// Trả về:
//   - nil: console handler không giữ tài nguyên cần dọn dẹp
func (a *ConsoleHandler) Close() error {
	// Không cần dọn dẹp cho console I/O
	return nil
}

// colorize áp dụng mã màu ANSI vào thông điệp dựa trên cấp độ log.
//
// Tham số:
//   - level: Level - cấp độ nghiêm trọng xác định màu sắc
//   - message: string - thông điệp log đã định dạng
//
// Trả về:
//   - string: thông điệp với mã màu ANSI đã áp dụng
func (a *ConsoleHandler) colorize(level Level, message string) string {
	// Chọn mã màu ANSI dựa trên cấp độ log
	var colorCode string

	switch level {
	case DebugLevel:
		colorCode = "\033[36m" // Cyan cho debug
	case InfoLevel:
		colorCode = "\033[32m" // Green cho info
	case WarningLevel:
		colorCode = "\033[33m" // Yellow cho warning
	case ErrorLevel:
		colorCode = "\033[31m" // Red cho error
	case FatalLevel:
		colorCode = "\033[35m" // Magenta cho fatal
	default:
		colorCode = "\033[0m" // Mặc định (reset)
	}

	// Áp dụng màu và đảm bảo reset ở cuối
	return fmt.Sprintf("%s%s\033[0m", colorCode, message)
}
