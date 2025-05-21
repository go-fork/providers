package handler

// Level đại diện cho cấp độ nghiêm trọng của một log entry.
//
// Các cấp độ được sắp xếp từ thấp đến cao nhất, cho phép lọc dựa trên
// ngưỡng cấp độ tối thiểu. Các cấp độ tiêu chuẩn là Debug, Info, Warning, Error,
// và Fatal.
type Level int

const (
	// DebugLevel đại diện cho thông tin chẩn đoán chi tiết (mức nghiêm trọng thấp nhất).
	DebugLevel Level = iota

	// InfoLevel đại diện cho thông tin hoạt động chung.
	InfoLevel

	// WarningLevel đại diện cho vấn đề tiềm ẩn hoặc điều kiện không mong đợi.
	WarningLevel

	// ErrorLevel đại diện cho lỗi ảnh hưởng đến hoạt động bình thường.
	ErrorLevel

	// FatalLevel đại diện cho lỗi nghiêm trọng cần chú ý ngay lập tức (mức nghiêm trọng cao nhất).
	FatalLevel
)

// String trả về biểu diễn chuỗi của một cấp độ log.
//
// Method này hữu ích khi cần gồm tên cấp độ trong output log
// và để chuyển đổi chuỗi trong giao diện người dùng.
//
// Trả về:
//   - string: tên dễ đọc của cấp độ log
//
// Ví dụ:
//
//	fmt.Printf("Cấp độ log hiện tại: %s\n", currentLevel.String())
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARNING"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Handler là interface mà tất cả các log handler phải triển khai.
//
// Handler chịu trách nhiệm xử lý các log entry và ghi chúng vào
// một đích đến output (console, file, network, v.v.). Các handlers nên
// được triển khai theo cách thread-safe.
type Handler interface {
	// Log xử lý một log entry với cấp độ và thông điệp được chỉ định.
	//
	// Tham số:
	//   - level: Level - cấp độ nghiêm trọng của log entry
	//   - message: string - thông điệp log đã được định dạng
	//   - args: ...interface{} - tham số định dạng tùy chọn
	//
	// Trả về:
	//   - error: một lỗi nếu log entry không thể được xử lý
	Log(level Level, message string, args ...interface{}) error

	// Close giải phóng bất kỳ tài nguyên nào được sử dụng bởi handler.
	//
	// Method này nên được gọi khi không cần handler nữa
	// để đảm bảo dọn dẹp đúng cách (đóng file, kết nối, v.v.).
	//
	// Trả về:
	//   - error: một lỗi nếu dọn dẹp thất bại
	Close() error
}
