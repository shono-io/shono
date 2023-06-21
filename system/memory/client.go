package memory

import (
	"context"
	"fmt"
	"github.com/shono-io/shono/inventory"
)

type inMemoryStore map[string]map[string]map[string]any

var Stores = map[string]inMemoryStore{}

func NewClient(cfg map[string]any) *Client {
	if _, ok := Stores[cfg[StorageIdField].(string)]; !ok {
		Stores[cfg[StorageIdField].(string)] = inMemoryStore{}
	}

	return &Client{
		storageId: cfg[StorageIdField].(string),
	}
}

type Client struct {
	storageId string
}

func (c *Client) List(ctx context.Context, collection string, filters map[string]any, paging *inventory.PagingOpts) (inventory.Cursor, error) {
	c.ensureCollection(collection)
	var results []map[string]any
	for _, item := range Stores[c.storageId][collection] {
		matches := true
		for k, v := range filters {
			if item[k] != v {
				matches = false
				break
			}
		}

		if matches {
			results = append(results, item)
		}
	}
	return &Cursor{items: results, currentIdx: 0}, nil
}

func (c *Client) Get(ctx context.Context, collection string, key string) (map[string]any, error) {
	c.ensureCollection(collection)
	return Stores[c.storageId][collection][key], nil
}

func (c *Client) Set(ctx context.Context, collection string, key string, value map[string]any) error {
	c.ensureCollection(collection)
	Stores[c.storageId][collection][key] = value
	return nil
}

func (c *Client) Add(ctx context.Context, collection string, key string, value map[string]any) error {
	c.ensureCollection(collection)

	if _, ok := Stores[c.storageId][collection][key]; ok {
		return fmt.Errorf("key already exists: %s", key)
	}

	Stores[c.storageId][collection][key] = value
	return nil
}

func (c *Client) Delete(ctx context.Context, collection string, key string) error {
	c.ensureCollection(collection)

	delete(Stores[c.storageId], key)

	return nil
}

func (c *Client) Close() error { return nil }

func (c *Client) ensureCollection(collection string) {
	if _, ok := Stores[c.storageId][collection]; !ok {
		Stores[c.storageId][collection] = map[string]map[string]any{}
	}
}

type Cursor struct {
	items      []map[string]any
	currentIdx int
}

func (c *Cursor) HasNext() bool {
	return c.currentIdx < len(c.items)
}

func (c *Cursor) Read() (map[string]any, error) {
	if !c.HasNext() {
		return nil, fmt.Errorf("no more items")
	}

	item := c.items[c.currentIdx]
	c.currentIdx++

	return item, nil
}

func (c *Cursor) Close() error {
	return nil
}
