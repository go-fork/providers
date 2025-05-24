// Package adapter cung cấp các triển khai khác nhau cho queue backend.
package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisQueue triển khai interface QueueAdapter sử dụng Redis.
// Struct này sử dụng Redis List để lưu trữ và quản lý hàng đợi,
// cho phép các hoạt động queue có tính mở rộng cao và phân tán.
type redisQueue struct {
	client *redis.Client
	prefix string
}

// NewRedisQueue tạo một instance mới của redisQueue.
// Hàm này khởi tạo kết nối Redis và áp dụng prefix cho key.
//
// Trả về:
//   - QueueAdapter: Instance mới của redisQueue
func NewRedisQueue(client *redis.Client, prefix string) QueueAdapter {
	if prefix == "" {
		prefix = "queue:"
	}
	return &redisQueue{
		client: client,
		prefix: prefix,
	}
}

// prefixKey thêm prefix đã cấu hình vào tên hàng đợi.
// Hàm này đảm bảo tất cả các key queue trong Redis có cùng prefix
// để dễ dàng quản lý và tránh xung đột tên.
//
// Tham số:
//   - queueName (string): Tên gốc của hàng đợi
//
// Trả về:
//   - string: Tên hàng đợi có prefix
func (q *redisQueue) prefixKey(queueName string) string {
	return q.prefix + queueName
}

// Enqueue thêm một item vào cuối hàng đợi.
// Hàm này serialize item thành JSON và thêm vào Redis list
// sử dụng lệnh RPUSH.
//
// Tham số:
//   - ctx (context.Context): Context cho request
//   - queueName (string): Tên của hàng đợi
//   - item (interface{}): Đối tượng cần đưa vào hàng đợi
//
// Trả về:
//   - error: Lỗi nếu có khi thêm item vào hàng đợi
func (q *redisQueue) Enqueue(ctx context.Context, queueName string, item interface{}) error {
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("error marshaling queue item: %w", err)
	}

	return q.client.RPush(ctx, q.prefixKey(queueName), data).Err()
}

// Dequeue lấy và xóa item ở đầu hàng đợi.
// Hàm này sử dụng lệnh LPOP của Redis để lấy item đầu tiên
// từ list, sau đó deserialize thành đối tượng định sẵn.
//
// Tham số:
//   - ctx (context.Context): Context cho request
//   - queueName (string): Tên của hàng đợi
//   - dest (interface{}): Con trỏ đến đối tượng sẽ nhận dữ liệu
//
// Trả về:
//   - error: Lỗi nếu có khi lấy item từ hàng đợi hoặc khi hàng đợi rỗng
func (q *redisQueue) Dequeue(ctx context.Context, queueName string, dest interface{}) error {
	data, err := q.client.LPop(ctx, q.prefixKey(queueName)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("queue is empty: %s", queueName)
		}
		return err
	}

	return json.Unmarshal(data, dest)
}

// EnqueueBatch thêm nhiều item vào cuối hàng đợi.
// Hàm này serialize mỗi item thành JSON và thêm tất cả
// vào Redis list trong một lệnh RPUSH duy nhất.
//
// Tham số:
//   - ctx (context.Context): Context cho request
//   - queueName (string): Tên của hàng đợi
//   - items ([]interface{}): Slice các đối tượng cần đưa vào hàng đợi
//
// Trả về:
//   - error: Lỗi nếu có khi thêm items vào hàng đợi
func (q *redisQueue) EnqueueBatch(ctx context.Context, queueName string, items []interface{}) error {
	if len(items) == 0 {
		return nil
	}

	// Convert items to JSON strings
	values := make([]interface{}, len(items))
	for i, item := range items {
		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("error marshaling queue item at index %d: %w", i, err)
		}
		values[i] = data
	}

	return q.client.RPush(ctx, q.prefixKey(queueName), values...).Err()
}

// Size trả về số lượng item trong hàng đợi.
// Hàm này sử dụng lệnh LLEN của Redis để đếm số lượng
// phần tử trong list.
//
// Tham số:
//   - ctx (context.Context): Context cho request
//   - queueName (string): Tên của hàng đợi
//
// Trả về:
//   - int64: Số lượng item trong hàng đợi
//   - error: Lỗi nếu có khi truy vấn Redis
func (q *redisQueue) Size(ctx context.Context, queueName string) (int64, error) {
	return q.client.LLen(ctx, q.prefixKey(queueName)).Result()
}

// IsEmpty kiểm tra xem hàng đợi có rỗng không.
// Hàm này kiểm tra kích thước của hàng đợi và trả về true
// nếu kích thước bằng 0.
//
// Tham số:
//   - ctx (context.Context): Context cho request
//   - queueName (string): Tên của hàng đợi
//
// Trả về:
//   - bool: true nếu hàng đợi rỗng, ngược lại là false
//   - error: Lỗi nếu có khi kiểm tra
func (q *redisQueue) IsEmpty(ctx context.Context, queueName string) (bool, error) {
	size, err := q.Size(ctx, queueName)
	if err != nil {
		return false, err
	}
	return size == 0, nil
}

// Clear xóa tất cả các item trong hàng đợi.
// Hàm này sử dụng lệnh DEL của Redis để xóa hoàn toàn
// Redis list tương ứng với hàng đợi.
//
// Tham số:
//   - ctx (context.Context): Context cho request
//   - queueName (string): Tên của hàng đợi
//
// Trả về:
//   - error: Lỗi nếu có khi xóa hàng đợi
func (q *redisQueue) Clear(ctx context.Context, queueName string) error {
	return q.client.Del(ctx, q.prefixKey(queueName)).Err()
}

// DequeueWithTimeout lấy và xóa item ở đầu hàng đợi, với khả năng chờ đợi
// nếu hàng đợi đang rỗng. Hàm này sử dụng lệnh BLPOP của Redis để chờ tối đa
// một khoảng thời gian nhất định cho đến khi có item mới trong hàng đợi.
//
// Tham số:
//   - ctx (context.Context): Context cho request
//   - queueName (string): Tên của hàng đợi
//   - timeout (time.Duration): Thời gian tối đa chờ đợi
//   - dest (interface{}): Con trỏ đến đối tượng sẽ nhận dữ liệu
//
// Trả về:
//   - error: Lỗi nếu có khi lấy item hoặc khi hết thời gian chờ
func (q *redisQueue) DequeueWithTimeout(ctx context.Context, queueName string, timeout time.Duration, dest interface{}) error {
	data, err := q.client.BLPop(ctx, timeout, q.prefixKey(queueName)).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("timeout waiting for queue item: %s", queueName)
		}
		return err
	}

	// BLPop returns a slice with [key, value]
	if len(data) < 2 {
		return fmt.Errorf("unexpected response format from redis")
	}

	return json.Unmarshal([]byte(data[1]), dest)
}
