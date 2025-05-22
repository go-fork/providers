// Package formatter cung cấp các định dạng (formatter) khác nhau để nạp cấu hình từ nhiều nguồn.
//
// Package này chứa các formatter cho các nguồn cấu hình phổ biến như YAML, JSON, và
// biến môi trường (environment variables). Các formatter thực thi interface Formatter,
// cho phép hệ thống quản lý cấu hình nạp dữ liệu từ nhiều nguồn khác nhau một cách thống nhất.
package formatter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-fork/providers/config/utils"
	"gopkg.in/yaml.v3"
)

// YamlFormatter là một Formatter cho phép nạp cấu hình từ file YAML.
//
// YamlFormatter đọc nội dung từ file YAML, chuyển đổi thành cấu trúc dữ liệu map,
// và làm phẳng các key phân cấp thành dot notation. Ví dụ, một cấu trúc YAML như:
//
//	database:
//	  host: localhost
//	  port: 3306
//
// sẽ được chuyển đổi thành map với các key: "database.host": "localhost" và "database.port": 3306.
//
// YamlFormatter xử lý các kiểu dữ liệu scalar (string, number, boolean), mảng và map lồng nhau.
// Formatter này phù hợp cho việc quản lý cấu hình ứng dụng trong file cấu hình chính.
type YamlFormatter struct {
	path string // Đường dẫn đến file YAML chứa dữ liệu cấu hình
}

// NewYamlFormatter khởi tạo và trả về một YamlFormatter mới.
//
// Hàm này tạo một instance mới của YamlFormatter với đường dẫn file được chỉ định.
// Đường dẫn có thể là tương đối hoặc tuyệt đối, và file phải tồn tại và có quyền đọc.
//
// Params:
//   - path: string - Đường dẫn đến file YAML cần nạp cấu hình.
//
// Returns:
//   - *YamlFormatter: Con trỏ đến đối tượng YamlFormatter mới được khởi tạo.
//
// Examples:
//
//	formatter := NewYamlFormatter("config/app.yaml")
//	config, err := formatter.Load()
func NewYamlFormatter(path string) *YamlFormatter {
	return &YamlFormatter{
		path: path,
	}
}

// Load nạp và chuyển đổi cấu hình từ file YAML thành map các giá trị.
//
// Phương thức này thực hiện các bước xử lý sau:
//  1. Kiểm tra sự tồn tại của file YAML
//  2. Đọc nội dung của file vào bộ nhớ
//  3. Parse nội dung YAML thành cấu trúc dữ liệu Go
//  4. Làm phẳng cấu trúc phân cấp thành map với key dạng dot notation
//
// Quá trình làm phẳng sẽ chuyển đổi cấu trúc lồng nhau thành các key phẳng.
// Ví dụ, cấu trúc YAML như sau:
//
//	app:
//	  name: "MyApp"
//	  settings:
//	    debug: true
//
// sẽ được chuyển đổi thành:
//
//	"app.name": "MyApp"
//	"app.settings.debug": true
//
// Params: Không yêu cầu tham số đầu vào.
//
// Returns:
//   - map[string]interface{}: Map chứa các cặp key-value đã được làm phẳng từ file YAML.
//   - error: Lỗi nếu xảy ra vấn đề trong quá trình đọc hoặc parse file.
//
// Exceptions:
//   - "config file not found: [path]": File YAML không tồn tại tại đường dẫn đã chỉ định
//   - "failed to read config file: [error]": Không thể đọc nội dung file (lỗi quyền truy cập, v.v.)
//   - "failed to parse YAML config: [error]": Nội dung file không phải định dạng YAML hợp lệ
func (p *YamlFormatter) Load() (map[string]interface{}, error) {
	// Kiểm tra file tồn tại
	if _, err := os.Stat(p.path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", p.path)
	}

	// Đọc nội dung file vào bộ nhớ
	data, err := os.ReadFile(p.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Xử lý trường hợp file rỗng
	if len(data) == 0 {
		return make(map[string]interface{}), nil
	}

	// Parse dữ liệu YAML thành map[string]interface{}
	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %v", err)
	}

	// Chuyển đổi cấu trúc phân cấp thành map phẳng với key dạng dot notation
	flattenMap := make(map[string]interface{})
	utils.FlattenMapRecursive(flattenMap, result, "")

	return flattenMap, nil
}

