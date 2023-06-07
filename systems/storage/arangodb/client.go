package arangodb

import (
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/benthosdev/benthos/v4/public/service"
)

func clientFields() []*service.ConfigField {
	return []*service.ConfigField{
		service.NewStringListField("urls").
			Default("http://localhost:8529").
			Description("The URLs of the ArangoDB server"),
		service.NewStringField("username").
			Default("root").
			Description("The username to use when connecting to the ArangoDB server"),
		service.NewStringField("password").
			Default("").
			Description("The password to use when connecting to the ArangoDB server").
			Secret(),
	}
}

func getClientFromConfig(conf *service.ParsedConfig) (driver.Client, error) {
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

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: url,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create arangodb connection: %w", err)
	}

	return driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(username, password),
	})
}
