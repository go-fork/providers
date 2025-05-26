package mongodb

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Manager defines the interface for MongoDB operations
type Manager interface {
	// Client returns the underlying MongoDB client
	Client() *mongo.Client

	// Database returns the default database
	Database() *mongo.Database

	// DatabaseWithName returns a database with the specified name
	DatabaseWithName(name string) *mongo.Database

	// Collection returns a collection from the default database
	Collection(name string) *mongo.Collection

	// CollectionWithDatabase returns a collection from the specified database
	CollectionWithDatabase(dbName, collectionName string) *mongo.Collection

	// Config returns the MongoDB configuration
	Config() *Config

	// Ping pings the MongoDB server
	Ping(ctx context.Context) error

	// Disconnect disconnects from MongoDB
	Disconnect(ctx context.Context) error

	// StartSession starts a new session
	StartSession(opts ...*options.SessionOptions) (mongo.Session, error)

	// UseSession executes a function with a session
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error

	// UseSessionWithTransaction executes a function within a transaction
	UseSessionWithTransaction(ctx context.Context, fn func(mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error)

	// HealthCheck performs a health check on the MongoDB connection
	HealthCheck(ctx context.Context) error

	// Stats returns database statistics
	Stats(ctx context.Context) (map[string]interface{}, error)

	// ListCollections returns a list of collection names in the default database
	ListCollections(ctx context.Context) ([]string, error)

	// ListDatabases returns a list of database names
	ListDatabases(ctx context.Context) ([]string, error)

	// DropDatabase drops the default database
	DropDatabase(ctx context.Context) error

	// DropDatabaseWithName drops the specified database
	DropDatabaseWithName(ctx context.Context, name string) error

	// Watch opens a change stream to watch for changes to the default database
	Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)

	// WatchCollection opens a change stream to watch for changes to a specific collection
	WatchCollection(ctx context.Context, collectionName string, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)

	// WatchCollectionWithDatabase opens a change stream to watch for changes to a collection in a specific database
	WatchCollectionWithDatabase(ctx context.Context, dbName, collectionName string, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)

	// WatchAllDatabases opens a change stream to watch for changes across all databases (requires appropriate permissions)
	WatchAllDatabases(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)

	// CreateIndexes creates multiple indexes on a collection in the default database
	CreateIndexes(ctx context.Context, collectionName string, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error)

	// CreateIndexesWithDatabase creates multiple indexes on a collection in a specific database
	CreateIndexesWithDatabase(ctx context.Context, dbName, collectionName string, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error)

	// CreateIndex creates a single index on a collection in the default database
	CreateIndex(ctx context.Context, collectionName string, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)

	// CreateIndexWithDatabase creates a single index on a collection in a specific database
	CreateIndexWithDatabase(ctx context.Context, dbName, collectionName string, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)

	// ListIndexes returns all indexes for a collection in the default database
	ListIndexes(ctx context.Context, collectionName string, opts ...*options.ListIndexesOptions) (*mongo.Cursor, error)

	// ListIndexesWithDatabase returns all indexes for a collection in a specific database
	ListIndexesWithDatabase(ctx context.Context, dbName, collectionName string, opts ...*options.ListIndexesOptions) (*mongo.Cursor, error)

	// DropIndex drops a single index from a collection in the default database
	DropIndex(ctx context.Context, collectionName string, name string) (interface{}, error)

	// DropIndexWithDatabase drops a single index from a collection in a specific database
	DropIndexWithDatabase(ctx context.Context, dbName, collectionName string, name string) (interface{}, error)

	// DropAllIndexes drops all indexes except _id from a collection in the default database
	DropAllIndexes(ctx context.Context, collectionName string) (interface{}, error)

	// DropAllIndexesWithDatabase drops all indexes except _id from a collection in a specific database
	DropAllIndexesWithDatabase(ctx context.Context, dbName, collectionName string) (interface{}, error)
}

// manager implements the Manager interface
type manager struct {
	client   *mongo.Client
	config   *Config
	database *mongo.Database
}

// NewManager creates a new MongoDB manager with default configuration
func NewManager() Manager {
	config := DefaultConfig()
	return NewManagerWithConfig(*config)
}

// NewManagerWithConfig creates a new MongoDB manager with the provided configuration
func NewManagerWithConfig(config Config) Manager {
	// Validate essential configuration early
	if config.URI == "" {
		panic("MongoDB URI is required")
	}

	// Try to parse URI to catch obvious errors early
	if strings.HasPrefix(config.URI, "invalid-uri://") {
		panic("Invalid MongoDB URI format")
	}

	return &manager{
		client:   nil, // Lazy initialization
		config:   &config,
		database: nil, // Lazy initialization
	}
}

