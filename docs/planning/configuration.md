# Phân tích triển khai Service Provider Configuration

## 1. Tổng quan

**Service Provider Configuration** là một hệ thống quản lý cấu hình trong Go, cho phép ứng dụng truy xuất và ánh xạ cấu hình từ nhiều định dạng (ENV, JSON, YAML) thông qua các **formatter**. Hệ thống sử dụng **flatten map** với **dot notation** để biểu diễn cấu hình phẳng, hỗ trợ các phương thức **Get**, **Set**, **Has**, **Unmarshal** để thao tác, và đảm bảo **thread-safe** trong môi trường đa luồng.

- **Mục tiêu nghiệp vụ**:
  - Cung cấp API type-safe, dễ sử dụng để truy xuất và ánh xạ cấu hình.
  - Hỗ trợ nhiều định dạng cấu hình mà không cần sửa mã nguồn.
  - Đảm bảo tính nhất quán khi merge cấu hình từ các formatter theo thứ tự ưu tiên (ENV > JSON > YAML).
  - Tối ưu hiệu năng, bảo trì, và khả năng mở rộng.
- **Giá trị nghiệp vụ**:
  - **Tính linh hoạt**: Hỗ trợ ENV, JSON, YAML; dễ mở rộng cho TOML, XML.
  - **Tính dễ sử dụng**: API đơn giản, type-safe, giảm lỗi runtime.
  - **Tính bảo trì**: Tuân theo **SOLID** principles (Single Responsibility, Open/Closed).
  - **Tính mở rộng**: Thêm formatter, validation, hot reload dễ dàng.
---

## 2. Phân tích yêu cầu (Business Analysis - BA)

### 2.1. Yêu cầu chức năng (Functional Requirements)

1. **Formatter (ENV, JSON, YAML)**:
   - **Định nghĩa**: Thành phần xử lý dữ liệu cấu hình, tích hợp ba tính năng:
     - **Load**: Tải dữ liệu từ nguồn (biến môi trường, file).
     - **Parse**: Phân tích cú pháp dữ liệu theo định dạng.
     - **Flatten Recursive**: Chuyển cấu trúc lồng nhau thành **flatten map** với **dot notation**.
   - **ENV Formatter**:
     - **Load**: Đọc biến môi trường từ `os.Environ()`.
     - **Parse**: Chuyển key từ `UPPER_CASE`/`CAMEL_CASE` sang `lower.case` (ví dụ: `DB_HOST` → `db.host`).
     - **Flatten Recursive**: Tạo `ConfigMap` từ key phẳng.
   - **JSON Formatter**:
     - **Load**: Đọc file JSON từ hệ thống file (`os.ReadFile`).
     - **Parse**: Sử dụng `encoding/json` để parse thành `map`/`slice`.
     - **Flatten Recursive**: Chuyển cấu trúc lồng nhau thành `ConfigMap`.
   - **YAML Formatter**:
     - **Load**: Đọc file YAML từ hệ thống file.
     - **Parse**: Sử dụng `gopkg.in/yaml.v3` để parse.
     - **Flatten Recursive**: Chuyển cấu trúc lồng nhau thành `ConfigMap`.
   - **Ghi chú**: Formatter đảm bảo `ConfigMap` chứa key dạng dot notation và giá trị với kiểu dữ liệu rõ ràng.

2. **Flatten Map với Dot Notation**:
   - **Định nghĩa**: Chuyển dữ liệu lồng nhau (map, slice, struct) thành `map[string]ConfigValue`.
   - **Dot Notation**: Biểu diễn key dạng `prefix.subkey` (ví dụ: `database.host`).
   - **Yêu cầu**:
     - Lưu giá trị gốc (slice, map) và các phần tử con.
     - Hỗ trợ tùy chọn: separator, skip nil, case sensitivity, handle empty key.
   - **Ghi chú**: Flatten recursive là thuật toán đệ quy để duyệt cấu trúc dữ liệu, tạo key phẳng.

