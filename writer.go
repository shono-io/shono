package go_shono

type Writer interface {
	MustWrite(evt *EventMeta, payload any)
	Write(evt *EventMeta, payload any) error
}
