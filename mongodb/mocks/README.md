# MongoDB Provider Mocks

Thư mục này chứa các triển khai mock cho các interface định nghĩa trong gói MongoDB provider.

## Mock có sẵn

- **MockManager**: Triển khai mock của interface `Manager` để kiểm thử các thao tác MongoDB.

## Cách sử dụng

### MockManager

```go
import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "go.fork.vn/providers/mongodb/mocks"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

func TestYourFunction(t *testing.T) {
    // Tạo mock manager mới
    mockManager := mocks.NewMockManager(t)
    
    // Thiết lập expectations
    mockManager.On("Ping", mock.Anything).Return(nil)
    mockManager.On("ListDatabases", mock.Anything).Return([]string{"db1", "db2"}, nil)
    
    // Mock các phương thức CRUD phức tạp hơn
    mockCollection := &mongo.Collection{} // Bạn có thể cần mock thêm cho đối tượng này
    mockManager.On("Collection", "users").Return(mockCollection)
    
    mockManager.On("QueryContext", mock.Anything, "users", mock.Anything).Return([]bson.M{
        {"_id": "1", "name": "User 1"},
        {"_id": "2", "name": "User 2"},
    }, nil)
    
    // Sử dụng mock trong tests của bạn
    err := YourFunction(mockManager)
    
    // Kiểm tra kết quả
    assert.NoError(t, err)
    mockManager.AssertExpectations(t)
}

// Hàm cần kiểm thử
func YourFunction(manager mongodb.Manager) error {
    ctx := context.Background()
    
    // Gọi các phương thức của MongoDB
    if err := manager.Ping(ctx); err != nil {
        return err
    }
    
    // Truy vấn dữ liệu
    users, err := manager.QueryContext(ctx, "users", bson.M{"active": true})
    if err != nil {
        return err
    }
    
    // Xử lý dữ liệu...
    
    return nil
}
```

## Tạo lại Mocks

Các mock trong thư mục này được tạo bằng [mockery](https://github.com/vektra/mockery).
Để tạo lại mocks, chạy lệnh sau từ thư mục gốc của dự án:

```bash
mockery
```

Lệnh này sẽ sử dụng cấu hình từ file `.mockery.yaml` trong thư mục gốc.

## Cấu hình

Việc tạo mock được cấu hình trong file `.mockery.yaml` ở thư mục gốc của dự án.
Cấu hình hiện tại tạo mock cho interface `Manager`.

## Mock các phương thức phức tạp

### Mock transactions

```go
mockManager.On("UseSessionWithTransaction", mock.Anything, mock.AnythingOfType("func(mongo.SessionContext) (interface{}, error)")).
    Run(func(args mock.Arguments) {
        // Trích xuất hàm callback transaction
        callback := args.Get(1).(func(mongo.SessionContext) (interface{}, error))
        
        // Tạo SessionContext mock
        mockSC := &mockSessionContext{} // Cần tạo một struct giả lập SessionContext
        
        // Gọi callback với session context giả lập
        callback(mockSC)
    }).
    Return(bson.M{"insertedID": "123"}, nil)
```

### Mock cursor và kết quả truy vấn

```go
mockCursor := &mongo.Cursor{} // Mock cursor
mockManager.On("Find", mock.Anything, "users", mock.Anything).Return(mockCursor, nil)
```

## Thông tin thêm

- Đối tượng mock triển khai cùng interface như triển khai thật nhưng cho phép bạn kiểm soát hành vi của chúng để kiểm thử.
- Sử dụng `On()` để thiết lập mong đợi cho các lời gọi phương thức.
- Sử dụng `Return()` để chỉ định giá trị nên được trả về khi phương thức được gọi.
- Sử dụng `Run()` để thực thi mã tùy chỉnh khi một phương thức được gọi.
- Sử dụng `AssertExpectations()` để xác minh rằng tất cả phương thức mong đợi đã được gọi.
- Sử dụng `mock.Anything` cho các tham số không quan trọng trong kiểm thử.
- Sử dụng `mock.AnythingOfType()` khi cần chỉ định kiểu dữ liệu nhưng không quan tâm giá trị cụ thể.