3. **Phương thức thao tác cấu hình**:
   - **Get**: Truy xuất giá trị type-safe:
     ```go
     GetString(key string) (string, bool)
     GetInt(key string) (int, bool)
     GetBool(key string) (bool, bool)
     GetSlice(key string) ([]interface{}, bool)
     GetMap(key string) (map[string]interface{}, bool)
     ```
     - **Ghi chú**: Type-safe đảm bảo kiểu dữ liệu khớp, trả về `false` nếu key không tồn tại hoặc kiểu sai.
   - **Set**: Cập nhật/thêm giá trị vào `ConfigMap`:
     ```go
     Set(key string, value interface{}) error
     ```
     - **Ghi chú**: Kiểm tra key hợp lệ, xác định kiểu dữ liệu bằng `reflect`.
   - **Has**: Kiểm tra sự tồn tại của key:
     ```go
     Has(key string) bool
     ```
     - **Ghi chú**: Trả về `true` nếu key tồn tại, O(1) nhờ map phẳng.
   - **Unmarshal**: Ánh xạ `ConfigMap` vào Go struct:
     ```go
     Unmarshal(key string, target interface{}) error
     ```
     - **Hành vi**:
       - `key == ""`: Ánh xạ toàn bộ `ConfigMap` (root path).
       - `key != ""`: Ánh xạ dữ liệu tại `key` và các key con (ví dụ: `database` → `database.host`, `database.port`).
     - **Yêu cầu**:
       - Hỗ trợ tag `config` (ví dụ: `config:"database.host"`).
       - Hỗ trợ struct lồng nhau, slice, map.
       - Type conversion an toàn (string → int).
       - Báo lỗi nếu `key` không tồn tại (khi `key != ""`).

4. **Thứ tự ưu tiên**:
   - Mặc định: **ENV > JSON > YAML** (có thể cấu hình lại).
   - Giá trị từ formatter ưu tiên cao ghi đè giá trị từ formatter thấp hơn.
   - **Ghi chú**: Merge cấu hình đảm bảo tính nhất quán.

5. **Thread-safe**:
   - Sử dụng `sync.RWMutex`:
     - `Lock()`: Khi load, merge, set, unmarshal.
     - `RLock()`: Khi get, has.
   - **Ghi chú**: Phù hợp với ứng dụng đa luồng (web server, microservices).

6. **Error Handling**:
   - Báo lỗi chi tiết với `fmt.Errorf` và `%w`:
     - `ErrInvalidKey`: Key không hợp lệ (rỗng, ký tự đặc biệt).
     - `ErrKeyNotFound`: Key không tồn tại.
     - `ErrTypeMismatch`: Kiểu dữ liệu không khớp.
     - `ErrInvalidTarget`: Target không phải con trỏ hoặc struct.
     - `ErrParseFailed`: Parse dữ liệu thất bại.
   - **Ghi chú**: Wrap lỗi để hỗ trợ debug.

7. **Khả năng mở rộng**:
   - Thêm formatter mới (TOML, XML) thông qua interface `Formatter`.
   - Hỗ trợ validation (schema, key bắt buộc).
   - Hỗ trợ hot reload (theo dõi file với `fsnotify`).
   - Hỗ trợ giá trị mặc định.

### 2.2. Yêu cầu phi chức năng (Non-functional Requirements)

1. **Hiệu năng**:
   - **Load/Parse/Flatten**: O(n), n là số phần tử trong dữ liệu.
   - **Get/Set/Has**: O(1) nhờ map phẳng.
   - **Unmarshal**: O(k + f), k là số key trong `ConfigMap`, f là số field trong struct.
   - Tối ưu bộ nhớ: Bỏ qua `nil`, cache ánh xạ tag-key.
