package logic

type EngineId string

type Logic interface {
	EngineId() EngineId
	Processor() (map[string]any, error)
}
