package mongodb

import (
	"github.com/shono-io/shono/systems"
)

func init() {
	systems.RegisterStorageSystem("mongodb", &system{})
}