2. **Bảo trì**:
   - Tuân theo **SOLID**:
     - **Single Responsibility**: Formatter xử lý load, parse, flatten; Config quản lý API.
     - **Open/Closed**: Dễ mở rộng formatter/tính năng.
   - Tài liệu API với **godoc**.
3. **Tính tương thích**:
   - Hỗ trợ Go 1.22+ (tính đến 23/05/2025).
   - Phụ thuộc tối thiểu: `gopkg.in/yaml.v3` (YAML), thư viện chuẩn Go (`encoding/json`, `os`, `reflect`, `sync`, `strconv`).
4. **Bảo mật**:
   - Kiểm tra đầu vào (file path, key) để tránh path traversal, key không hợp lệ.
   - Kiểm tra `target` trong `Unmarshal` để tránh panic.
5. **Khả năng mở rộng**:
   - Thêm formatter mới mà không sửa mã hiện tại.
   - Hỗ trợ tính năng nâng cao (validation, hot reload).

### 2.3. Các bên liên quan (Stakeholders)
- **Nhà phát triển ứng dụng**:
  - Sử dụng **Get**, **Set**, **Has**, **Unmarshal** để quản lý cấu hình.
  - Mong muốn API type-safe, dễ tích hợp.
- **DevOps**:
  - Cung cấp biến môi trường, file JSON/YAML.
  - Yêu cầu cấu hình dễ cập nhật.
- **Quản trị hệ thống**:
  - Yêu cầu hot reload để cập nhật cấu hình động.
  - Mong muốn hệ thống an toàn, không lỗi runtime.
- **Nhà phát triển thư viện**:
  - Mở rộng formatter (TOML, XML) hoặc tính năng (validation).

### 2.4. Ràng buộc
- **Thư viện**:
  - `gopkg.in/yaml.v3`: Parse YAML.
  - `encoding/json`: Parse JSON.
  - Thư viện chuẩn Go: `os` (đọc file/env), `reflect` (xử lý động), `sync` (thread-safe), `strconv` (type conversion).
- **Hạn chế**:
  - **Reflection**: Hiệu năng thấp, cần cache.
  - **Key trùng lặp**: Xử lý khi merge.
  - **Type mismatch**: Đảm bảo kiểu dữ liệu khớp.

---

## 3. Cấu trúc giải pháp

### 3.1. Cấu trúc mô-đun
Hệ thống được thiết kế **modular** để đảm bảo bảo trì và mở rộng:
- **`model`**:
  - `ConfigValue`:
    ```go
    type ConfigValue struct {
        Value interface{} // Giá trị gốc (string, int, bool, slice, map)
        Type  ValueType   // Kiểu: TypeString, TypeInt, TypeBool, TypeSlice, TypeMap
    }
    ```
    - **Ghi chú**: `ValueType` là enum để hỗ trợ type-safe access.
  - `ConfigMap`:
    ```go
    type ConfigMap map[string]ConfigValue
    ```
    - **Ghi chú**: Lưu key dot notation và giá trị với kiểu.
- **`formatter`**:
  - Interface `Formatter`:
    ```go
    type Formatter interface {
        Load() (interface{}, error)                     // Tải dữ liệu
        Parse(data interface{}) (interface{}, error)    // Phân tích cú pháp
        Flatten(data interface{}, opts FlattenOptions) (model.ConfigMap, error) // Flatten recursive
    }
    ```
    - **Ghi chú**: Formatter tích hợp **Load**, **Parse**, **Flatten Recursive**.
  - Triển khai: `EnvFormatter`, `JSONFormatter`, `YAMLFormatter`.
- **`config`**:
  - `Config`:
    ```go
    type Config struct {
        formatters []Formatter       // Danh sách formatter
        configMap  model.ConfigMap   // Map phẳng
        mu         sync.RWMutex      // Thread-safe
        opts       FlattenOptions    // Tùy chọn flatten
    }
    ```
    - **Ghi chú**: Quản lý formatter, `ConfigMap`, và cung cấp API.
  - Phương thức:
    - `GetString`, `GetInt`, `GetBool`, `GetSlice`, `GetMap`
    - `Set`, `Has`, `Unmarshal`
