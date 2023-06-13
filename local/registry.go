package local

import "github.com/shono-io/shono/graph"

type Registry struct {
	backboneRepo
	scopeRepo
	conceptRepo
	eventRepo
	reaktorRepo
	storageRepo
}

type Opt func(*Registry)

func WithBackbone(bb graph.Backbone) Opt {
	return func(r *Registry) {
		r.bb = bb
	}
}

func WithScope(scope ...graph.Scope) Opt {
	return func(r *Registry) {
		for _, v := range scope {
			r.scopes[v.Code] = v
		}
	}
}

func WithConcept(concept ...graph.Concept) Opt {
	return func(r *Registry) {
		for _, v := range concept {
			r.concepts[v.ConceptReference.String()] = v
		}
	}
}

func WithEvent(event ...graph.Event) Opt {
	return func(r *Registry) {
		for _, v := range event {
			r.events[v.EventReference.String()] = v
		}
	}
}

func WithReaktor(reaktor ...graph.Reaktor) Opt {
	return func(r *Registry) {
		for _, v := range reaktor {
			r.reaktors[v.ReaktorReference.String()] = v
		}
	}
}

func WithStorage(storage ...graph.Storage) Opt {
	return func(r *Registry) {
		for _, v := range storage {
			r.storage[v.Key()] = v
		}
	}
}

func NewRegistry(opts ...Opt) *Registry {
	r := &Registry{
		backboneRepo: backboneRepo{},
		scopeRepo: scopeRepo{
			scopes: map[string]graph.Scope{},
		},
		conceptRepo: conceptRepo{
			concepts: map[string]graph.Concept{},
		},
		eventRepo: eventRepo{
			events: map[string]graph.Event{},
		},
		reaktorRepo: reaktorRepo{
			reaktors: map[string]graph.Reaktor{},
		},
		storageRepo: storageRepo{
			storage: map[string]graph.Storage{},
		},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
