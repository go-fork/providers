package driver

import (
	"context"
	"time"

	"go.fork.vn/providers/cache/config"
	"go.fork.vn/providers/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCacheItem đại diện cho một mục trong MongoDB cache.
//
// Cấu trúc này lưu trữ dữ liệu cache dưới dạng document trong MongoDB,
// với các trường cần thiết như key, value, thời gian hết hạn và thời gian tạo.
type MongoCacheItem struct {
	Key        string      `bson:"_id"`        // Cache key, sử dụng như primary key
	Value      interface{} `bson:"value"`      // Giá trị được lưu trong cache
	Expiration int64       `bson:"expiration"` // Thời điểm hết hạn (UnixNano), 0 nếu không hết hạn
	CreatedAt  time.Time   `bson:"created_at"` // Thời điểm tạo cache item
}

type MongoDBDriver interface {
	Driver
	// ensureIndexes tạo các index cần thiết cho MongoDB collection.
	ensureIndexes(ctx context.Context) error
}

// mongoDBDriver cài đặt cache driver sử dụng MongoDB.
//
// mongoDBDriver lưu trữ dữ liệu cache trong một collection MongoDB.
// Driver này phù hợp cho các ứng dụng yêu cầu persistence, phân tán dữ liệu
// và cần tìm kiếm trong dữ liệu cache. MongoDB TTL index được sử dụng để tự động
// xóa các document đã hết hạn.
type mongoDBDriver struct {
	mongodb    *mongodb.Manager // Service Provider mongoDB manager
	config     config.DriverMongodbConfig
	database   *mongo.Database   // MongoDB database để lưu trữ cache
	collection *mongo.Collection // MongoDB collection để lưu trữ cache
}

// NewMongoDBDriver tạo một MongoDB driver mới với cấu hình mặc định.
//
// Phương thức này khởi tạo một MongoDBDriver mới với thông tin kết nối cơ bản.
// Nó sử dụng giá trị mặc định cho defaultExpiration là 5 phút.
//
// Params:
//   - config: config.DriverMongodbConfig,
//   - manager: mongodb.Manager
//
// Returns:
//   - *MongoDBDriver: Driver đã được khởi tạo
//   - error: Lỗi nếu không thể kết nối đến MongoDB hoặc tạo indices
func NewMongoDBDriver(cfg config.DriverMongodbConfig, manager mongodb.Manager) (MongoDBDriver, error) {
	driver := &mongoDBDriver{
		mongodb:    &manager,
		config:     cfg,
		database:   manager.DatabaseWithName(cfg.Database),
		collection: manager.DatabaseWithName(cfg.Database).Collection(cfg.Collection),
	}

	// Tạo indices cần thiết
	if err := driver.ensureIndexes(context.Background()); err != nil {
		return nil, err
	}

	return driver, nil
}

// ensureIndexes tạo các index cần thiết cho MongoDB collection.
//
// Phương thức này tạo TTL index trên trường expiration để MongoDB tự động
// xóa các document đã hết hạn. TTL index kiểm tra các document có expiration > 0
// và xóa chúng khi thời gian hiện tại vượt quá giá trị expiration.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//
// Returns:
//   - error: Lỗi nếu có trong quá trình tạo index
func (d *mongoDBDriver) ensureIndexes(ctx context.Context) error {
	// Tạo TTL index trên trường expiration
	// Index này sẽ tự động xóa documents khi expiration time đến
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "expiration", Value: 1}, // Index trên trường expiration
		},
		Options: options.Index().
			SetExpireAfterSeconds(0). // TTL index với 0 seconds để sử dụng giá trị trong document
			SetPartialFilterExpression(bson.D{
				{Key: "expiration", Value: bson.D{{Key: "$gt", Value: 0}}}, // Chỉ áp dụng cho document có expiration > 0
			}).
			SetName("cache_expiration_ttl"), // Tên index
	}

	// Tạo index
	_, err := d.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	// Tạo index cho _id nếu chưa có (thường MongoDB tự tạo)
	// Nhưng đảm bảo performance cho cache key lookups
	keyIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "_id", Value: 1},
		},
		Options: options.Index().
			SetName("cache_key_index"),
	}

	_, err = d.collection.Indexes().CreateOne(ctx, keyIndexModel)
	if err != nil {
		// Ignore error nếu index đã tồn tại
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}
	return nil
}

