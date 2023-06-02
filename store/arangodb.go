package store

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/shono-io/shono"
	"net/http"
	"strings"
)

func NewArangodbStore(concept shono.Concept, code, url, database, collection, username, password string) *ArangodbStore {
	return &ArangodbStore{
		store: store{
			concept:     concept,
			key:         concept.Key().Child("store", code),
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
	store

	urls       []string
	username   string
	password   string
	database   string
	collection string
}

func (s *ArangodbStore) Operations() map[shono.StoreOpertationId]shono.StoreOperation {
	return map[shono.StoreOpertationId]shono.StoreOperation{
		shono.ExistsOperation: newStoreOperation(shono.ExistsOperation,
			fmt.Sprintf("%sExists", strings.ToTitle(s.concept.Key().Code())),
			fmt.Sprintf("%s Exists", strings.ToTitle(s.concept.Key().Code())),
			fmt.Sprintf("Checks if %s exists", strings.ToTitle(s.concept.Key().Code())),
			false,
			nil),
	}
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
		"label":    s.Key().Code(),
		"arangodb": acfg,
	}, nil
}

func (s *ArangodbStore) Urls() []string {
	return s.urls
}

func (s *ArangodbStore) Credentials() (string, string) {
	return s.username, s.password
}

func (s *ArangodbStore) Database() string {
	return s.database
}

func (s *ArangodbStore) Collection() string {
	return s.collection
}

type arangodbClient struct {
	client     driver.Client
	db         driver.Database
	collection driver.Collection
}

func (c *arangodbClient) exists(ctx context.Context, request shono.StoreOperationRequest) shono.StoreOperationResponse {
	adbKey := keyToArangodbKey(request.Key())
	b, err := c.collection.DocumentExists(ctx, adbKey)
	if err != nil {
		return newResponse(http.StatusInternalServerError, err)
	}

	return newResponse(http.StatusOK, b)
}

func keyToArangodbKey(key shono.Key) string {
	return strings.ReplaceAll(key.String(), ":", "_")
}
