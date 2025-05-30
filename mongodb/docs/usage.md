# MongoDB Provider Usage Guide

## Quick Start

### 1. Basic Installation & Setup

```bash
# Install package
go get go.fork.vn/providers/mongodb@v0.1.0
```

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "go.fork.vn/providers/mongodb"
)

func main() {
    // Create configuration
    config := &mongodb.Config{
        URI:      "mongodb://localhost:27017",
        Database: "myapp",
    }
    
    // Create manager
    manager := mongodb.NewManager(config)
    defer manager.Disconnect(context.Background())
    
    // Test connection
    ctx := context.Background()
    if err := manager.Ping(ctx); err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }
    
    fmt.Println("Connected to MongoDB successfully!")
}
```

### 2. Configuration với Environment Variables

```bash
# .env file
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=myapp
MONGODB_MAX_POOL_SIZE=50
```

```go
import (
    "go.fork.vn/providers/config"
    "go.fork.vn/providers/mongodb"
)

func setupMongoDB() mongodb.Manager {
    // Load config từ environment
    configManager := config.NewManager()
    configManager.LoadFromEnv()
    
    // Parse MongoDB config
    var mongoConfig mongodb.Config
    if err := configManager.UnmarshalKey("mongodb", &mongoConfig); err != nil {
        log.Fatal("Failed to parse MongoDB config:", err)
    }
    
    return mongodb.NewManager(&mongoConfig)
}
```

## Dependency Injection Integration

### 1. Complete DI Setup

```go
package main

import (
    "context"
    "log"
    
    "go.fork.vn/di"
    "go.fork.vn/providers/config"
    "go.fork.vn/providers/mongodb"
)

func main() {
    container := di.NewContainer()
    
    // Register providers
    container.RegisterServiceProvider(config.NewServiceProvider())
    container.RegisterServiceProvider(mongodb.NewServiceProvider())
    
    // Build container
    if err := container.Build(); err != nil {
        log.Fatal("Failed to build container:", err)
    }
    
    // Use MongoDB
    var mongoManager mongodb.Manager
    if err := container.Resolve(&mongoManager); err != nil {
        log.Fatal("Failed to resolve MongoDB manager:", err)
    }
    
    // Test connection
    ctx := context.Background()
    if err := mongoManager.Ping(ctx); err != nil {
        log.Fatal("MongoDB ping failed:", err)
    }
    
    log.Println("MongoDB ready!")
}
```

### 2. Service Layer Integration

```go
// models/user.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name     string             `bson:"name" json:"name"`
    Email    string             `bson:"email" json:"email"`
    Age      int                `bson:"age" json:"age"`
    Active   bool               `bson:"active" json:"active"`
}
```

```go
// services/user_service.go
package services

import (
    "context"
    "fmt"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.fork.vn/providers/mongodb"
    "yourapp/models"
)

type UserService struct {
    mongoManager mongodb.Manager
}

func NewUserService(mongoManager mongodb.Manager) *UserService {
    return &UserService{
        mongoManager: mongoManager,
    }
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
    collection, err := s.mongoManager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return fmt.Errorf("failed to get collection: %w", err)
    }
    
    result, err := collection.InsertOne(ctx, user)
    if err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }
    
    user.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*models.User, error) {
    collection, err := s.mongoManager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return nil, fmt.Errorf("failed to get collection: %w", err)
    }
    
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }
    
    var user models.User
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    
    return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, updates bson.M) error {
    collection, err := s.mongoManager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return fmt.Errorf("failed to get collection: %w", err)
    }
    
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("invalid user ID: %w", err)
    }
    
    _, err = collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$set": updates},
    )
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
    collection, err := s.mongoManager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return fmt.Errorf("failed to get collection: %w", err)
    }
    
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("invalid user ID: %w", err)
    }
    
    _, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    return nil
}

func (s *UserService) ListUsers(ctx context.Context, limit int64) ([]*models.User, error) {
    collection, err := s.mongoManager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return nil, fmt.Errorf("failed to get collection: %w", err)
    }
    
    cursor, err := collection.Find(ctx, bson.M{}, 
        options.Find().SetLimit(limit))
    if err != nil {
        return nil, fmt.Errorf("failed to find users: %w", err)
    }
    defer cursor.Close(ctx)
    
    var users []*models.User
    if err = cursor.All(ctx, &users); err != nil {
        return nil, fmt.Errorf("failed to decode users: %w", err)
    }
    
    return users, nil
}
```

### 3. Register Services với DI

```go
// main.go
package main

