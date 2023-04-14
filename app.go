package go_shono

import "github.com/shono-io/go-shono/reaktor"

type App struct {
	Config   AppConfig
	Handlers reaktor.Handlers
}

func (a App) Run() error {
	r, err := reaktor.NewReaktor(a.Handlers, a.Config)
	if err != nil {
		return err
	}

	return r.Start()
}
