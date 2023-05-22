package go_shono

type Writer interface {
	MustWrite(correlationId string, evt *EventMeta, payload any)
	Write(correlationId string, evt *EventMeta, payload any) error
}