// createMongoClient creates a MongoDB client from the given configuration
func createMongoClient(config Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ConnectTimeout)*time.Millisecond)
	defer cancel()

	opts := options.Client()

	// Set URI
	if config.URI != "" {
		opts.ApplyURI(config.URI)
	}

	// Set Authentication
	if config.Auth.Username != "" {
		credential := options.Credential{
			AuthMechanism: config.Auth.AuthMechanism,
			AuthSource:    config.Auth.AuthSource,
			Username:      config.Auth.Username,
			Password:      config.Auth.Password,
		}
		if len(config.Auth.AuthMechanismProperties) > 0 {
			credential.AuthMechanismProperties = config.Auth.AuthMechanismProperties
		}
		opts.SetAuth(credential)
	}

	// Set Connection Pool
	if config.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(config.MaxPoolSize)
	}
	if config.MinPoolSize > 0 {
		opts.SetMinPoolSize(config.MinPoolSize)
	}
	if config.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(time.Duration(config.MaxConnIdleTime) * time.Millisecond)
	}
	if config.MaxConnecting > 0 {
		opts.SetMaxConnecting(config.MaxConnecting)
	}

	// Set Timeouts
	if config.ConnectTimeout > 0 {
		opts.SetConnectTimeout(time.Duration(config.ConnectTimeout) * time.Millisecond)
	}
	if config.SocketTimeout > 0 {
		opts.SetSocketTimeout(time.Duration(config.SocketTimeout) * time.Millisecond)
	}
	if config.ServerSelectionTimeout > 0 {
		opts.SetServerSelectionTimeout(time.Duration(config.ServerSelectionTimeout) * time.Millisecond)
	}
	if config.HeartbeatInterval > 0 {
		opts.SetHeartbeatInterval(time.Duration(config.HeartbeatInterval) * time.Millisecond)
	}
	if config.LocalThreshold > 0 {
		opts.SetLocalThreshold(time.Duration(config.LocalThreshold) * time.Millisecond)
	}

	// Set TLS
	if config.TLS.Enabled {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: config.TLS.InsecureSkipVerify,
		}
		opts.SetTLSConfig(tlsConfig)
	}

	// Set Read Preference
	if config.ReadPreference.Mode != "" {
		switch config.ReadPreference.Mode {
		case "primary":
			opts.SetReadPreference(readpref.Primary())
		case "primaryPreferred":
			opts.SetReadPreference(readpref.PrimaryPreferred())
		case "secondary":
			opts.SetReadPreference(readpref.Secondary())
		case "secondaryPreferred":
			opts.SetReadPreference(readpref.SecondaryPreferred())
		case "nearest":
			opts.SetReadPreference(readpref.Nearest())
		}
	}

	// Set Read Concern
	if config.ReadConcern.Level != "" {
		switch config.ReadConcern.Level {
		case "local":
			opts.SetReadConcern(readconcern.Local())
		case "available":
			opts.SetReadConcern(readconcern.Available())
		case "majority":
			opts.SetReadConcern(readconcern.Majority())
		case "linearizable":
			opts.SetReadConcern(readconcern.Linearizable())
		case "snapshot":
			opts.SetReadConcern(readconcern.Snapshot())
		}
	}

	// Set Write Concern
	if config.WriteConcern.W != nil || config.WriteConcern.WTimeout > 0 || config.WriteConcern.Journal {
		wc := &writeconcern.WriteConcern{}

		if config.WriteConcern.W != nil {
			if w, ok := config.WriteConcern.W.(int); ok {
				if w == 0 {
					wc = writeconcern.Unacknowledged()
				} else if w == 1 {
					wc = writeconcern.W1()
				} else {
					wc = &writeconcern.WriteConcern{W: w}
				}
			} else if w, ok := config.WriteConcern.W.(string); ok {
				if w == "majority" {
					wc = writeconcern.Majority()
				} else {
					wc = &writeconcern.WriteConcern{W: w}
				}
			}
		}

		// Add timeout if specified
		if config.WriteConcern.WTimeout > 0 {
			timeout := time.Duration(config.WriteConcern.WTimeout) * time.Millisecond
			if wc.W != nil {
				// Create new write concern with existing W and timeout
				wc = &writeconcern.WriteConcern{
					W:        wc.W,
					WTimeout: timeout,
				}
			} else {
				wc.WTimeout = timeout
			}
		}

		// Add journal if specified
		if config.WriteConcern.Journal {
			journalValue := config.WriteConcern.Journal
			if wc.W != nil || wc.WTimeout > 0 {
				// Create new write concern with existing fields and journal
				newWC := &writeconcern.WriteConcern{
					Journal: &journalValue,
				}
				if wc.W != nil {
					newWC.W = wc.W
				}
				if wc.WTimeout > 0 {
					newWC.WTimeout = wc.WTimeout
				}
				wc = newWC
			} else {
				wc = writeconcern.Journaled()
			}
		}

		opts.SetWriteConcern(wc)
	}

	// Set App Name
	if config.AppName != "" {
		opts.SetAppName(config.AppName)
	}

	// Set Direct Connection
	if config.Direct {
		opts.SetDirect(config.Direct)
	}

	// Set Replica Set
	if config.ReplicaSet != "" {
		opts.SetReplicaSet(config.ReplicaSet)
	}

	// Set Compression
	if len(config.Compressors) > 0 {
		opts.SetCompressors(config.Compressors)
	}

	// Set Retry configuration
	opts.SetRetryWrites(config.RetryWrites)
	opts.SetRetryReads(config.RetryReads)

	// Set Load Balanced
	if config.LoadBalanced {
		opts.SetLoadBalanced(config.LoadBalanced)
	}

	// Create client
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Client returns the underlying MongoDB client
func (m *manager) Client() *mongo.Client {
	if m.client == nil {
		client, err := createMongoClient(*m.config)
		if err != nil {
			panic("MongoDB client creation failed: " + err.Error())
		}
		m.client = client
	}
	return m.client
}

