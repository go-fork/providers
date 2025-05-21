// Package driver cung cấp các interface và thực thi cụ thể cho hệ thống cache.
//
// Package này định nghĩa Driver interface chính và cung cấp các implementation
// cho các loại cache khác nhau như memory, file, redis và mongodb.
//
// Các driver được thiết kế để hoạt động với nhiều loại lưu trữ từ bộ nhớ máy tính đến
// các hệ thống lưu trữ phân tán. Mỗi driver cài đặt đầy đủ các phương thức của
// interface Driver và có thể được sử dụng trực tiếp hoặc thông qua cache.Manager.
package driver

import (
	"context"
	"time"
)

// Driver định nghĩa các thao tác cần thiết cho một cache driver.
// Interface này cung cấp một tập các phương thức tiêu chuẩn cho các
// thao tác lưu trữ và truy xuất dữ liệu từ cache, độc lập với implementation
// cụ thể bên dưới.
type Driver interface {
	// Get lấy một giá trị từ cache.
	//
	// Phương thức này thực hiện tìm kiếm một key trong cache và trả về giá trị
	// tương ứng nếu tìm thấy. Nếu key không tồn tại hoặc đã hết hạn, phương thức
	// sẽ trả về false ở giá trị thứ hai.
	//
	// Params:
	//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
	//   - key: Khóa cần tìm trong cache
	//
	// Returns:
	//   - interface{}: Giá trị được lưu trong cache (nil nếu không tìm thấy)
	//   - bool: true nếu tìm thấy key và chưa hết hạn, false nếu ngược lại
	Get(ctx context.Context, key string) (interface{}, bool)

	// Set đặt một giá trị vào cache với TTL (Time To Live) tùy chọn.
	//
	// Phương thức này lưu trữ một cặp key-value vào cache với thời gian sống
	// được chỉ định. Nếu key đã tồn tại, giá trị sẽ bị ghi đè.
	//
	// Params:
	//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
	//   - key: Khóa để lưu giá trị trong cache
	//   - value: Giá trị cần lưu trữ
	//   - ttl: Thời gian sống của giá trị (0 để sử dụng mặc định của driver, -1 để không hết hạn)
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình thực hiện
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Has kiểm tra xem một key có tồn tại trong cache không.
	//
	// Phương thức này xác định liệu một key có tồn tại trong cache và chưa hết hạn hay không.
	//
	// Params:
	//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
	//   - key: Khóa cần kiểm tra
	//
	// Returns:
	//   - bool: true nếu key tồn tại và chưa hết hạn, false nếu ngược lại
	Has(ctx context.Context, key string) bool

	// Delete xóa một key khỏi cache.
	//
	// Phương thức này loại bỏ một key và giá trị tương ứng khỏi cache nếu tồn tại.
	// Nếu key không tồn tại, thao tác này không có tác dụng và không trả về lỗi.
	//
	// Params:
	//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
	//   - key: Khóa cần xóa
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình thực hiện
	Delete(ctx context.Context, key string) error

	// Flush xóa tất cả các key khỏi cache.
	//
	// Phương thức này xóa tất cả dữ liệu trong cache, làm trống hoàn toàn bộ nhớ cache.
	//
	// Params:
	//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình thực hiện
	Flush(ctx context.Context) error

	// GetMultiple lấy nhiều giá trị từ cache.
	//
	// Phương thức này lấy các giá trị tương ứng với nhiều key trong một lần gọi.
	// Các key không tìm thấy hoặc đã hết hạn sẽ được thêm vào danh sách missed.
	//
	// Params:
	//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
	//   - keys: Danh sách các khóa cần lấy
	//
	// Returns:
	//   - map[string]interface{}: Map chứa các key tìm thấy và giá trị tương ứng
	//   - []string: Danh sách các key không tìm thấy hoặc đã hết hạn
	GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, []string)

	// SetMultiple đặt nhiều giá trị vào cache.
	//
	// Phương thức này lưu trữ nhiều cặp key-value vào cache trong một lần gọi
	// với cùng một thời gian sống.
	//
	// Params:
	//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
	//   - values: Map chứa các key và giá trị tương ứng cần lưu trữ
	//   - ttl: Thời gian sống chung cho tất cả các giá trị
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình thực hiện
	SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error

	// DeleteMultiple xóa nhiều key khỏi cache.
	//
	// Phương thức này xóa nhiều key và giá trị tương ứng khỏi cache trong một lần gọi.
	//
	// Params:
	//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
	//   - keys: Danh sách các khóa cần xóa
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình thực hiện
	DeleteMultiple(ctx context.Context, keys []string) error

	// Remember lấy một giá trị từ cache hoặc thực thi callback nếu không tìm thấy.
	//
	// Phương thức này kiểm tra xem một key có tồn tại trong cache không, nếu có thì
	// trả về giá trị tương ứng. Nếu key không tồn tại hoặc đã hết hạn, phương thức
	// sẽ gọi hàm callback để lấy dữ liệu, lưu kết quả vào cache và trả về giá trị đó.
	//
	// Params:
	//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
	//   - key: Khóa cần tìm hoặc lưu vào cache
	//   - ttl: Thời gian sống của giá trị nếu phải lấy từ callback
	//   - callback: Hàm được gọi để lấy dữ liệu khi key không có trong cache
	//
	// Returns:
	//   - interface{}: Giá trị từ cache hoặc từ callback
	//   - error: Lỗi nếu có trong quá trình thực hiện hoặc từ callback
	Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error)

	// Stats trả về thông tin thống kê về cache.
	//
	// Phương thức này thu thập và trả về các thông tin thống kê về trạng thái
	// hiện tại của cache như số lượng item, dung lượng, số lần hit/miss, v.v.
	//
	// Params:
	//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
	//
	// Returns:
	//   - map[string]interface{}: Map chứa các thông tin thống kê
	Stats(ctx context.Context) map[string]interface{}

	// Close giải phóng tài nguyên của driver.
	//
	// Phương thức này giải phóng các tài nguyên được sử dụng bởi driver
	// như kết nối, file handlers, goroutines, v.v.
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình giải phóng tài nguyên
	Close() error
}