import (
    "go.fork.vn/di"
    "yourapp/services"
)

func main() {
    container := di.NewContainer()
    
    // Register providers
    container.RegisterServiceProvider(config.NewServiceProvider())
    container.RegisterServiceProvider(mongodb.NewServiceProvider())
    
    // Register services
    container.RegisterSingleton(func(mongoManager mongodb.Manager) *services.UserService {
        return services.NewUserService(mongoManager)
    })
    
    // Build container
    if err := container.Build(); err != nil {
        log.Fatal("Failed to build container:", err)
    }
    
    // Use service
    var userService *services.UserService
    if err := container.Resolve(&userService); err != nil {
        log.Fatal("Failed to resolve user service:", err)
    }
    
    // Create user
    ctx := context.Background()
    user := &models.User{
        Name:   "John Doe",
        Email:  "john@example.com",
        Age:    30,
        Active: true,
    }
    
    if err := userService.CreateUser(ctx, user); err != nil {
        log.Fatal("Failed to create user:", err)
    }
    
    log.Printf("Created user with ID: %s", user.ID.Hex())
}
```

## Configuration Examples

### 1. Development Configuration

```yaml
# config/app.yaml
mongodb:
  uri: "mongodb://localhost:27017"
  database: "myapp_dev"
  max_pool_size: 10
  min_pool_size: 2
  max_conn_idle_time: "30s"
  connect_timeout: "10s"
  socket_timeout: "30s"
```

### 2. Production Configuration với SSL

```yaml
# config/app.yaml
mongodb:
  uri: "mongodb+srv://cluster.mongodb.net"
  database: "myapp_prod"
  max_pool_size: 100
  min_pool_size: 20
  max_conn_idle_time: "60s"
  server_selection_timeout: "30s"
  connect_timeout: "10s"
  socket_timeout: "30s"
  ssl:
    enabled: true
    ca_file: "/etc/ssl/certs/mongodb-ca.pem"
    certificate_file: "/etc/ssl/certs/mongodb-cert.pem"
    private_key_file: "/etc/ssl/private/mongodb-key.pem"
    insecure_skip_verify: false
  auth:
    username: "${MONGODB_USERNAME}"
    password: "${MONGODB_PASSWORD}"
    auth_db: "admin"
```

### 3. Environment Variables Setup

```bash
# .env
MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net
MONGODB_DATABASE=myapp_prod
MONGODB_MAX_POOL_SIZE=100
MONGODB_MIN_POOL_SIZE=20
MONGODB_MAX_CONN_IDLE_TIME=60s
MONGODB_SERVER_SELECTION_TIMEOUT=30s
MONGODB_CONNECT_TIMEOUT=10s
MONGODB_SOCKET_TIMEOUT=30s

# SSL Settings
MONGODB_SSL_ENABLED=true
MONGODB_SSL_CA_FILE=/etc/ssl/certs/mongodb-ca.pem
MONGODB_SSL_CERTIFICATE_FILE=/etc/ssl/certs/mongodb-cert.pem
MONGODB_SSL_PRIVATE_KEY_FILE=/etc/ssl/private/mongodb-key.pem
MONGODB_SSL_INSECURE_SKIP_VERIFY=false

