# MongoDB Provider

## Giới thiệu

MongoDB Provider là một package cung cấp tích hợp MongoDB cho framework dependency injection go-fork. Provider này cung cấp các tính năng quản lý kết nối MongoDB, pool connection và các thao tác cơ bản với MongoDB trong ứng dụng Go. Package này được thiết kế để giúp đơn giản hóa việc tích hợp MongoDB vào ứng dụng Go của bạn, đồng thời hỗ trợ các tính năng nâng cao như transaction và change streams.

## Tổng quan

MongoDB Provider hỗ trợ:
- Tích hợp dễ dàng với framework dependency injection go-fork
- Quản lý kết nối và connection pool
- Hỗ trợ xác thực và SSL/TLS
- Giao diện đơn giản cho các thao tác MongoDB phổ biến
- Hỗ trợ transaction và change streams
- Các tiện ích kiểm tra sức khỏe (health check) và thống kê

## Cài đặt

```bash
go get github.com/go-fork/providers/mongodb
```

## Cấu hình

Sao chép file cấu hình mẫu và chỉnh sửa theo nhu cầu:

```bash
cp configs/app.sample.yaml configs/app.yaml
```

### Ví dụ cấu hình

```yaml
mongodb:
  # Connection URI for MongoDB
  uri: "mongodb://localhost:27017"
  
  # Default database name
  database: "myapp"
  
  # Application name to identify the connection
  app_name: "my-golang-app"
  
  # Connection pool settings
  max_pool_size: 100
  min_pool_size: 5
  max_connecting: 10
  max_conn_idle_time: 600000
  
  # Timeout settings (all in milliseconds)
  connect_timeout: 30000
  server_selection_timeout: 30000
  socket_timeout: 0
  
  # TLS/SSL configuration
  tls:
    enabled: false
    insecure_skip_verify: false
    
  # Authentication configuration
  auth:
    username: ""
    password: ""
    auth_source: "admin"
    auth_mechanism: "SCRAM-SHA-256"
  
  # Read & write concerns
  read_preference:
    mode: "primary"
  read_concern:
    level: "majority"
  write_concern:
    w: "majority"
    journal: true
    
  # Retry configuration
  retry_writes: true
  retry_reads: true
```

## Sử dụng

### Thiết lập cơ bản

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/go-fork/di"
    "github.com/go-fork/providers/config"
    "github.com/go-fork/providers/mongodb"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

func main() {
    // Tạo DI container
    container := di.New()
    
    // Đăng ký provider config (nếu sử dụng service config)
    configProvider := config.NewServiceProvider()
    container.Register(configProvider)
    
    // Đăng ký MongoDB provider
    mongoProvider := mongodb.NewProvider()
    container.Register(mongoProvider)
    
    // Boot các service providers
    container.Boot()
    
    // Lấy MongoDB manager sử dụng MustMake
    // MustMake sẽ panic nếu service không tồn tại hoặc không thể tạo được
    mongoManager := container.MustMake("mongodb").(mongodb.Manager)
    
    // Tạo context với timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Ping database để kiểm tra kết nối
    if err := mongoManager.Ping(ctx); err != nil {
        log.Fatal("Không thể kết nối đến MongoDB:", err)
    }
    
    // Lấy collection
    collection := mongoManager.Collection("users")
    
    // Thêm document
    result, err := collection.InsertOne(ctx, bson.M{
        "name":      "Nguyễn Văn A",
        "email":     "nguyen@example.com",
        "createdAt": time.Now(),
    })
    if err != nil {
        log.Fatal("Không thể thêm document:", err)
    }
    
    log.Printf("Đã thêm document với ID: %v\n", result.InsertedID)
    log.Println("Kết nối MongoDB thành công!")
}
```

### Sử dụng các phương thức của Manager

```go
// Khởi tạo context với timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Kiểm tra sức khỏe
if err := mongoManager.HealthCheck(ctx); err != nil {
    log.Fatal("Kiểm tra sức khỏe MongoDB thất bại:", err)
}

// Lấy thống kê database
stats, err := mongoManager.Stats(ctx)
if err != nil {
    log.Fatal("Không thể lấy thống kê:", err)
}
log.Printf("Thống kê database: %+v\n", stats)

// Liệt kê collections
collections, err := mongoManager.ListCollections(ctx)
if err != nil {
    log.Fatal("Không thể liệt kê collections:", err)
}
log.Printf("Danh sách collections: %v\n", collections)

// Liệt kê tất cả databases
databases, err := mongoManager.ListDatabases(ctx)
if err != nil {
    log.Fatal("Không thể liệt kê databases:", err)
}
log.Printf("Danh sách databases: %v\n", databases)

