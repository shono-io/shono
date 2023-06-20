package arangodb

import "github.com/shono-io/shono/inventory"

var (
	UrlsField     = "urls"
	UsernameField = "username"
	PasswordField = "password"
	DatabaseField = "database"
)

var configFields = []inventory.IOConfigSpecField{
	{UrlsField, "string", "list", false, nil},
	{UsernameField, "string", "scalar", false, nil},
	{PasswordField, "string", "scalar", false, nil},
	{DatabaseField, "string", "scalar", false, nil},
}