# Authentication
MONGODB_AUTH_USERNAME=admin
MONGODB_AUTH_PASSWORD=password123
MONGODB_AUTH_DB=admin
```

## Advanced Operations

### 1. Transactions

```go
func (s *UserService) TransferCredits(ctx context.Context, fromUserID, toUserID string, amount int) error {
    client, err := s.mongoManager.GetClient(ctx)
    if err != nil {
        return fmt.Errorf("failed to get client: %w", err)
    }
    
    // Start session
    session, err := client.StartSession()
    if err != nil {
        return fmt.Errorf("failed to start session: %w", err)
    }
    defer session.EndSession(ctx)
    
    // Execute transaction
    _, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
        collection, err := s.mongoManager.GetCollection(sc, "myapp", "users")
        if err != nil {
            return nil, err
        }
        
        // Deduct from sender
        fromObjID, _ := primitive.ObjectIDFromHex(fromUserID)
        _, err = collection.UpdateOne(sc, 
            bson.M{"_id": fromObjID}, 
            bson.M{"$inc": bson.M{"credits": -amount}})
        if err != nil {
            return nil, fmt.Errorf("failed to deduct credits: %w", err)
        }
        
        // Add to receiver
        toObjID, _ := primitive.ObjectIDFromHex(toUserID)
        _, err = collection.UpdateOne(sc, 
            bson.M{"_id": toObjID}, 
            bson.M{"$inc": bson.M{"credits": amount}})
        if err != nil {
            return nil, fmt.Errorf("failed to add credits: %w", err)
        }
        
        return nil, nil
    })
    
    return err
}
```

### 2. Aggregation Pipeline

```go
func (s *UserService) GetUserStats(ctx context.Context) (*UserStats, error) {
    collection, err := s.mongoManager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return nil, fmt.Errorf("failed to get collection: %w", err)
    }
    
    pipeline := []bson.M{
        {"$group": bson.M{
            "_id": nil,
            "total_users": bson.M{"$sum": 1},
            "active_users": bson.M{
                "$sum": bson.M{
                    "$cond": []interface{}{
                        "$active", 1, 0,
                    },
                },
            },
            "avg_age": bson.M{"$avg": "$age"},
            "max_age": bson.M{"$max": "$age"},
            "min_age": bson.M{"$min": "$age"},
        }},
    }
    
    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, fmt.Errorf("failed to aggregate: %w", err)
    }
    defer cursor.Close(ctx)
    
    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("failed to decode results: %w", err)
    }
    
    if len(results) == 0 {
        return &UserStats{}, nil
    }
    
    result := results[0]
    return &UserStats{
        TotalUsers:  result["total_users"].(int32),
        ActiveUsers: result["active_users"].(int32),
        AvgAge:      result["avg_age"].(float64),
        MaxAge:      result["max_age"].(int32),
        MinAge:      result["min_age"].(int32),
    }, nil
}

type UserStats struct {
    TotalUsers  int32   `json:"total_users"`
    ActiveUsers int32   `json:"active_users"`
    AvgAge      float64 `json:"avg_age"`
    MaxAge      int32   `json:"max_age"`
    MinAge      int32   `json:"min_age"`
}
```

### 3. Indexing

```go
func (s *UserService) EnsureIndexes(ctx context.Context) error {
    collection, err := s.mongoManager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return fmt.Errorf("failed to get collection: %w", err)
    }
    
    indexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys: bson.D{
                {Key: "name", Value: 1},
                {Key: "age", Value: -1},
            },
        },
        {
            Keys: bson.D{{Key: "active", Value: 1}},
            Options: options.Index().SetPartialFilterExpression(
                bson.M{"active": true},
            ),
        },
    }
    
    _, err = collection.Indexes().CreateMany(ctx, indexes)
    if err != nil {
        return fmt.Errorf("failed to create indexes: %w", err)
    }
    
    return nil
}
```

## Testing Examples

### 1. Unit Testing với Mocks

```go
package services_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "go.mongodb.org/mongo-driver/mongo"
    
    "go.fork.vn/providers/mongodb/mocks"
    "yourapp/models"
    "yourapp/services"
)

func TestUserService_CreateUser(t *testing.T) {
    // Setup mock
    mockManager := &mocks.Manager{}
    mockCollection := &mongo.Collection{} // Mock this too
    
    mockManager.On("GetCollection", mock.Anything, "myapp", "users").
        Return(mockCollection, nil)
    
    // Test service
    service := services.NewUserService(mockManager)
    
    user := &models.User{
        Name:   "Test User",
        Email:  "test@example.com",
        Age:    25,
        Active: true,
    }
    
    err := service.CreateUser(context.Background(), user)
    
    assert.NoError(t, err)
    mockManager.AssertExpectations(t)
}
```

### 2. Integration Testing

```go
package integration_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "go.fork.vn/providers/mongodb"
    "yourapp/models"
    "yourapp/services"
)

