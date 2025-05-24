package queue

import (
	"encoding/json"
	"fmt"
	"time"
)

// Task đại diện cho một tác vụ cần được xử lý.
type Task struct {
	// ID là định danh duy nhất của tác vụ
	ID string

	// Name là tên của loại tác vụ
	Name string

	// Payload là dữ liệu của tác vụ dưới dạng bytes
	Payload []byte

	// Queue là tên của hàng đợi chứa tác vụ
	Queue string

	// MaxRetry là số lần thử lại tối đa nếu tác vụ thất bại
	MaxRetry int

	// RetryCount là số lần tác vụ đã được thử lại
	RetryCount int

	// CreatedAt là thời điểm tác vụ được tạo
	CreatedAt time.Time

	// ProcessAt là thời điểm tác vụ sẽ được xử lý
	ProcessAt time.Time
}

// Unmarshal giải mã payload thành một struct.
func (t *Task) Unmarshal(v interface{}) error {
	return json.Unmarshal(t.Payload, v)
}

// GetName trả về tên của tác vụ.
func (t *Task) GetName() string {
	return t.Name
}

// GetPayload trả về payload của tác vụ.
func (t *Task) GetPayload() []byte {
	return t.Payload
}

// GetID trả về ID của tác vụ.
func (t *Task) GetID() string {
	return t.ID
}

// TaskInfo chứa thông tin về một tác vụ đã được đưa vào hàng đợi.
type TaskInfo struct {
	// ID là định danh duy nhất của tác vụ
	ID string

	// Name là tên của loại tác vụ
	Name string

	// Queue là tên của hàng đợi chứa tác vụ
	Queue string

	// MaxRetry là số lần thử lại tối đa nếu tác vụ thất bại
	MaxRetry int

	// State là trạng thái hiện tại của tác vụ (ví dụ: "pending", "scheduled", "processing", "completed")
	State string

	// CreatedAt là thời điểm tác vụ được tạo
	CreatedAt time.Time

	// ProcessAt là thời điểm tác vụ sẽ được xử lý
	ProcessAt time.Time
}

// NewTask tạo một tác vụ mới với tên và payload được cung cấp.
func NewTask(name string, payload []byte) *Task {
	return &Task{
		Name:      name,
		Payload:   payload,
		CreatedAt: time.Now(),
		ProcessAt: time.Now(),
	}
}

// String trả về biểu diễn chuỗi của TaskInfo.
func (info *TaskInfo) String() string {
	return fmt.Sprintf("TaskInfo{ID: %s, Name: %s, Queue: %s, State: %s, ProcessAt: %v}",
		info.ID, info.Name, info.Queue, info.State, info.ProcessAt)
}