// Name trả về tên định danh của formatter cho mục đích ghi log và debug.
//
// Phương thức này tạo một tên định danh duy nhất cho formatter, bao gồm loại formatter
// và tên file cấu hình. Tên này được sử dụng cho mục đích ghi log, debugging và theo dõi
// nguồn của các giá trị cấu hình.
//
// Format của tên trả về là "yaml:<tên-file>", trong đó <tên-file> là tên file
// (không bao gồm đường dẫn thư mục) được cung cấp khi khởi tạo formatter.
//
// Params: Không yêu cầu tham số đầu vào.
//
// Returns:
//   - string: Tên định danh của formatter theo format "yaml:<tên-file>".
//
// Examples:
//   - Nếu path = "/etc/app/config.yaml" thì Name() trả về "yaml:config.yaml"
//   - Nếu path = "./settings.yml" thì Name() trả về "yaml:settings.yml"
func (p *YamlFormatter) Name() string {
	return "yaml:" + filepath.Base(p.path)
}

// LoadFromDirectory nạp cấu hình từ tất cả các file YAML trong một thư mục.
//
// Hàm này quét qua tất cả các file có phần mở rộng .yaml hoặc .yml trong thư mục được chỉ định,
// nạp cấu hình từ mỗi file, và hợp nhất (merge) chúng vào một map kết quả duy nhất.
// Nếu có các key trùng nhau ở các file khác nhau, giá trị từ file được đọc sau sẽ ghi đè
// lên giá trị từ file được đọc trước đó.
//
// Quá trình xử lý:
//  1. Kiểm tra đường dẫn thư mục hợp lệ
//  2. Quét tất cả file có phần mở rộng .yaml hoặc .yml trong thư mục
//  3. Nạp từng file bằng YamlFormatter
//  4. Merge các giá trị cấu hình từ tất cả file vào một map kết quả
//
// Params:
//   - directory: string - Đường dẫn đến thư mục chứa các file YAML cần nạp.
//
// Returns:
//   - map[string]interface{}: Map hợp nhất chứa các giá trị cấu hình từ tất cả file YAML.
//   - error: Lỗi nếu có vấn đề trong quá trình đọc thư mục hoặc nạp file.
//
// Exceptions:
//   - "empty directory path": Đường dẫn thư mục rỗng
//   - "config directory not found: [path]": Thư mục không tồn tại
//   - "path is not a directory: [path]": Đường dẫn không phải là thư mục
//   - "failed to read config directory: [error]": Lỗi khi đọc thư mục
//   - "failed to load YAML file [path]: [error]": Lỗi khi nạp một file YAML cụ thể
//
// Examples:
//
//	configs, err := LoadFromDirectory("./configs")
//	if err != nil {
//	    log.Fatalf("Could not load config files: %v", err)
//	}
func LoadFromDirectory(directory string) (map[string]interface{}, error) {
	// Kiểm tra đường dẫn thư mục không rỗng
	if directory == "" {
		return nil, fmt.Errorf("empty directory path")
	}

	// Khởi tạo map kết quả
	result := make(map[string]interface{})

	// Kiểm tra thư mục tồn tại
	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config directory not found: %s", directory)
	}

	// Kiểm tra đường dẫn là thư mục
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", directory)
	}

	// Đọc danh sách các file trong thư mục
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %v", err)
	}

	// Lặp qua từng mục trong thư mục
	for _, file := range files {
		// Bỏ qua các thư mục con
		if file.IsDir() {
			continue
		}

		// Chỉ xử lý các file có phần mở rộng .yaml hoặc .yml
		ext := filepath.Ext(file.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		// Xây dựng đường dẫn đầy đủ đến file
		filePath := filepath.Join(directory, file.Name())

		// Tạo formatter và nạp file
		formatter := NewYamlFormatter(filePath)
		values, err := formatter.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to load YAML file %s: %w", filePath, err)
		}

		// Merge các giá trị từ file vào map kết quả
		for k, v := range values {
			result[k] = v
		}
	}

	return result, nil
}