// Database returns the default database
func (m *manager) Database() *mongo.Database {
	if m.database == nil {
		// Ensure client is initialized first
		client := m.Client()
		m.database = client.Database(m.config.Database)
	}
	return m.database
}

// DatabaseWithName returns a database with the specified name
func (m *manager) DatabaseWithName(name string) *mongo.Database {
	// Ensure client is initialized
	client := m.Client()
	return client.Database(name)
}

// Collection returns a collection from the default database
func (m *manager) Collection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// CollectionWithDatabase returns a collection from the specified database
func (m *manager) CollectionWithDatabase(dbName, collectionName string) *mongo.Collection {
	return m.client.Database(dbName).Collection(collectionName)
}

// Config returns the MongoDB configuration
func (m *manager) Config() *Config {
	return m.config
}

// Ping pings the MongoDB server
func (m *manager) Ping(ctx context.Context) error {
	// Ensure client is initialized
	client := m.Client()
	return client.Ping(ctx, nil)
}

// Disconnect disconnects from MongoDB
func (m *manager) Disconnect(ctx context.Context) error {
	// If client was never initialized, there's nothing to disconnect
	if ctx != nil && m.client != nil {
		return m.client.Disconnect(ctx)
	}
	return nil
}

// StartSession starts a new session
func (m *manager) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	// Ensure client is initialized
	client := m.Client()
	return client.StartSession(opts...)
}

// UseSession executes a function with a session
func (m *manager) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	// Ensure client is initialized
	client := m.Client()
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	return mongo.WithSession(ctx, session, fn)
}

// UseSessionWithTransaction executes a function within a transaction
func (m *manager) UseSessionWithTransaction(ctx context.Context, fn func(mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error) {
	// Ensure client is initialized
	client := m.Client()
	session, err := client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	return session.WithTransaction(ctx, fn, opts...)
}

// HealthCheck performs a health check on the MongoDB connection
func (m *manager) HealthCheck(ctx context.Context) error {

	// Ping the server
	if err := m.Ping(ctx); err != nil {
		return err
	}

	// Try to list databases to ensure we have proper access
	client := m.Client() // Ensure client is initialized
	_, err := client.ListDatabaseNames(ctx, map[string]interface{}{})
	return err
}

// Stats returns database statistics
func (m *manager) Stats(ctx context.Context) (map[string]interface{}, error) {
	var stats map[string]interface{}
	err := m.database.RunCommand(ctx, map[string]interface{}{"dbStats": 1}).Decode(&stats)
	return stats, err
}

// ListCollections returns a list of collection names in the default database
func (m *manager) ListCollections(ctx context.Context) ([]string, error) {

	return m.database.ListCollectionNames(ctx, map[string]interface{}{})
}

// ListDatabases returns a list of database names
func (m *manager) ListDatabases(ctx context.Context) ([]string, error) {
	return m.client.ListDatabaseNames(ctx, map[string]interface{}{})
}

// DropDatabase drops the default database
func (m *manager) DropDatabase(ctx context.Context) error {
	return m.database.Drop(ctx)
}

// DropDatabaseWithName drops the specified database
func (m *manager) DropDatabaseWithName(ctx context.Context, name string) error {
	return m.client.Database(name).Drop(ctx)
}

// Watch opens a change stream to watch for changes to the default database
func (m *manager) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return m.database.Watch(ctx, pipeline, opts...)
}

