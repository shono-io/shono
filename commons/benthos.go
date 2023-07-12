package commons

type BenthosMarshaller interface {
	MarshalBenthos() (map[string]any, error)
}