- **`utils`**:
  - Hàm tiện ích: Merge map, xử lý tag `config`, type conversion, validation.

### 3.2. Luồng xử lý (High-level Workflow)
1. **Khởi tạo Config**:
   - Tạo `Config` với danh sách formatter (ENV, JSON, YAML) và `FlattenOptions`.
2. **Load dữ liệu**:
   - Mỗi formatter gọi `Load()`:
     - ENV: Đọc `os.Environ()`.
     - JSON/YAML: Đọc file (`os.ReadFile`).
3. **Parse dữ liệu**:
   - Formatter gọi `Parse(data)`:
     - ENV: Chuyển key sang dot notation.
     - JSON: `json.Unmarshal`.
     - YAML: `yaml.Unmarshal`.
4. **Flatten dữ liệu**:
   - Formatter gọi `Flatten(data, opts)` để tạo `ConfigMap` với dot notation.
5. **Merge cấu hình**:
   - Kết hợp `ConfigMap` từ các formatter, ghi đè key theo thứ tự ưu tiên (ENV > JSON > YAML).
6. **Thao tác cấu hình**:
   - **Get**: Truy xuất giá trị type-safe.
   - **Set**: Cập nhật/thêm giá trị.
   - **Has**: Kiểm tra key.
   - **Unmarshal**: Ánh xạ vào struct (toàn bộ hoặc tại key).

### 3.3. Cấu trúc dữ liệu
- **`ConfigValue`**:
  - Lưu giá trị gốc và kiểu (`TypeString`, `TypeInt`, v.v.).
- **`ConfigMap`**:
  - `map[string]ConfigValue`, lưu key dot notation.
- **Go Struct**:
  - Hỗ trợ tag `config` (ví dụ: `config:"database.host"`).
  - Hỗ trợ struct lồng nhau, slice, map.
- **`FlattenOptions`**:
  ```go
  type FlattenOptions struct {
      Separator      string // Separator (mặc định ".")
      SkipNil        bool   // Bỏ qua nil
      HandleEmptyKey bool   // Xử lý key rỗng
      CaseSensitive  bool   // Phân biệt hoa thường
  }
  ```

---

## 4. Tính năng chi tiết

### 4.1. Formatter (ENV, JSON, YAML)

#### 4.1.1. ENV Formatter
- **Load**:
  - Đọc biến môi trường từ `os.Environ()` → `[]string` (ví dụ: `["DB_HOST=localhost"]`).
  - **Thuật toán**: Gọi `os.Environ()`, xử lý lỗi truy cập (hiếm gặp).
  - **Độ phức tạp**: O(n), n là số biến môi trường.
- **Parse**:
  - Chuyển key sang dot notation (ví dụ: `DB_HOST` → `db.host`).
  - **Thuật toán**:
    - Tách key/value từ chuỗi (`DB_HOST=localhost` → `key=DB_HOST`, `value=localhost`).
    - Thay `_` bằng `.` (hoặc separator tùy chỉnh).
    - Chuyển key sang lowercase nếu `CaseSensitive` là `false`.
    - Lưu vào `map[string]interface{}`.
  - **Độ phức tạp**: O(n).
- **Flatten Recursive**:
  - Tạo `ConfigMap` từ key phẳng (ENV thường không lồng nhau).
  - **Ví dụ**:
    ```bash
    export DB_HOST=localhost
    export APP_NAME=myapp
    ```
    → `ConfigMap`:
    ```go
    {
        "db.host": ConfigValue{Value: "localhost", Type: TypeString},
        "app.name": ConfigValue{Value: "myapp", Type: TypeString}
    }
    ```

