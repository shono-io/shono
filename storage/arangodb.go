package storage

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/shono-io/shono/core"
	"github.com/shono-io/shono/graph"
	"strings"
)

type ArangodbConfig struct {
	Urls     []string `json:"urls"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Database string   `json:"database"`
}

type arangodb struct {
	key    string
	config ArangodbConfig
}

func (s arangodb) GetClient() (graph.StorageClient, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: s.config.Urls,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create arangodb connection: %w", err)
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(s.config.Username, s.config.Password),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create arangodb client: %w", err)
	}

	db, err := c.Database(context.Background(), s.config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to get arangodb database: %w", err)
	}

	return arangodbClient{c, db}, nil
}

func (s arangodb) Key() string {
	return s.key
}

type arangodbClient struct {
	cl driver.Client
	db driver.Database
}

func (c arangodbClient) List(ctx context.Context, collection string, filters map[string]any, paging *graph.PagingOpts) (graph.Cursor, error) {
	q, bv, err := c.buildQuery(collection, filters, paging)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	cursor, err := c.db.Query(ctx, q, bv)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &arangodbCursorWrapper{cursor, ctx}, nil
}

func (c arangodbClient) Get(ctx context.Context, collection string, key string) (map[string]any, error) {
	col, err := c.db.Collection(ctx, collection)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	var target map[string]any
	_, err = col.ReadDocument(ctx, key, &target)
	return target, err
}

func (c arangodbClient) Set(ctx context.Context, collection string, key string, value map[string]any) error {
	col, err := c.db.Collection(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}

	fnd, err := col.DocumentExists(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to check if document exists: %w", err)
	}

	// -- override the key
	value["_key"] = key

	if fnd {
		_, err = col.ReplaceDocument(ctx, key, value)
		return err
	} else {
		_, err = col.CreateDocument(ctx, value)
		return err
	}
}

func (c arangodbClient) Add(ctx context.Context, collection string, key string, value map[string]any) error {
	col, err := c.db.Collection(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}

	// -- override the key
	value["_key"] = key

	_, err = col.CreateDocument(ctx, value)
	return err
}

func (c arangodbClient) Delete(ctx context.Context, collection string, key string) error {
	col, err := c.db.Collection(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}

	_, err = col.RemoveDocument(ctx, key)
	return err
}

func (c arangodbClient) Close() error {
	return nil
}

func (c arangodbClient) collectionForKey(ctx context.Context, key core.Reference) (driver.Collection, error) {
	return c.db.Collection(ctx, key.Kind())
}

func (c arangodbClient) buildQuery(collection string, filters map[string]any, paging *graph.PagingOpts) (string, map[string]any, error) {
	result := fmt.Sprintf("FOR d IN %s", collection)

	var parts []string

	var fl []string
	vars := map[string]any{}
	for k, v := range filters {

		ke := strings.ReplaceAll(k, ".", "_")
		fl = append(fl, fmt.Sprintf("d.%s == @%s", k, ke))
		vars[ke] = v
	}

	if len(fl) > 0 {
		parts = append(parts, fmt.Sprintf("FILTER %s", strings.Join(fl, " AND ")))
	}

	if paging != nil {
		parts = append(parts, fmt.Sprintf("LIMIT %d, %d", paging.Offset, paging.Size))
	}

	result += " RETURN d"

	return result, vars, nil
}

type arangodbCursorWrapper struct {
	c   driver.Cursor
	ctx context.Context
}

func (c *arangodbCursorWrapper) Close() error {
	return c.c.Close()
}

func (c *arangodbCursorWrapper) Count() int64 {
	return c.c.Count()
}

func (c *arangodbCursorWrapper) HasNext() bool {
	return c.c.HasMore()
}

func (c *arangodbCursorWrapper) Read() (map[string]any, error) {
	var target map[string]any
	_, err := c.c.ReadDocument(c.ctx, target)
	return target, err
}
