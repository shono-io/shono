package logic

type Engine interface {
	EngineId() EngineId
	Parse(source string) (map[string]any, error)
}
