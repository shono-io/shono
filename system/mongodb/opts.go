package mongodb

type Opt func(map[string]any)

func WithUri(uri string) Opt {
	return func(config map[string]any) {
		config[UriField] = uri
	}
}

func WithDatabase(database string) Opt {
	return func(config map[string]any) {
		config[DatabaseField] = database
	}
}