// Tạo bản ghi sử dụng transaction
// Yêu cầu MongoDB ReplicaSet
result, err := mongoManager.UseSessionWithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
    // Các thao tác transaction ở đây
    collection := mongoManager.Collection("users")
    return collection.InsertOne(sc, bson.M{
        "name":      "Trần Thị B",
        "email":     "tran@example.com",
        "createdAt": time.Now(),
    })
})
if err != nil {
    log.Fatal("Transaction thất bại:", err)
}
```

### Các Services đăng ký

Provider đăng ký các services sau trong DI container:

- `mongodb` - Instance MongoDB Manager
- `mongo.client` - Client MongoDB gốc
- `mongo` - Alias cho MongoDB Manager

Ví dụ truy xuất các services này với MustMake:

```go
// Lấy MongoDB Manager
mongoManager := container.MustMake("mongodb").(mongodb.Manager)

// Lấy MongoDB client gốc
client := container.MustMake("mongo.client").(*mongo.Client)

// Lấy MongoDB Manager thông qua alias
manager := container.MustMake("mongo").(mongodb.Manager)
```

## Phát triển

### Mock cho Testing

Package này cung cấp mock cho việc testing trong thư mục `mocks`. Sử dụng MockManager để test các thành phần phụ thuộc vào MongoDB mà không cần kết nối đến database thật:

```go
import (
    "testing"
    "context"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/go-fork/providers/mongodb/mocks"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

func TestYourFunction(t *testing.T) {
    // Tạo mock manager
    mockManager := mocks.NewMockManager(t)
    
    // Thiết lập expectations cơ bản
    mockManager.On("Ping", mock.Anything).Return(nil)
    mockManager.On("ListDatabases", mock.Anything).Return([]string{"db1", "db2"}, nil)
    
    // Thiết lập mock cho collection và các thao tác CRUD
    mockCollection := &mongo.Collection{}
    mockManager.On("Collection", "users").Return(mockCollection)
    
    // Mock cho ListCollections (với kết quả)
    mockManager.On("ListCollections", mock.Anything).Return([]string{"users", "products"}, nil)
    
    // Mock cho QueryContext với wildcard params
    mockManager.On("QueryContext", mock.Anything, "users", mock.Anything).Return([]bson.M{
        {"_id": "1", "name": "User 1"},
        {"_id": "2", "name": "User 2"},
    }, nil)
    
    // Mock cho hàm nhận vào callback transaction
    mockManager.On("UseSessionWithTransaction", mock.Anything, mock.AnythingOfType("func(mongo.SessionContext) (interface{}, error)")).
        Run(func(args mock.Arguments) {
            // Giả lập thực thi transaction function
            callback := args.Get(1).(func(mongo.SessionContext) (interface{}, error))
            // Tạo SessionContext giả lập
            mockSC := mock.AnythingOfType("mongo.SessionContext").(mongo.SessionContext)
            // Chạy callback
            callback(mockSC)
        }).
        Return(bson.M{"insertedID": "123"}, nil)
    
    // Sử dụng mock trong tests
    err := YourFunction(mockManager)
    
    // Kiểm tra kết quả
    assert.NoError(t, err)
    mockManager.AssertExpectations(t)
}

// Ví dụ về hàm sử dụng MongoDB Manager
func YourFunction(m mongodb.Manager) error {
    ctx := context.Background()
    
    // Ping database
    if err := m.Ping(ctx); err != nil {
        return err
    }
    
    // Liệt kê databases
    dbs, err := m.ListDatabases(ctx)
    if err != nil {
        return err
    }
    
    // Do something with dbs...
    
    return nil
}
```

### Tạo lại Mocks

Mocks được tạo bằng [mockery](https://github.com/vektra/mockery). Để tạo lại mocks, chạy lệnh sau từ thư mục gốc của project:

```bash
mockery
```

Lệnh này sẽ sử dụng cấu hình từ file `.mockery.yaml`.

### Phương pháp cải thiện test coverage

Để cải thiện test coverage của package, hãy chú ý đến các phương pháp sau:

1. **Thiết lập test helper**: Tạo các hàm helper để thiết lập và dọn dẹp môi trường test một cách nhất quán.

2. **Mock external dependencies**: Sử dụng mock cho các dependency bên ngoài như mongo.Client để không phụ thuộc vào MongoDB thật trong unit tests.

3. **Kiểm tra cả happy path và error path**: Đảm bảo kiểm tra cả trường hợp thành công và thất bại của mỗi hàm.

4. **Sử dụng testify**: Sử dụng các assertion của package testify để làm cho tests dễ đọc hơn.

5. **Docker containers cho integration tests**: Sử dụng Docker để chạy MongoDB tạm thời cho integration tests.

Ví dụ thiết lập test helper:

```go
// testHelper.go
package mongodb_test

import (
    "context"
    "testing"
    
    "github.com/go-fork/providers/mongodb"
    "github.com/stretchr/testify/require"
)

func setupTestManager(t *testing.T) (mongodb.Manager, func()) {
    // Tạo config cho test
    config := mongodb.NewConfig()
    config.URI = "mongodb://localhost:27017"
    config.Database = "test_db"
    
    // Khởi tạo manager
    manager, err := mongodb.NewManager(config)
    require.NoError(t, err)
    
    // Tạo cleanup function
    cleanup := func() {
        ctx := context.Background()
        // Xóa database test
        err := manager.DropDatabase(ctx)
        require.NoError(t, err)
        // Đóng kết nối
        err = manager.Disconnect(ctx)
        require.NoError(t, err)
    }
    
    return manager, cleanup
}
```

## Danh sách phương thức

### Các phương thức kết nối

| Phương thức | Mô tả |
|------------|-------|
| `Ping(ctx context.Context) error` | Kiểm tra kết nối tới MongoDB |
| `Disconnect(ctx context.Context) error` | Đóng kết nối tới MongoDB |
| `Client() *mongo.Client` | Lấy MongoDB client gốc |
| `Database() *mongo.Database` | Lấy đối tượng database mặc định |

### Các phương thức quản lý database

| Phương thức | Mô tả |
|------------|-------|
| `ListDatabases(ctx context.Context) ([]string, error)` | Liệt kê tất cả các database |
| `DropDatabase(ctx context.Context) error` | Xóa database mặc định |
| `DropDatabaseWithName(ctx context.Context, name string) error` | Xóa database theo tên |
| `Stats(ctx context.Context) (bson.M, error)` | Lấy thống kê của database |

### Các phương thức quản lý collection

| Phương thức | Mô tả |
|------------|-------|
| `Collection(name string) *mongo.Collection` | Lấy collection theo tên |
| `ListCollections(ctx context.Context) ([]string, error)` | Liệt kê tất cả collection trong database |
| `CreateCollection(ctx context.Context, name string) error` | Tạo collection mới |
| `DropCollection(ctx context.Context, name string) error` | Xóa collection theo tên |
| `RenameCollection(ctx context.Context, oldName, newName string) error` | Đổi tên collection |

### Các phương thức quản lý index

| Phương thức | Mô tả |
|------------|-------|
| `CreateIndex(ctx context.Context, coll string, keys any, opts ...*options.IndexOptions) (string, error)` | Tạo index cho collection |
| `CreateIndexes(ctx context.Context, coll string, models []mongo.IndexModel) ([]string, error)` | Tạo nhiều indexes cho collection |
| `DropIndex(ctx context.Context, coll string, name string) error` | Xóa index theo tên |
| `DropAllIndexes(ctx context.Context, coll string) error` | Xóa tất cả indexes của collection |
| `ListIndexes(ctx context.Context, coll string) ([]bson.M, error)` | Liệt kê tất cả indexes của collection |

### Các phương thức transaction

| Phương thức | Mô tả |
|------------|-------|
| `StartSession() (mongo.Session, error)` | Bắt đầu MongoDB session mới |
| `UseSessionWithTransaction(ctx context.Context, fn func(mongo.SessionContext) (any, error)) (any, error)` | Thực thi một hàm trong transaction |

### Các phương thức tiện ích

| Phương thức | Mô tả |
|------------|-------|
| `HealthCheck(ctx context.Context) error` | Kiểm tra trạng thái kết nối MongoDB |
| `QueryContext(ctx context.Context, coll string, filter any, opts ...*options.FindOptions) ([]bson.M, error)` | Thực hiện truy vấn và trả về kết quả dạng bson.M |

## Lưu ý

1. **Connection Pooling**: Mặc định package tạo một connection pool. Điều chỉnh `maxPoolSize` và `minPoolSize` phù hợp với nhu cầu ứng dụng.

2. **Timeout**: Cân nhắc thiết lập giá trị timeout phù hợp để tránh các kết nối bị treo.

3. **Authentication**: Đảm bảo sử dụng database `authSource` chính xác khi cấu hình xác thực.

4. **Transactions**: Transactions chỉ hoạt động với MongoDB 4.0+ và yêu cầu cấu hình replica set.

5. **Change Streams**: Change streams yêu cầu MongoDB 4.0+ và cấu hình replica set.

6. **Memory Usage**: Theo dõi việc sử dụng bộ nhớ khi làm việc với tập dữ liệu lớn và điều chỉnh `maxPoolSize` phù hợp.

7. **Context Management**: Luôn sử dụng context có timeout khi gọi các phương thức của MongoDB để tránh treo ứng dụng.
