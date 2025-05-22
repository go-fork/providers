// Package formatter cung cấp các implementation cho việc nạp cấu hình từ các nguồn khác nhau.
package formatter

import (
	"os"
	"strconv"
	"strings"
)

// EnvFormatter là một Formatter cho phép nạp cấu hình từ biến môi trường (environment variables).
//
// EnvFormatter quét tất cả các biến môi trường của hệ thống, lọc theo prefix được cung cấp,
// và chuyển đổi chúng thành map cấu hình. Quá trình này bao gồm:
//   - Lọc biến môi trường bắt đầu bằng prefix đã chỉ định
//   - Loại bỏ prefix và gạch dưới đầu tiên sau prefix (nếu có)
//   - Chuyển key sang chữ thường (lowercase) để chuẩn hóa
//   - Chuyển đổi dấu gạch dưới (_) trong key thành dấu chấm (.) để hỗ trợ dot notation
//   - Tự động chuyển đổi giá trị thành kiểu dữ liệu phù hợp (bool, int, float, string)
//
// Ví dụ, với prefix là "APP_":
//   - Biến môi trường "APP_SERVER_HOST=localhost" sẽ trở thành key "server.host" với giá trị "localhost"
//   - Biến môi trường "APP_DATABASE_PORT=5432" sẽ trở thành key "database.port" với giá trị số nguyên 5432
//   - Biến môi trường "APP_FEATURES_DEBUG=true" sẽ trở thành key "features.debug" với giá trị boolean true
//
// EnvFormatter đặc biệt hữu ích cho việc ghi đè cấu hình trong các môi trường container
// hoặc cloud, nơi biến môi trường là cách chính để cấu hình ứng dụng.
type EnvFormatter struct {
	prefix string // Tiền tố để lọc các biến môi trường liên quan đến ứng dụng
}

// NewEnvFormatter khởi tạo và trả về một EnvFormatter mới.
//
// Hàm này tạo một formatter để nạp cấu hình từ biến môi trường với tiền tố được chỉ định.
// Việc sử dụng tiền tố giúp phân biệt các biến môi trường của ứng dụng với các biến
// môi trường hệ thống khác, tránh xung đột và đảm bảo an toàn.
//
// Params:
//   - prefix: string - Tiền tố cho các biến môi trường liên quan đến ứng dụng.
//     Ví dụ: "APP_", "MYAPP_". Nếu để trống, formatter sẽ không nạp bất kỳ biến nào
//     (vì lý do bảo mật).
//
// Returns:
//   - *EnvFormatter: Con trỏ đến đối tượng EnvFormatter mới được khởi tạo.
//
// Examples:
//
//	formatter := NewEnvFormatter("APP_")
//	config, err := formatter.Load()
//	// Với biến môi trường APP_DATABASE_HOST=localhost
//	// config sẽ chứa {"database.host": "localhost"}
func NewEnvFormatter(prefix string) *EnvFormatter {
	return &EnvFormatter{
		prefix: prefix,
	}
}

// Load nạp cấu hình từ các biến môi trường và chuyển đổi thành map cấu hình.
//
// Phương thức này thực hiện các bước xử lý sau:
//  1. Kiểm tra tính hợp lệ của prefix (từ chối prefix rỗng vì lý do bảo mật)
//  2. Quét tất cả các biến môi trường hiện có trong hệ thống
//  3. Lọc các biến môi trường có tiền tố phù hợp
//  4. Xử lý key (loại bỏ prefix, chuyển _ thành ., chuyển sang lowercase)
//  5. Chuyển đổi giá trị sang kiểu dữ liệu phù hợp (bool, int, float, string)
//  6. Xây dựng map kết quả
//
// Quá trình chuyển đổi kiểu dữ liệu tự động sẽ nhận diện:
//   - Boolean: Các giá trị "true", "yes", "1", "on" được chuyển thành true,
//     các giá trị "false", "no", "0", "off" được chuyển thành false
//   - Số nguyên: Chuỗi chỉ chứa các chữ số được chuyển thành kiểu int
//   - Số thực: Chuỗi chứa dấu thập phân được chuyển thành kiểu float64
//   - Chuỗi: Mọi giá trị khác được giữ nguyên dưới dạng chuỗi
//
// Params: Không yêu cầu tham số đầu vào.
//
// Returns:
//   - map[string]interface{}: Map chứa các cặp key-value được nạp từ biến môi trường.
//   - error: Luôn là nil vì việc nạp từ biến môi trường không phát sinh lỗi.
//
// Examples:
//
//	Với các biến môi trường:
//	  APP_SERVER_HOST=localhost
//	  APP_SERVER_PORT=8080
//	  APP_DEBUG=true
//
//	formatter := NewEnvFormatter("APP_")
//	config, _ := formatter.Load()
//	// config sẽ chứa:
//	// {
//	//   "server.host": "localhost",
//	//   "server.port": 8080,
//	//   "debug": true
//	// }
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

// processEnvValue chuyển đổi giá trị chuỗi từ biến môi trường thành kiểu dữ liệu phù hợp.
//
// Hàm này phân tích giá trị chuỗi và cố gắng chuyển đổi nó thành kiểu dữ liệu
// phù hợp nhất dựa trên nội dung của chuỗi. Thứ tự ưu tiên chuyển đổi là:
//  1. Boolean: Nhận diện các giá trị chuỗi biểu thị true/false
//  2. Integer: Thử chuyển đổi thành số nguyên
//  3. Float: Thử chuyển đổi thành số thực
//  4. String: Giữ nguyên dạng chuỗi nếu không khớp với các kiểu trên
//
// Params:
//   - value: string - Giá trị chuỗi cần chuyển đổi kiểu dữ liệu
//
// Returns:
//   - interface{}: Giá trị sau khi chuyển đổi sang kiểu dữ liệu phù hợp (bool, int, float64, hoặc string)
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

// Name trả về tên định danh của formatter cho mục đích ghi log và debug.
//
// Phương thức này trả về một chuỗi định danh duy nhất cho formatter môi trường.
// Khác với JsonFormatter hay YamlFormatter, EnvFormatter không liên kết với một file
// cụ thể, nên Name() chỉ đơn giản trả về "env" để xác định nguồn cấu hình
// là từ biến môi trường.
//
// Tên này được sử dụng bởi Manager để xác định nguồn cấu hình trong quá trình gỡ lỗi,
// ghi log, và theo dõi quá trình nạp cấu hình.
//
// Params: Không yêu cầu tham số đầu vào.
//
// Returns:
//   - string: Tên định danh cố định "env".
func (p *EnvFormatter) Name() string {
	return "env"
}