func TestUserService_Integration(t *testing.T) {
    // Setup test MongoDB
    config := &mongodb.Config{
        URI:      "mongodb://localhost:27017",
        Database: "test_db",
    }
    
    manager := mongodb.NewManager(config)
    defer manager.Disconnect(context.Background())
    
    // Test connection
    ctx := context.Background()
    require.NoError(t, manager.Ping(ctx))
    
    // Create service
    service := services.NewUserService(manager)
    
    // Test create user
    user := &models.User{
        Name:   "Integration Test User",
        Email:  "integration@test.com",
        Age:    30,
        Active: true,
    }
    
    err := service.CreateUser(ctx, user)
    require.NoError(t, err)
    require.NotEmpty(t, user.ID)
    
    // Test get user
    retrievedUser, err := service.GetUser(ctx, user.ID.Hex())
    require.NoError(t, err)
    assert.Equal(t, user.Name, retrievedUser.Name)
    assert.Equal(t, user.Email, retrievedUser.Email)
    
    // Test update user
    err = service.UpdateUser(ctx, user.ID.Hex(), bson.M{"age": 31})
    require.NoError(t, err)
    
    // Verify update
    updatedUser, err := service.GetUser(ctx, user.ID.Hex())
    require.NoError(t, err)
    assert.Equal(t, 31, updatedUser.Age)
    
    // Test delete user
    err = service.DeleteUser(ctx, user.ID.Hex())
    require.NoError(t, err)
    
    // Verify deletion
    _, err = service.GetUser(ctx, user.ID.Hex())
    assert.Error(t, err) // Should not find user
}
```

## Health Checks

### 1. Basic Health Check

```go
func (s *UserService) HealthCheck(ctx context.Context) error {
    // Check MongoDB connection
    if !s.mongoManager.IsConnected(ctx) {
        return errors.New("mongodb not connected")
    }
    
    // Ping database
    if err := s.mongoManager.Ping(ctx); err != nil {
        return fmt.Errorf("mongodb ping failed: %w", err)
    }
    
    return nil
}
```

### 2. Detailed Health Check

```go
type HealthStatus struct {
    MongoDB DatabaseHealth `json:"mongodb"`
}

type DatabaseHealth struct {
    Status      string        `json:"status"`
    Latency     time.Duration `json:"latency"`
    Error       string        `json:"error,omitempty"`
    Connected   bool          `json:"connected"`
    Database    string        `json:"database"`
    Collections int           `json:"collections"`
}

func (s *UserService) DetailedHealthCheck(ctx context.Context) *HealthStatus {
    health := &HealthStatus{}
    
    start := time.Now()
    
    // Check connection
    health.MongoDB.Connected = s.mongoManager.IsConnected(ctx)
    health.MongoDB.Database = "myapp"
    
    if !health.MongoDB.Connected {
        health.MongoDB.Status = "disconnected"
        health.MongoDB.Error = "not connected to database"
        return health
    }
    
    // Ping database
    if err := s.mongoManager.Ping(ctx); err != nil {
        health.MongoDB.Status = "unhealthy"
        health.MongoDB.Error = err.Error()
        return health
    }
    
    health.MongoDB.Latency = time.Since(start)
    
    // Get collection count
    db, err := s.mongoManager.GetDatabase(ctx, "myapp")
    if err != nil {
        health.MongoDB.Status = "unhealthy"
        health.MongoDB.Error = err.Error()
        return health
    }
    
    collections, err := db.ListCollectionNames(ctx, bson.M{})
    if err != nil {
        health.MongoDB.Status = "unhealthy"
        health.MongoDB.Error = err.Error()
        return health
    }
    
    health.MongoDB.Collections = len(collections)
    health.MongoDB.Status = "healthy"
    
    return health
}
```

## Performance Monitoring

### 1. Connection Pool Monitoring

```go
func (s *UserService) GetConnectionStats(ctx context.Context) map[string]interface{} {
    client, err := s.mongoManager.GetClient(ctx)
    if err != nil {
        return map[string]interface{}{
            "error": err.Error(),
        }
    }
    
    return map[string]interface{}{
        "connection_string": s.mongoManager.GetConnectionString(),
        "connected":         s.mongoManager.IsConnected(ctx),
        // Add more stats as needed
    }
}
```

### 2. Query Performance Logging

```go
func (s *UserService) FindUsersWithLogging(ctx context.Context, filter bson.M) ([]*models.User, error) {
    start := time.Now()
    
    collection, err := s.mongoManager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return nil, err
    }
    
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var users []*models.User
    if err = cursor.All(ctx, &users); err != nil {
        return nil, err
    }
    
    duration := time.Since(start)
    
    // Log query performance
    log.Printf("Query executed in %v, returned %d users, filter: %+v", 
        duration, len(users), filter)
    
    return users, nil
}
```

## Best Practices Summary

### 1. Connection Management
- Sử dụng single Manager instance
- Implement graceful shutdown
- Monitor connection health

### 2. Error Handling  
- Always check IsConnected() trước operations
- Implement retry logic
- Log connection events

### 3. Performance
- Configure connection pool appropriately
- Use indexes effectively
- Monitor query performance

### 4. Security
- Always use SSL/TLS trong production
- Store credentials securely
- Implement proper authentication

### 5. Testing
- Use mocks cho unit tests
- Implement integration tests
- Test error scenarios

Đây là hướng dẫn comprehensive để sử dụng MongoDB Provider hiệu quả trong ứng dụng Go của bạn.
