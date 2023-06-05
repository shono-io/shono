package graph

type ScopeOpt func(s *Scope)

func WithScopeName(name string) ScopeOpt {
	return func(s *Scope) {
		s.name = name
	}
}

func WithScopeDescription(description string) ScopeOpt {
	return func(s *Scope) {
		s.description = description
	}
}

func NewScope(key Key, opts ...ScopeOpt) Scope {
	result := Scope{
		key,
		key.Code(),
		"",
	}

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

type Scope struct {
	key         Key
	name        string
	description string
}

func (s Scope) Key() Key {
	return s.key
}

func (s Scope) Name() string {
	return s.name
}

func (s Scope) Description() string {
	return s.description
}
