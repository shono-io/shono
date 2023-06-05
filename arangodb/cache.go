package arangodb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/benthosdev/benthos/v4/public/service"
	"time"
)

func cacheFromConfig(conf *service.ParsedConfig) (*arangodbCache, error) {
	database, err := conf.FieldString("database")
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	collection, err := conf.FieldString("collection")
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	adb, err := getClientFromConfig(conf)
	if err != nil {
		return nil, err
	}

	// -- get the database
	db, err := adb.Database(context.Background(), database)
	if err != nil {
		return nil, fmt.Errorf("failed to get arangodb database: %w", err)
	}

	// -- get the collection
	col, err := db.Collection(context.Background(), collection)
	if err != nil {
		return nil, fmt.Errorf("failed to get arangodb collection: %w", err)
	}

	return &arangodbCache{
		collection: col,
	}, nil
}

func arangodbCacheConfig() *service.ConfigSpec {
	spec := service.NewConfigSpec().
		Stable()

	for _, f := range clientFields() {
		spec = spec.Field(f)
	}

	return spec.
		Field(service.NewStringField("database").
			Default("_system").
			Description("The database to use when connecting to the ArangoDB server")).
		Field(service.NewStringField("collection").
			Description("The collection to use when connecting to the ArangoDB server"))
}

type arangodbCache struct {
	collection driver.Collection
}

func (a *arangodbCache) Close(ctx context.Context) error {
	return nil
}

func (a *arangodbCache) Get(ctx context.Context, key string) ([]byte, error) {
	var result map[string]any

	if _, err := a.collection.ReadDocument(ctx, key, &result); err != nil {
		if driver.IsNotFoundGeneral(err) {
			return nil, nil
		}

		return nil, err
	}

	return json.Marshal(result)
}

func (a *arangodbCache) Set(ctx context.Context, key string, value []byte, ttl *time.Duration) error {
	var doc map[string]any
	if err := json.Unmarshal(value, &doc); err != nil {
		return err
	}

	fnd, err := a.collection.DocumentExists(ctx, key)
	if err != nil {
		return err
	}
	if fnd {
		_, err := a.collection.ReplaceDocument(ctx, key, doc)
		return err
	} else {
		_, err := a.collection.CreateDocument(ctx, doc)
		return err
	}
}

func (a *arangodbCache) Add(ctx context.Context, key string, value []byte, ttl *time.Duration) error {
	var doc map[string]any
	if err := json.Unmarshal(value, &doc); err != nil {
		return err
	}

	_, err := a.collection.CreateDocument(ctx, doc)
	return err
}

func (a *arangodbCache) Delete(ctx context.Context, key string) error {
	_, err := a.collection.RemoveDocument(ctx, key)
	return err
}