#### 4.1.2. JSON Formatter
- **Load**:
  - Đọc file JSON (`os.ReadFile`) → `[]byte`.
  - **Thuật toán**: Mở file, đọc nội dung, báo lỗi nếu không tồn tại (`os.ErrNotExist`).
  - **Độ phức tạp**: O(n), n là kích thước file.
- **Parse**:
  - Parse `[]byte` thành `interface{}` (`json.Unmarshal`).
  - **Thuật toán**:
    - Gọi `json.Unmarshal(data, &result)`.
    - Kiểm tra `result` là `map[string]interface{}` hoặc `[]interface{}`.
    - Báo lỗi cú pháp (`ErrParseFailed`).
  - **Độ phức tạp**: O(n), n là số phần tử JSON.
- **Flatten Recursive**:
  - Chuyển cấu trúc lồng nhau thành `ConfigMap`.
  - **Ví dụ**:
    ```json
    {
      "database": {
        "host": "localhost",
        "port": 5432
      },
      "servers": [
        {"host": "server1"},
        {"host": "server2"}
      ]
    }
    ```
    → `ConfigMap`:
    ```go
    {
        "database": ConfigValue{Value: map[string]interface{}{"host": "localhost", "port": 5432}, Type: TypeMap},
        "database.host": ConfigValue{Value: "localhost", Type: TypeString},
        "database.port": ConfigValue{Value: 5432, Type: TypeInt},
        "servers": ConfigValue{Value: []interface{}{map[string]interface{}{"host": "server1"}, map[string]interface{}{"host": "server2"}}, Type: TypeSlice},
        "servers.0.host": ConfigValue{Value: "server1", Type: TypeString},
        "servers.1.host": ConfigValue{Value: "server2", Type: TypeString}
    }
    ```

#### 4.1.3. YAML Formatter
- **Load**:
  - Đọc file YAML (`os.ReadFile`) → `[]byte`.
  - **Độ phức tạp**: O(n).
- **Parse**:
  - Parse `[]byte` thành `interface{}` (`yaml.Unmarshal`).
  - **Thuật toán**: Tương tự JSON, sử dụng `gopkg.in/yaml.v3`.
  - **Độ phức tạp**: O(n).
- **Flatten Recursive**:
  - Tương tự JSON, tạo `ConfigMap` với dot notation.
  - **Ví dụ**:
    ```yaml
    database:
      host: localhost
      port: 5432
    servers:
      - host: server1
      - host: server2
    ```
    → `ConfigMap`: Tương tự JSON.

### 4.2. Phương thức thao tác cấu hình

#### 4.2.1. Get
- **Mục đích**: Truy xuất giá trị type-safe từ `ConfigMap`.
- **Phương thức**:
  ```go
  GetString(key string) (string, bool)
  GetInt(key string) (int, bool)
  GetBool(key string) (bool, bool)
  GetSlice(key string) ([]interface{}, bool)
  GetMap(key string) (map[string]interface{}, bool)
  ```
- **Thuật toán**:
  1. Sử dụng `RLock()` để thread-safe.
  2. Kiểm tra key trong `ConfigMap` (O(1)).
  3. Kiểm tra `ConfigValue.Type` so với kiểu mong muốn.
  4. Chuyển đổi kiểu nếu cần (string → int với `strconv`).
  5. Trả về giá trị và `true`, hoặc `false` nếu key không tồn tại/sai kiểu.
- **Độ phức tạp**: O(1).
- **Ví dụ**:
  ```go
  if val, ok := cfg.GetString("database.host"); ok {
      fmt.Println("Host:", val) // Host: localhost
  }
  ```

#### 4.2.2. Set
- **Mục đích**: Cập nhật/thêm giá trị vào `ConfigMap`.
- **Phương thức**:
  ```go
  Set(key string, value interface{}) error
  ```