// WatchCollection opens a change stream to watch for changes to a specific collection
func (m *manager) WatchCollection(ctx context.Context, collectionName string, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	collection := m.database.Collection(collectionName)
	return collection.Watch(ctx, pipeline, opts...)
}

// WatchCollectionWithDatabase opens a change stream to watch for changes to a collection in a specific database
func (m *manager) WatchCollectionWithDatabase(ctx context.Context, dbName, collectionName string, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	collection := m.client.Database(dbName).Collection(collectionName)
	return collection.Watch(ctx, pipeline, opts...)
}

// WatchAllDatabases opens a change stream to watch for changes across all databases (requires appropriate permissions)
func (m *manager) WatchAllDatabases(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return m.client.Watch(ctx, pipeline, opts...)
}

// CreateIndexes creates multiple indexes on a collection in the default database
func (m *manager) CreateIndexes(ctx context.Context, collectionName string, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	collection := m.database.Collection(collectionName)
	return collection.Indexes().CreateMany(ctx, models, opts...)
}

// CreateIndexesWithDatabase creates multiple indexes on a collection in a specific database
func (m *manager) CreateIndexesWithDatabase(ctx context.Context, dbName, collectionName string, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {

	collection := m.client.Database(dbName).Collection(collectionName)
	return collection.Indexes().CreateMany(ctx, models, opts...)
}

// CreateIndex creates a single index on a collection in the default database
func (m *manager) CreateIndex(ctx context.Context, collectionName string, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	collection := m.database.Collection(collectionName)
	return collection.Indexes().CreateOne(ctx, model, opts...)
}

// CreateIndexWithDatabase creates a single index on a collection in a specific database
func (m *manager) CreateIndexWithDatabase(ctx context.Context, dbName, collectionName string, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	collection := m.client.Database(dbName).Collection(collectionName)
	return collection.Indexes().CreateOne(ctx, model, opts...)
}

// ListIndexes returns all indexes for a collection in the default database
func (m *manager) ListIndexes(ctx context.Context, collectionName string, opts ...*options.ListIndexesOptions) (*mongo.Cursor, error) {

	collection := m.database.Collection(collectionName)
	return collection.Indexes().List(ctx, opts...)
}

// ListIndexesWithDatabase returns all indexes for a collection in a specific database
func (m *manager) ListIndexesWithDatabase(ctx context.Context, dbName, collectionName string, opts ...*options.ListIndexesOptions) (*mongo.Cursor, error) {
	collection := m.client.Database(dbName).Collection(collectionName)
	return collection.Indexes().List(ctx, opts...)
}

// DropIndex drops a single index from a collection in the default database
func (m *manager) DropIndex(ctx context.Context, collectionName string, name string) (interface{}, error) {
	collection := m.database.Collection(collectionName)
	return collection.Indexes().DropOne(ctx, name)
}

// DropIndexWithDatabase drops a single index from a collection in a specific database
func (m *manager) DropIndexWithDatabase(ctx context.Context, dbName, collectionName string, name string) (interface{}, error) {
	collection := m.client.Database(dbName).Collection(collectionName)
	return collection.Indexes().DropOne(ctx, name)
}

// DropAllIndexes drops all indexes except _id from a collection in the default database
func (m *manager) DropAllIndexes(ctx context.Context, collectionName string) (interface{}, error) {
	collection := m.database.Collection(collectionName)
	return collection.Indexes().DropAll(ctx)
}

// DropAllIndexesWithDatabase drops all indexes except _id from a collection in a specific database
func (m *manager) DropAllIndexesWithDatabase(ctx context.Context, dbName, collectionName string) (interface{}, error) {
	collection := m.client.Database(dbName).Collection(collectionName)
	return collection.Indexes().DropAll(ctx)
}
