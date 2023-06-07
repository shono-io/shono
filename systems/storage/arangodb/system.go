package arangodb

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
	"strings"
)

type system struct {
}

func (s system) GetClient(config map[string]any) (graph.StorageClient, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: config["urls"].([]string),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create arangodb connection: %w", err)
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(config["username"].(string), config["password"].(string)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create arangodb client: %w", err)
	}

	db, err := c.Database(context.Background(), config["database"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to get arangodb database: %w", err)
	}

	return client{c, db}, nil
}

type client struct {
	cl driver.Client
	db driver.Database
}

func (c client) List(ctx context.Context, collection string, filters map[string]any, paging *graph.PagingOpts) (graph.Cursor, error) {
	q, bv, err := c.buildQuery(collection, filters, paging)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	cursor, err := c.db.Query(ctx, q, bv)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &cursorWrapper{cursor, ctx}, nil
}

func (c client) Get(ctx context.Context, collection string, key string) (map[string]any, error) {
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

func (c client) Set(ctx context.Context, collection string, key string, value map[string]any) error {
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

func (c client) Add(ctx context.Context, collection string, key string, value map[string]any) error {
	col, err := c.db.Collection(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}

	// -- override the key
	value["_key"] = key

	_, err = col.CreateDocument(ctx, value)
	return err
}

func (c client) Delete(ctx context.Context, collection string, key string) error {
	col, err := c.db.Collection(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}

	_, err = col.RemoveDocument(ctx, key)
	return err
}

func (c client) Close() error {
	return nil
}

func (c client) collectionForKey(ctx context.Context, key commons.Key) (driver.Collection, error) {
	return c.db.Collection(ctx, key.Kind())
}

func (c client) buildQuery(collection string, filters map[string]any, paging *graph.PagingOpts) (string, map[string]any, error) {
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

type cursorWrapper struct {
	c   driver.Cursor
	ctx context.Context
}

func (c *cursorWrapper) Close() error {
	return c.c.Close()
}

func (c *cursorWrapper) Count() int64 {
	return c.c.Count()
}

func (c *cursorWrapper) HasNext() bool {
	return c.c.HasMore()
}

func (c *cursorWrapper) Read() (map[string]any, error) {
	var target map[string]any
	_, err := c.c.ReadDocument(c.ctx, target)
	return target, err
}
