package arangodb

type Opt func(map[string]any)

func WithUrls(urls ...string) Opt {
	return func(config map[string]any) {
		config[UrlsField] = urls
	}
}

func WithUsername(username string) Opt {
	return func(config map[string]any) {
		config[UsernameField] = username
	}
}

func WithPassword(password string) Opt {
	return func(config map[string]any) {
		config[PasswordField] = password
	}
}

func WithDatabase(database string) Opt {
	return func(config map[string]any) {
		config[DatabaseField] = database
	}
}
