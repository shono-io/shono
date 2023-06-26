package arangodb

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"strings"
)

func NewClient(cfg map[string]any) (inventory.StorageClient, error) {
	urls, ok := cfg[UrlsField]
	if !ok {
		return nil, fmt.Errorf("missing urls field")
	}

	username, ok := cfg[UsernameField]
	if !ok {
		return nil, fmt.Errorf("missing username field")
	}

	password, ok := cfg[PasswordField]
	if !ok {
		return nil, fmt.Errorf("missing password field")
	}

	database, ok := cfg[DatabaseField]
	if !ok {
		return nil, fmt.Errorf("missing database field")
	}

	var u []string
	for _, v := range urls.([]any) {
		u = append(u, v.(string))
	}

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: u,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create arangodb connection: %w", err)
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(username.(string), password.(string)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create arangodb client: %w", err)
	}

	db, err := c.Database(context.Background(), database.(string))
	if err != nil {
		return nil, fmt.Errorf("failed to get arangodb database: %w", err)
	}

	return Client{c, db}, nil
}

type Client struct {
	cl driver.Client
	db driver.Database
}

func (c Client) List(ctx context.Context, collection string, filters map[string]any, paging *inventory.PagingOpts) (inventory.Cursor, error) {
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

func (c Client) Get(ctx context.Context, collection string, key string) (map[string]any, error) {
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

func (c Client) Set(ctx context.Context, collection string, key string, value map[string]any) error {
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

func (c Client) Add(ctx context.Context, collection string, key string, value map[string]any) error {
	col, err := c.db.Collection(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}

	// -- override the key
	value["_key"] = key

	_, err = col.CreateDocument(ctx, value)
	return err
}

func (c Client) Delete(ctx context.Context, collection string, key string) error {
	col, err := c.db.Collection(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}

	_, err = col.RemoveDocument(ctx, key)
	return err
}

func (c Client) Close() error {
	return nil
}

func (c Client) collectionForKey(ctx context.Context, key commons.Reference) (driver.Collection, error) {
	return c.db.Collection(ctx, key.Kind())
}

func (c Client) buildQuery(collection string, filters map[string]any, paging *inventory.PagingOpts) (string, map[string]any, error) {
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