- **Thuật toán**:
  1. Sử dụng `Lock()` để thread-safe.
  2. Kiểm tra key hợp lệ (không rỗng nếu `HandleEmptyKey` là `false`).
  3. Xác định kiểu dữ liệu của `value` (`reflect.TypeOf`).
  4. Cập nhật `ConfigMap[key]` với `ConfigValue{Value, Type}`.
- **Độ phức tạp**: O(1).
- **Ví dụ**:
  ```go
  err := cfg.Set("database.host", "newhost")
  // ConfigMap["database.host"] = ConfigValue{Value: "newhost", Type: TypeString}
  ```

#### 4.2.3. Has
- **Mục đích**: Kiểm tra sự tồn tại của key.
- **Phương thức**:
  ```go
  Has(key string) bool
  ```
- **Thuật toán**:
  1. Sử dụng `RLock()` để thread-safe.
  2. Kiểm tra `ConfigMap[key]` (O(1)).
  3. Trả về `true` nếu tồn tại.
- **Độ phức tạp**: O(1).
- **Ví dụ**:
  ```go
  if cfg.Has("database.host") {
      fmt.Println("Key exists")
  }
  ```

#### 4.2.4. Unmarshal
- **Mục đích**: Ánh xạ `ConfigMap` vào Go struct.
- **Phương thức**:
  ```go
  Unmarshal(key string, target interface{}) error
  ```
- **Hành vi**:
  - `key == ""`: Ánh xạ toàn bộ `ConfigMap`.
  - `key != ""`: Ánh xạ dữ liệu tại `key` và key con.
- **Thuật toán**:
  1. **Kiểm tra đầu vào**:
     - `target` là con trỏ đến struct (`reflect`).
     - `key` hợp lệ (không chứa ký tự đặc biệt nếu yêu cầu).
  2. **Khóa thread-safe**:
     - Sử dụng `Lock()`.
  3. **Lọc dữ liệu**:
     - Nếu `key == ""`: Sử dụng toàn bộ `ConfigMap`.
     - Nếu `key != ""`: Lọc key và key con (dùng `strings.HasPrefix`).
  4. **Ánh xạ**:
     - Duyệt field của struct (`reflect`).
     - Tìm key trong `ConfigMap` (tag `config` hoặc tên field).
     - Chuyển đổi kiểu (string → int với `strconv`).
     - Set giá trị vào field.
     - Gọi đệ quy cho struct lồng nhau, slice, map.
  5. **Xử lý lỗi**:
     - Báo lỗi nếu type mismatch, key không tồn tại, target không hợp lệ.
- **Độ phức tạp**: O(k + f), k là số key, f là số field.
- **Ví dụ**:
  ```go
  type Config struct {
      Database struct {
          Host string `config:"database.host"`
          Port int    `config:"database.port"`
      } `config:"database"`
      Servers []struct {
          Host string `config:"host"`
      } `config:"servers"`
  }
  var cfgStruct Config
  cfg.Unmarshal("", &cfgStruct) // Toàn bộ ConfigMap
  // cfgStruct: {Database: {Host: "localhost", Port: 5432}, Servers: [{Host: "server1"}, {Host: "server2"}]}
  var dbStruct struct {
      Host string `config:"database.host"`
      Port int    `config:"database.port"`
  }
  cfg.Unmarshal("database", &dbStruct) // Chỉ database
  // dbStruct: {Host: "localhost", Port: 5432}
  ```

### 4.3. Merge cấu hình
- **Mục đích**: Kết hợp `ConfigMap` từ các formatter theo thứ tự ưu tiên.
- **Thuật toán**:
  1. Gọi `Load`, `Parse`, `Flatten` từ formatter (YAML → JSON → ENV).
  2. Ghi đè key trong `ConfigMap` chính theo thứ tự ưu tiên.
