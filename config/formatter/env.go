package formatter

import (
	"os"
	"strconv"
	"strings"
)

// EnvFormatter cung cấp Formatter cho phép nạp cấu hình từ biến môi trường (environment variable).
//
// Các biến môi trường sẽ được lọc theo prefix, chuyển đổi key sang lowercase và thay _ thành . để hỗ trợ dot notation.
type EnvFormatter struct {
	prefix string // Prefix để lọc biến môi trường
}

// NewEnvFormatter khởi tạo một EnvFormatter mới.
// NewEnvFormatter nhận vào prefix để lọc biến môi trường liên quan.
func NewEnvFormatter(prefix string) *EnvFormatter {
	return &EnvFormatter{
		prefix: prefix,
	}
}

// Load tải cấu hình từ biến môi trường.
// Load trả về map[string]interface{} chứa các giá trị cấu hình, hoặc error nếu có lỗi khi nạp.
func (p *EnvFormatter) Load() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	prefix := p.prefix

	if prefix == "" {
		// Nếu không có prefix, lấy tất cả biến môi trường có thể gây ra vấn đề bảo mật
		// hoặc xung đột với các ứng dụng khác. Tốt nhất là luôn yêu cầu prefix.
		return result, nil
	}

	prefixLen := len(prefix)

	// Lấy tất cả biến môi trường
	for _, env := range os.Environ() {
		// Tách key và value
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]

		// Chỉ lấy các biến có prefix đúng
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		// Bỏ qua biến có value rỗng (để hỗ trợ test isolation tuyệt đối)
		if value == "" {
			continue
		}

		// Loại bỏ prefix và dấu gạch dưới đầu tiên sau prefix (nếu có)
		key = key[prefixLen:]
		if strings.HasPrefix(key, "_") {
			key = key[1:]
		}

		// Chuyển key sang lowercase và thay _ thành . để hỗ trợ dot notation
		key = strings.ToLower(key)
		key = strings.ReplaceAll(key, "_", ".")

		// Chuyển đổi value sang kiểu dữ liệu phù hợp
		processedValue := processEnvValue(value)

		// Lưu vào map kết quả
		result[key] = processedValue
	}

	return result, nil
}

// processEnvValue cố gắng chuyển đổi giá trị env sang kiểu dữ liệu phù hợp (bool, int, float)
func processEnvValue(value string) interface{} {
	switch strings.ToLower(value) {
	case "true", "yes", "1", "on":
		return true
	case "false", "no", "0", "off":
		return false
	}
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}
	return value
}

// Name trả về tên định danh của EnvFormatter.
// Name trả về string "env".
func (p *EnvFormatter) Name() string {
	return "env"
}