// Get lấy một giá trị từ cache.
//
// Phương thức này tìm kiếm document theo key trong MongoDB collection
// và giải mã dữ liệu. Nếu document không tồn tại hoặc đã hết hạn, phương thức
// trả về false ở giá trị thứ hai và cập nhật bộ đếm miss. Nếu tìm thấy và
// còn hạn, phương thức trả về giá trị và cập nhật bộ đếm hit.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần tìm
//
// Returns:
//   - interface{}: Giá trị được lưu trong cache (nil nếu không tìm thấy)
//   - bool: true nếu tìm thấy key và chưa hết hạn, false nếu ngược lại
func (d *mongoDBDriver) Get(ctx context.Context, key string) (interface{}, bool) {
	var cacheItem MongoCacheItem
	err := d.collection.FindOne(ctx, bson.M{"_id": key}).Decode(&cacheItem)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			d.config.Misses++
			return nil, false
		}
		return nil, false
	}

	// Kiểm tra expiration (TTL index sẽ tự động xóa expired documents,
	// nhưng chúng ta vẫn kiểm tra để đảm bảo tính nhất quán)
	if cacheItem.Expiration > 0 && time.Now().UnixNano() > cacheItem.Expiration {
		d.config.Misses++
		// TTL index sẽ tự động xóa, không cần xóa thủ công
		return nil, false
	}

	d.config.Hits++
	return cacheItem.Value, true
}

// Set đặt một giá trị vào cache với TTL tùy chọn.
//
// Phương thức này tạo hoặc cập nhật một document trong MongoDB collection
// để lưu trữ cặp key-value với thời gian sống được chỉ định. Nếu key đã tồn tại,
// document cũ sẽ bị thay thế.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key để lưu giá trị
//   - value: Giá trị cần lưu trữ
//   - ttl: Thời gian sống của giá trị (0 để sử dụng mặc định, -1 để không hết hạn)
//
// Returns:
//   - error: Lỗi nếu có trong quá trình lưu trữ vào MongoDB
func (d *mongoDBDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var exp int64
	now := time.Now()

	if ttl == 0 {
		if d.config.GetDefaultExpiration() > 0 {
			exp = now.Add(d.config.GetDefaultExpiration()).UnixNano()
		}
	} else if ttl > 0 {
		exp = now.Add(ttl).UnixNano()
	}

	// Tạo cache item
	cacheItem := MongoCacheItem{
		Key:        key,
		Value:      value,
		Expiration: exp,
		CreatedAt:  now,
	}

	// Nếu có expiration > 0, đặt thời gian hết hạn
	opts := options.ReplaceOptions{}
	opts.SetUpsert(true)

	// Lưu vào MongoDB
	_, err := d.collection.ReplaceOne(
		ctx,
		bson.M{"_id": key},
		cacheItem,
		&opts,
	)

	return err
}

// Has kiểm tra xem một key có tồn tại trong cache không.
//
// Phương thức này xác định liệu một key có tồn tại trong cache và chưa hết hạn hay không.
// Nó sử dụng phương thức Get để thực hiện kiểm tra và do đó cũng cập nhật bộ đếm hit/miss.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần kiểm tra
//
// Returns:
//   - bool: true nếu key tồn tại và chưa hết hạn, false nếu ngược lại
func (d *mongoDBDriver) Has(ctx context.Context, key string) bool {
	_, exists := d.Get(ctx, key)
	return exists
}

// Delete xóa một key khỏi cache.
//
// Phương thức này xóa document tương ứng với key được chỉ định khỏi MongoDB collection.
// Nếu key không tồn tại, thao tác này không có tác dụng và không trả về lỗi.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần xóa
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa
func (d *mongoDBDriver) Delete(ctx context.Context, key string) error {
	_, err := d.collection.DeleteOne(ctx, bson.M{"_id": key})
	return err
}

// Flush xóa tất cả các key khỏi cache.
//
// Phương thức này xóa tất cả documents trong MongoDB collection,
// làm trống hoàn toàn bộ nhớ cache.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa
func (d *mongoDBDriver) Flush(ctx context.Context) error {
	_, err := d.collection.DeleteMany(ctx, bson.M{})
	return err
}

// GetMultiple lấy nhiều giá trị từ cache.
//
// Phương thức này tìm kiếm và trả về nhiều giá trị từ cache dựa trên danh sách key.
// Các key không tìm thấy hoặc đã hết hạn sẽ được thêm vào danh sách missed.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - keys: Danh sách các khóa cần lấy
//
// Returns:
//   - map[string]interface{}: Map chứa các key tìm thấy và giá trị tương ứng
//   - []string: Danh sách các key không tìm thấy hoặc đã hết hạn
func (d *mongoDBDriver) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, []string) {
	results := make(map[string]interface{})
	missed := make([]string, 0)

	// Tạo filter cho nhiều key
	filter := bson.M{"_id": bson.M{"$in": keys}}

	// Tìm tất cả các document khớp với filter
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return results, keys
	}
	defer cursor.Close(ctx)

	// Tạo map để theo dõi các key đã tìm thấy
	found := make(map[string]bool)
	now := time.Now().UnixNano()

	// Giải mã các kết quả
	for cursor.Next(ctx) {
		var cacheItem MongoCacheItem
		if err := cursor.Decode(&cacheItem); err != nil {
			continue
		}

		// Kiểm tra expiration
		if cacheItem.Expiration > 0 && now > cacheItem.Expiration {
			missed = append(missed, cacheItem.Key)
			continue
		}

		results[cacheItem.Key] = cacheItem.Value
		found[cacheItem.Key] = true
	}

	// Thêm các key không tìm thấy vào danh sách missed
	for _, key := range keys {
		if !found[key] {
			missed = append(missed, key)
		}
	}

	return results, missed
}

