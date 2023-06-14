package core

type System interface {
	Node
	AsInput(config map[string]any) (map[string]any, error)
	AsOutput(config map[string]any) (map[string]any, error)
}
