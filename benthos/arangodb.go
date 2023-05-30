package benthos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/sirupsen/logrus"
	"time"
)

var arangodbSpec = service.NewConfigSpec().
	Field(service.NewStringListField("urls").Default("http://localhost:8529").Description("The URLs of the ArangoDB server")).
	Field(service.NewStringField("username").Default("root").Description("The username to use when connecting to the ArangoDB server")).
	Field(service.NewStringField("password").Default("").Description("The password to use when connecting to the ArangoDB server").Secret()).
	Field(service.NewStringField("database").Default("_system").Description("The database to use when connecting to the ArangoDB server")).
	Field(service.NewStringField("collection").Description("The collection to use when connecting to the ArangoDB server"))

func init() {
	err := service.RegisterCache("arangodb", arangodbSpec, func(conf *service.ParsedConfig, mgr *service.Resources) (service.Cache, error) {
		url, err := conf.FieldStringList("urls")
		if err != nil {
			return nil, fmt.Errorf("failed to get url: %w", err)
		}

		username, err := conf.FieldString("username")
		if err != nil {
			return nil, fmt.Errorf("failed to get username: %w", err)
		}

		password, err := conf.FieldString("password")
		if err != nil {
			return nil, fmt.Errorf("failed to get password: %w", err)
		}

		database, err := conf.FieldString("database")
		if err != nil {
			return nil, fmt.Errorf("failed to get database: %w", err)
		}

		collection, err := conf.FieldString("collection")
		if err != nil {
			return nil, fmt.Errorf("failed to get collection: %w", err)
		}

		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: url,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create arangodb connection: %w", err)
		}

		adb, err := driver.NewClient(driver.ClientConfig{
			Connection:     conn,
			Authentication: driver.BasicAuthentication(username, password),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create arangodb client: %w", err)
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
	})

	if err != nil {
		logrus.Panic(err)
	}
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
