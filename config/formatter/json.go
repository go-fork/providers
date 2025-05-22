// Package formatter cung cấp các implementation cho việc nạp cấu hình từ các nguồn khác nhau.
package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-fork/providers/config/utils"
)

// JsonFormatter là một Formatter cho phép nạp cấu hình từ file JSON.
//
// JsonFormatter đọc dữ liệu từ file JSON, phân tích cú pháp (parse) thành cấu trúc dữ liệu map,
// và chuyển đổi các key phân cấp thành dạng dot notation. Ví dụ, một cấu trúc JSON như:
//
//	{
//	  "database": {
//	    "host": "localhost",
//	    "port": 3306
//	  }
//	}
//
// sẽ được chuyển đổi thành map với các key: "database.host": "localhost" và "database.port": 3306.
//
// JsonFormatter hỗ trợ tất cả các kiểu dữ liệu JSON hợp lệ: objects, arrays, strings,
// numbers, booleans và null values. Kiểu dữ liệu này phù hợp cho cấu hình ứng dụng với
// các cấu trúc dữ liệu phức tạp và hỗ trợ tốt cho dữ liệu số.
type JsonFormatter struct {
	path string // Đường dẫn đến file JSON chứa dữ liệu cấu hình
}

// NewJsonFormatter khởi tạo và trả về một JsonFormatter mới.
//
// Hàm này tạo một instance mới của JsonFormatter với đường dẫn file JSON được chỉ định.
// Đường dẫn có thể là tương đối hoặc tuyệt đối. File sẽ được kiểm tra tồn tại và
// khả năng đọc khi phương thức Load được gọi.
//
// Params:
//   - path: string - Đường dẫn đến file JSON cần nạp cấu hình.
//
// Returns:
//   - *JsonFormatter: Con trỏ đến đối tượng JsonFormatter mới được khởi tạo.
//
// Examples:
//
//	formatter := NewJsonFormatter("config/app.json")
//	config, err := formatter.Load()
func NewJsonFormatter(path string) *JsonFormatter {
	return &JsonFormatter{
		path: path,
	}
}

// Load nạp và chuyển đổi cấu hình từ file JSON thành map các giá trị.
//
// Phương thức này thực hiện các bước xử lý sau:
//  1. Kiểm tra tính hợp lệ của đường dẫn file
//  2. Đọc nội dung của file JSON vào bộ nhớ
//  3. Parse nội dung JSON thành cấu trúc dữ liệu Go
//  4. Làm phẳng cấu trúc phân cấp thành map với key dạng dot notation
//
// Phương thức này có các xử lý đặc biệt cho các trường hợp:
//   - File JSON rỗng hoặc chỉ chứa object rỗng ({}) -> trả về map rỗng
//   - JSON không phải là object (như array hoặc giá trị nguyên thủy) -> trả về map rỗng
//
// Quá trình làm phẳng sẽ chuyển đổi cấu trúc lồng nhau thành các key phẳng.
// Ví dụ, cấu trúc JSON như sau:
//
//	{
//	  "app": {
//	    "name": "MyApp",
//	    "version": 1.2,
//	    "features": ["auth", "api"]
//	  }
//	}
//
// sẽ được chuyển đổi thành:
//
//	"app.name": "MyApp"
//	"app.version": 1.2
//	"app.features": ["auth", "api"]
//
// Params: Không yêu cầu tham số đầu vào.
//
// Returns:
//   - map[string]interface{}: Map chứa các cặp key-value đã được làm phẳng từ file JSON.
//   - error: Lỗi nếu xảy ra vấn đề trong quá trình đọc hoặc parse file.
//
// Exceptions:
//   - "empty file path": Đường dẫn file rỗng
//   - "config file not found: [path]": File JSON không tồn tại tại đường dẫn đã chỉ định
//   - "failed to read config file: [error]": Không thể đọc nội dung file (lỗi quyền truy cập, v.v.)
//   - "failed to parse JSON config at [path]: [error]": Nội dung file không phải định dạng JSON hợp lệ
func (p *JsonFormatter) Load() (map[string]interface{}, error) {
	// Kiểm tra đường dẫn file không rỗng
	if p.path == "" {
		return nil, fmt.Errorf("empty file path")
	}

	// Đọc nội dung file vào bộ nhớ
	data, err := os.ReadFile(p.path)
	if err != nil {
		// Xử lý trường hợp file không tồn tại
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", p.path)
		}
		// Xử lý các lỗi đọc file khác
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Xử lý trường hợp file rỗng hoặc chỉ chứa object rỗng
	if len(data) == 0 || (len(data) == 2 && string(data) == "{}") {
		return make(map[string]interface{}), nil
	}

	// Parse dữ liệu JSON thành map[string]interface{}
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		// Nếu không phải là object JSON, thử parse thành kiểu dữ liệu bất kỳ
		// (có thể là array, string, number, boolean, hoặc null)
		var anyValue interface{}
		if errAny := json.Unmarshal(data, &anyValue); errAny == nil {
			// Trường hợp JSON hợp lệ nhưng không phải là object -> trả về map rỗng
			return make(map[string]interface{}), nil
		}
		// Trường hợp JSON không hợp lệ -> trả về lỗi
		return nil, fmt.Errorf("failed to parse JSON config at %s: %v", p.path, err)
	}

	// Chuyển đổi cấu trúc phân cấp thành map phẳng với key dạng dot notation
	flattenMap := make(map[string]interface{})
	utils.FlattenMapRecursive(flattenMap, result, "")
	return flattenMap, nil
}

// Name trả về tên định danh của formatter cho mục đích ghi log và debug.
//
// Phương thức này tạo một tên định danh duy nhất cho formatter, bao gồm loại formatter
// và tên file cấu hình. Tên này được sử dụng bởi Manager để xác định nguồn cấu hình
// trong quá trình gỡ lỗi, ghi log, và theo dõi quá trình nạp cấu hình.
//
// Format của tên trả về là "json:<tên-file>", trong đó <tên-file> là tên file
// (không bao gồm đường dẫn thư mục) được cung cấp khi khởi tạo formatter.
//
// Params: Không yêu cầu tham số đầu vào.
//
// Returns:
//   - string: Tên định danh của formatter theo format "json:<tên-file>".
//
// Examples:
//   - Nếu path = "/etc/app/config.json" thì Name() trả về "json:config.json"
//   - Nếu path = "./settings.json" thì Name() trả về "json:settings.json"
func (p *JsonFormatter) Name() string {
	return "json:" + filepath.Base(p.path)
}
