// Package mongodb cung cấp một service provider cho MongoDB trong framework dependency injection go-fork.
//
// Package này tích hợp MongoDB vào ứng dụng Go thông qua cơ chế dependency injection,
// giúp đơn giản hóa việc quản lý kết nối MongoDB và cung cấp nhiều phương thức tiện ích
// để làm việc với cơ sở dữ liệu MongoDB.
//
// # Tổng quan
//
// Package mongodb cung cấp các tính năng chính sau:
//   - Tích hợp với DI container của go-fork thông qua ServiceProvider
//   - Quản lý kết nối MongoDB và connection pooling
//   - Hỗ trợ xác thực và SSL/TLS
//   - Interface đơn giản cho các thao tác MongoDB phổ biến
//   - Hỗ trợ transaction và change streams
//   - Các tiện ích kiểm tra sức khỏe (health check) và thống kê
//   - Mock cho kiểm thử (testing)
//
// # Cấu hình
//
// Package mongodb có thể được cấu hình thông qua YAML hoặc trực tiếp qua code. Dưới đây là một ví dụ cấu hình YAML:
//
//	mongodb:
//	  uri: "mongodb://localhost:27017"
//	  database: "myapp"
//	  app_name: "my-golang-app"
//	  max_pool_size: 100
//	  min_pool_size: 5
//	  connect_timeout: 30000
//	  auth:
//	    username: ""
//	    password: ""
//	    auth_source: "admin"
//	    auth_mechanism: "SCRAM-SHA-256"
//
// # Đăng ký Provider
//
// Đăng ký MongoDB provider với DI container của go-fork:
//
//	container := di.New()
//	mongoProvider := mongodb.NewProvider()
//	container.Register(mongoProvider)
//	container.Boot()
//
// # Sử dụng Manager
//
// Sau khi đăng ký, bạn có thể truy xuất MongoDB manager từ container:
//
//	mongoManager := container.MustMake("mongodb").(mongodb.Manager)
//
//	// Sử dụng manager
//	ctx := context.Background()
//	err := mongoManager.Ping(ctx)
//
//	// Thao tác với collection
//	collection := mongoManager.Collection("users")
//	result, err := collection.InsertOne(ctx, bson.M{"name": "Example"})
//
// # Các dịch vụ được đăng ký
//
// Provider đăng ký các dịch vụ sau vào container:
//   - "mongodb" - Instance của MongoDB Manager (kiểu mongodb.Manager)
//   - "mongo.client" - Client MongoDB gốc (kiểu *mongo.Client)
//   - "mongo" - Alias cho MongoDB Manager
//
// # Testing với Mocks
//
// Package này cung cấp mock cho interface Manager để hỗ trợ việc kiểm thử:
//
//	mockManager := mocks.NewMockManager(t)
//	mockManager.On("Ping", mock.Anything).Return(nil)
//	mockManager.On("ListDatabases", mock.Anything).Return([]string{"db1", "db2"}, nil)
//
//	// Test với mock
//	err := YourFunction(mockManager)
//	mockManager.AssertExpectations(t)
//
// # Yêu cầu
//
// - Go 1.23.x hoặc mới hơn
// - MongoDB 4.0 hoặc mới hơn (cho transactions và change streams)
//
// # Lưu ý
//
// - Transactions yêu cầu MongoDB 4.0+ và cấu hình replica set
// - Change Streams yêu cầu MongoDB 4.0+ và cấu hình replica set
// - Luôn sử dụng context.Context có timeout khi gọi các phương thức MongoDB
//
// Để biết thêm thông tin chi tiết, xem README.md hoặc mã nguồn của package.
package mongodb
