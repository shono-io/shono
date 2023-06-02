package benthos

import (
	"fmt"
	"github.com/shono-io/shono"
	"github.com/shono-io/shono/store"
)

func (g *Generator) generateCaches(result map[string]any, scope shono.Scope) (err error) {
	// -- make a list of all stores used by the reaktors
	stores := map[string]shono.Store{}
	for _, reaktor := range scope.Reaktors() {
		for _, s := range reaktor.Stores() {
			stores[s.Key().String()] = s
		}
	}

	var cr []map[string]any

	// -- convert each of these stores in their yaml
	for _, s := range stores {
		st, err := storeToCache(s)
		if err != nil {
			return fmt.Errorf("failed to convert store %q to cache: %w", s.Key().String(), err)
		}

		cr = append(cr, st)
	}

	result["cache_resources"] = cr

	return nil
}

func storeToCache(s shono.Store) (map[string]any, error) {
	switch st := s.(type) {
	case *store.ArangodbStore:
		return arangodbToCache(st)
	default:
		return nil, fmt.Errorf("unknown store type %T", s)
	}
}

func arangodbToCache(s *store.ArangodbStore) (map[string]any, error) {
	username, passwd := s.Credentials()

	return map[string]interface{}{
		"label": labelize(s.Key().CodeString()),
		"arangodb": map[string]interface{}{
			"urls":       s.Urls(),
			"username":   username,
			"password":   passwd,
			"database":   s.Database(),
			"collection": s.Collection(),
		},
	}, nil
}
