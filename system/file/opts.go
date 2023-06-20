package file

type Opt func(config map[string]any)

func WithOutputPath(path string) Opt {
	return func(config map[string]any) {
		config["path"] = path
	}
}

func WithInputPath(path string) Opt {
	return func(config map[string]any) {
		cp, fnd := config["paths"]
		if !fnd {
			cp = []string{}
		}
		cp = append(cp.([]string), path)
		config["paths"] = cp
	}
}

func WithCodec(codec string) Opt {
	return func(config map[string]any) {
		config["codec"] = codec
	}
}
