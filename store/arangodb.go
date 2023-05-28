package store

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"strings"
)

type Arangodb[T State] struct {
	db  driver.Database
	col driver.Collection
}

func (s *Arangodb[T]) List(ctx context.Context, filters map[string]interface{}, offset uint, size uint) ([]T, int64, error) {
	params := map[string]any{
		"col":    s.col.Name(),
		"offset": offset,
		"size":   size,
	}

	var f []string
	i := 0
	for k, v := range filters {
		i++
		f = append(f, fmt.Sprintf("%s == @filter%d", k, i))
		params[fmt.Sprintf("filter%d", i)] = v
	}

	q := fmt.Sprintf("FOR d IN @col FILTER %s LIMIT @offset, @size RETURN d", strings.Join(f, " AND "))

	cur, err := s.db.Query(ctx, q, params)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close()

	var result []T
	if cur == nil {
		return result, 0, nil
	}

	for cur.HasMore() {
		var doc T
		if _, err := cur.ReadDocument(ctx, &doc); err != nil {
			return nil, 0, err
		}

		result = append(result, doc)
	}

	return result, cur.Count(), nil
}

func (s *Arangodb[T]) Get(ctx context.Context, key string) (*T, error) {
	var state T
	_, err := s.col.ReadDocument(ctx, key, &state)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			return nil, nil
		}

		return nil, err
	}

	return &state, nil
}

func (s *Arangodb[T]) Persist(ctx context.Context, state T, mode PersistMode) error {
	// TODO eventually we need to take care of the document meta to make sure we don't run into concurrency issues
	switch mode {
	case CreatePersistMode:
		_, err := s.col.CreateDocument(ctx, state)
		if err != nil {
			return err
		}
	case ReplacePersistMode:
		_, err := s.col.ReplaceDocument(ctx, state.Key(), state)
		if err != nil {
			return err
		}
	case PatchPersistMode:
		_, err := s.col.UpdateDocument(ctx, state.Key(), state)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Arangodb[T]) Remove(ctx context.Context, key string) error {
	_, err := s.col.RemoveDocument(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
