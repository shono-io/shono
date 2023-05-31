package shono

import "context"

type Store interface {
	Entity

	AsBenthosComponent() (map[string]any, error)

	Operations() map[StoreOpertationId]StoreOperation
}

type StoreOpertationId string

var (
	ExistsOperation = StoreOpertationId("exists")
	GetOperation    = StoreOpertationId("get")
	ListOperation   = StoreOpertationId("list")

	CreateOperation  = StoreOpertationId("create")
	ReplaceOperation = StoreOpertationId("replace")
	UpdateOperation  = StoreOpertationId("update")
	DeleteOperation  = StoreOpertationId("delete")
)

type StoreOperationHandler func(ctx context.Context, request StoreOperationRequest) StoreOperationResponse

type StoreOperation interface {
	Id() StoreOpertationId
	Name() string
	Title() string
	Description() string
	IsModifier() bool

	Handler() StoreOperationHandler
}

type StoreOperationRequest interface {
	Key() Key
	HasParameter(name string) bool
	StringParameter(name string) string
	IntParameter(name string) int
	FloatParameter(name string) float64
	JsonParameter(name string, target any) error
}

type StoreOperationResponse interface {
	Code() int
	Payload() any
}

type StoreClient interface {
	Apply(ctx context.Context, request StoreOperationRequest) StoreOperationResponse
}
