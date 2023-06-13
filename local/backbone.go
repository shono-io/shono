package local

import (
	"github.com/shono-io/shono/graph"
)

type backboneRepo struct {
	bb graph.Backbone
}

func (b *backboneRepo) GetBackbone() (graph.Backbone, error) {
	return b.bb, nil
}
