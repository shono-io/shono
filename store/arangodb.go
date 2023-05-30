package store

import "fmt"

func NewArangodbStore(conceptCode, code, url, database, collection, username, password string) *ArangodbStore {
	return &ArangodbStore{
		store: &store{
			conceptCode: conceptCode,
			code:        code,
			name:        fmt.Sprintf("%s Arangodb Store", code),
			description: fmt.Sprintf("%s items are stored within an arangodb store inside the %q database and %q collection", code, database, collection),
		},
		url:        url,
		username:   username,
		password:   password,
		database:   database,
		collection: collection,
	}
}

type ArangodbStore struct {
	*store

	url        string
	username   string
	password   string
	database   string
	collection string
}

func (s *ArangodbStore) AsBenthosComponent() (map[string]interface{}, error) {
	acfg := map[string]interface{}{
		"url":        s.url,
		"username":   s.username,
		"password":   s.password,
		"database":   s.database,
		"collection": s.collection,
	}

	return map[string]interface{}{
		"arangodb": acfg,
	}, nil

}