- **Độ phức tạp**: O(k), k là tổng số key.
- **Ví dụ**:
  - YAML: `database.host: localhost`
  - JSON: `database.host: 127.0.0.1`
  - ENV: `DB_HOST=override.local`
  - Kết quả: `ConfigMap["database.host"]` → `ConfigValue{Value: "override.local", Type: TypeString}`.

### 4.4. Thread-safe
- **Cơ chế**:
  - `Lock()`: Load, merge, set, unmarshal.
  - `RLock()`: Get, has.
- **Tối ưu**: `sync.RWMutex` cho read-heavy scenarios.

### 4.5. Error Handling
- **Lỗi**:
  - `ErrInvalidKey`: Key không hợp lệ.
  - `ErrKeyNotFound`: Key không tồn tại.
  - `ErrTypeMismatch`: Kiểu không khớp.
  - `ErrInvalidTarget`: Target không hợp lệ.
  - `ErrParseFailed`: Parse thất bại.
- **Cơ chế**: Wrap lỗi với `fmt.Errorf` và `%w`.

### 4.6. Khả năng mở rộng
- **Formatter mới**: Implement `Formatter` (TOML, XML).
- **Validation**: Kiểm tra schema (`gojsonschema`).
- **Hot Reload**: Theo dõi file (`fsnotify`).
- **Default Values**: Hỗ trợ giá trị mặc định.

---

## 5. Logic triển khai chi tiết

### 5.1. Formatter Logic
- **Load**:
  - ENV: `os.Environ()` → O(n).
  - JSON/YAML: `os.ReadFile` → O(n).
- **Parse**:
  - ENV: Tách key/value, chuyển sang dot notation → O(n).
  - JSON: `json.Unmarshal` → O(n).
  - YAML: `yaml.Unmarshal` → O(n).
- **Flatten Recursive**:
  - **Thuật toán**:
    1. Duyệt dữ liệu bằng `reflect`.
    2. Xử lý kiểu:
       - **Map**: Lưu gốc, tạo key mới (prefix + separator + key).
       - **Slice**: Lưu gốc, tạo key với index.
       - **Scalar**: Lưu vào `ConfigMap`.
    3. Lưu `ConfigValue{Value, Type}`.
  - **Độ phức tạp**: O(n).
  - **Tối ưu**: Sử dụng `strings.Builder` để nối key.

### 5.2. Merge cấu hình
- **Thuật toán**:
  1. Load, parse, flatten từ formatter.
  2. Ghi đè key theo thứ tự ưu tiên.
- **Độ phức tạp**: O(k).

### 5.3. Get
- **Thuật toán**:
  1. `RLock()`.
  2. Kiểm tra `ConfigMap[key]` (O(1)).
  3. Kiểm tra `ConfigValue.Type`.
  4. Chuyển đổi kiểu nếu cần.
- **Độ phức tạp**: O(1).

### 5.4. Set
- **Thuật toán**:
  1. `Lock()`.
  2. Kiểm tra key hợp lệ.
  3. Xác định kiểu của `value`.
  4. Cập nhật `ConfigMap`.
- **Độ phức tạp**: O(1).

### 5.5. Has
- **Thuật toán**:
  1. `RLock()`.
  2. Kiểm tra `ConfigMap[key]` (O(1)).
- **Độ phức tạp**: O(1).

### 5.6. Unmarshal
- **Thuật toán**:
  ```go
  func (c *Config) Unmarshal(key string, target interface{}) error {
      c.mu.Lock()
      defer c.mu.Unlock()
      rv := reflect.ValueOf(target)
      if rv.Kind() != reflect.Ptr || rv.IsNil() {
          return fmt.Errorf("target must be a non-nil pointer: %w", ErrInvalidTarget)
      }
      rv = rv.Elem()
      if rv.Kind() != reflect.Struct {
          return fmt.Errorf("target must point to a struct: %w", ErrInvalidTarget)
      }
      var data ConfigMap
      if key == "" {
          data = c.configMap
      } else {
          data = make(ConfigMap)
          if val, exists := c.configMap[key]; exists {
              data[key] = val
          }
          for k, v := range c.configMap {
              if strings.HasPrefix(k, key+".") {
                  data[k] = v
              }
          }
          if len(data) == 0 {
              return fmt.Errorf("key %s not found: %w", key, ErrKeyNotFound)
          }
      }
      return c.unmarshalStruct(rv, data, key)
  }
  ```
