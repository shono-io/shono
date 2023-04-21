package app

type Module interface {
	Name() string
	Init(catalogs *Catalogs) error
}
