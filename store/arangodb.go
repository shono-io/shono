package store

import "fmt"

import _ "github.com/shono-io/shono/benthos"

func NewArangodbStore(scopeCode, conceptCode, code, url, database, collection, username, password string) *ArangodbStore {
	return &ArangodbStore{
		store: &store{
			scopeCode:   scopeCode,
			conceptCode: conceptCode,
			code:        code,
			name:        fmt.Sprintf("%s Arangodb Store", code),
			description: fmt.Sprintf("%s items are stored within an arangodb store inside the %q database and %q collection", code, database, collection),
		},
		urls:       []string{url},
		username:   username,
		password:   password,
		database:   database,
		collection: collection,
	}
}

type ArangodbStore struct {
	*store

	urls       []string
	username   string
	password   string
	database   string
	collection string
}

func (s *ArangodbStore) AsBenthosComponent() (map[string]interface{}, error) {
	acfg := map[string]interface{}{
		"urls":       s.urls,
		"username":   s.username,
		"password":   s.password,
		"database":   s.database,
		"collection": s.collection,
	}

	return map[string]interface{}{
		"label":    s.code,
		"arangodb": acfg,
	}, nil

}