// SetMultiple đặt nhiều giá trị vào cache.
//
// Phương thức này lưu trữ nhiều cặp key-value vào MongoDB collection trong một lần gọi
// với cùng một thời gian sống. Nó sử dụng BulkWrite để tối ưu hiệu suất.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - values: Map chứa các key và giá trị tương ứng cần lưu trữ
//   - ttl: Thời gian sống chung cho tất cả các giá trị
//
// Returns:
//   - error: Lỗi nếu có trong quá trình lưu trữ
func (d *mongoDBDriver) SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = d.config.GetDefaultExpiration()
	}
	// Chuẩn bị các document để chèn
	now := time.Now()
	exp := now.Add(ttl).UnixNano()

	var operations []mongo.WriteModel

	for key, value := range values {
		cacheItem := MongoCacheItem{
			Key:        key,
			Value:      value,
			Expiration: exp,
			CreatedAt:  now,
		}

		operation := mongo.NewReplaceOneModel().
			SetFilter(bson.M{"_id": key}).
			SetReplacement(cacheItem).
			SetUpsert(true)

		operations = append(operations, operation)
	}

	// Thực hiện bulk write
	_, err := d.collection.BulkWrite(ctx, operations)
	return err
}

// DeleteMultiple xóa nhiều key khỏi cache.
//
// Phương thức này xóa nhiều key và giá trị tương ứng khỏi MongoDB collection trong một lần gọi.
// Nó sử dụng một single query với $in operator để tối ưu hiệu suất.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - keys: Danh sách các khóa cần xóa
//
// Returns:
//   - error: Lỗi nếu có trong quá trình xóa
func (d *mongoDBDriver) DeleteMultiple(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	// Xóa tất cả các document với key trong danh sách
	_, err := d.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": keys}})
	return err
}

// Remember lấy một giá trị từ cache hoặc thực thi callback nếu không tìm thấy.
//
// Phương thức này kiểm tra xem một key có tồn tại trong cache không, nếu có thì
// trả về giá trị tương ứng. Nếu key không tồn tại hoặc đã hết hạn, phương thức
// sẽ gọi hàm callback để lấy dữ liệu, lưu kết quả vào cache và trả về giá trị đó.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//   - key: Cache key cần tìm hoặc lưu vào cache
//   - ttl: Thời gian sống của giá trị nếu phải lấy từ callback
//   - callback: Hàm được gọi để lấy dữ liệu khi key không có trong cache
//
// Returns:
//   - interface{}: Giá trị từ cache hoặc từ callback
//   - error: Lỗi nếu có trong quá trình thực hiện hoặc từ callback
func (d *mongoDBDriver) Remember(ctx context.Context, key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// Kiểm tra cache trước
	value, found := d.Get(ctx, key)
	if found {
		return value, nil
	}

	// Không tìm thấy, gọi callback
	value, err := callback()
	if err != nil {
		return nil, err
	}

	// Lưu kết quả vào cache
	err = d.Set(ctx, key, value, ttl)
	return value, err
}

// Stats trả về thông tin thống kê về cache.
//
// Phương thức này thu thập và trả về các thông tin thống kê về trạng thái
// hiện tại của MongoDB cache như số lượng document, kích thước, số lần hit/miss, v.v.
//
// Params:
//   - ctx: Context để kiểm soát thời gian thực thi của thao tác
//
// Returns:
//   - map[string]interface{}: Map chứa các thông tin thống kê
func (d *mongoDBDriver) Stats(ctx context.Context) map[string]interface{} {
	// Đếm số lượng document
	count, err := d.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		count = -1
	}

	// Lấy stats từ cơ sở dữ liệu
	var stats bson.M
	cmd := bson.D{{Key: "collStats", Value: d.collection.Name()}}
	err = d.database.RunCommand(ctx, cmd).Decode(&stats)
	if err != nil {
		stats = bson.M{}
	}

	return map[string]interface{}{
		"count":  count,
		"hits":   d.config.Hits,
		"misses": d.config.Misses,
		"type":   "mongodb",
		"stats":  stats,
	}
}

// Close giải phóng tài nguyên của driver.
//
// Phương thức này đóng kết nối tới MongoDB và giải phóng
// các tài nguyên khác được sử dụng bởi driver.
//
// Returns:
//   - error: Lỗi nếu có trong quá trình đóng kết nối
func (d *mongoDBDriver) Close() error {
	// disconnect MongoDB connection by service provider mongodb
	return nil
}