- **unmarshalStruct**: Duyệt field, ánh xạ scalar/struct/slice/map.
- **setField**: Ánh xạ scalar, hỗ trợ type conversion.
- **unmarshalSlice/unmarshalMap**: Xử lý slice/map đệ quy.
- **Độ phức tạp**: O(k + f).
- **Tối ưu**: Cache ánh xạ tag-key.

---

## 6. Phân tích kỹ thuật

### 6.1. Thư viện
- **`gopkg.in/yaml.v3`**: Parse YAML.
- **`encoding/json`**: Parse JSON.
- **`reflect`**: Xử lý động (flatten, unmarshal).
- **`sync.RWMutex`**: Thread-safe.
- **`os`**: Đọc file/env.
- **`strconv`**: Type conversion.
- **`strings`**: Xử lý dot notation.

### 6.2. Độ phức tạp
- **Load/Parse/Flatten**: O(n).
- **Merge**: O(k).
- **Get/Set/Has**: O(1).
- **Unmarshal**: O(k + f).

### 6.3. Tối ưu
- Cache ánh xạ tag-key.
- Sử dụng `strings.Builder` cho dot notation.
- `sync.RWMutex` cho read-heavy.
- Bỏ qua `nil` (`SkipNil`).

### 6.4. Rủi ro và giải pháp
- **Reflection**: Cache để giảm chi phí.
- **Type mismatch**: Kiểm tra kiểu, chuyển đổi an toàn.
- **Key không tồn tại**: Báo `ErrKeyNotFound`.
- **Race condition**: Sử dụng `sync.RWMutex`.

---

## 7. Ví dụ minh họa

### 7.1. ConfigMap
```go
ConfigMap{
    "database.host": ConfigValue{Value: "localhost", Type: TypeString},
    "database.port": ConfigValue{Value: 5432, Type: TypeInt},
    "servers.0.host": ConfigValue{Value: "server1", Type: TypeString},
    "servers.1.host": ConfigValue{Value: "server2", Type: TypeString}
}
```

### 7.2. Sử dụng
```go
type Config struct {
    Database struct {
        Host string `config:"database.host"`
        Port int    `config:"database.port"`
    } `config:"database"`
    Servers []struct {
        Host string `config:"host"`
    } `config:"servers"`
}
cfg := NewConfig([]Formatter{envFormatter, jsonFormatter, yamlFormatter})
var cfgStruct Config
cfg.Unmarshal("", &cfgStruct) // Toàn bộ
// cfgStruct: {Database: {Host: "localhost", Port: 5432}, Servers: [{Host: "server1"}, {Host: "server2"}]}
var dbStruct struct {
    Host string `config:"database.host"`
    Port int    `config:"database.port"`
}
cfg.Unmarshal("database", &dbStruct) // Chỉ database
// dbStruct: {Host: "localhost", Port: 5432}
```

---

## 8. Kết luận
**Service Provider Configuration** đáp ứng yêu cầu:
- **Formatter**: ENV, JSON, YAML với **Load**, **Parse**, **Flatten Recursive**.
- **Phương thức**: **Get**, **Set**, **Has**, **Unmarshal** (toàn bộ hoặc tại key).
- **Flatten Map**: Dot notation, type-safe.
- **Thread-safe**: `sync.RWMutex`.
- **Hiệu năng**: O(1) cho Get/Set/Has, O(k + f) cho Unmarshal.
- **Mở rộng**: Formatter mới, validation, hot reload.
