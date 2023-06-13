package storage

import (
	"context"
	"fmt"
	"github.com/shono-io/shono/graph"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbConfig struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type mongodb struct {
	key    string
	config MongodbConfig
}

func (s mongodb) GetClient() (graph.StorageClient, error) {
	opts := options.Client().
		ApplyURI(s.config.Uri).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	cl, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	db := cl.Database(s.config.Database)

	return client{cl, db}, nil
}

func (s mongodb) Key() string {
	return s.key
}

type client struct {
	cl *mongo.Client
	db *mongo.Database
}

func (c client) List(ctx context.Context, collection string, filters map[string]any, paging *graph.PagingOpts) (graph.Cursor, error) {
	col := c.db.Collection(collection)

	var fs []bson.D
	for k, v := range filters {
		fs = append(fs, bson.D{{fmt.Sprintf("data.%s", k), bson.D{{"$eq", v}}}})
	}

	filter := bson.D{
		{"$and", bson.A{fs}},
	}

	var opts []*options.FindOptions
	if paging != nil {
		opts = append(opts,
			options.Find().SetSkip(paging.Offset),
			options.Find().SetLimit(paging.Size),
		)

	}

	cur, err := col.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	return &cursorWrapper{ctx, cur}, nil
}

func (c client) Get(ctx context.Context, collection string, key string) (map[string]any, error) {
	col := c.db.Collection(collection)
	r := col.FindOne(ctx, key)
	if r == nil {
		return nil, nil
	}
	if r.Err() != nil {
		return nil, r.Err()
	}

	var d document
	err := r.Decode(&d)
	if err != nil {
		return nil, err
	}

	return d.Data, nil
}

func (c client) Set(ctx context.Context, collection string, key string, value map[string]any) error {
	col := c.db.Collection(collection)
	_, err := col.InsertOne(ctx, document{
		Key:  key,
		Data: value,
	})
	return err
}

func (c client) Add(ctx context.Context, collection string, key string, value map[string]any) error {
	col := c.db.Collection(collection)
	_, err := col.InsertOne(ctx, document{
		Key:  key,
		Data: value,
	})
	return err
}

func (c client) Delete(ctx context.Context, collection string, key string) error {
	filter := bson.D{{"_id", bson.D{{"$eq", key}}}}

	col := c.db.Collection(collection)
	r, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if r.DeletedCount != 1 {
		return fmt.Errorf("failed to delete document")
	}

	return nil
}

func (c client) Close() error {
	return c.cl.Disconnect(context.Background())
}

type document struct {
	Key  string         `json:"_id"`
	Data map[string]any `json:"data"`
}

type cursorWrapper struct {
	ctx    context.Context
	cursor *mongo.Cursor
}

func (c *cursorWrapper) HasNext() bool {
	return c.cursor.Next(c.ctx)
}

func (c *cursorWrapper) Read() (map[string]any, error) {
	var d document
	if err := c.cursor.Decode(&d); err != nil {
		return nil, err
	}

	return d.Data, nil
}

func (c *cursorWrapper) Close() error {
	return c.cursor.Close(c.ctx)
}
