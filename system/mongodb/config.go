package mongodb

import "github.com/shono-io/shono/inventory"

var (
	UriField      = "uri"
	DatabaseField = "database"
)

var configFields = []inventory.IOConfigSpecField{
	{UriField, "string", "scalar", false, nil},
	{DatabaseField, "string", "scalar", false, nil},
}
