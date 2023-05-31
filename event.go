package shono

type Event interface {
	Entity
	Id() EventId
}

type EventOpt func(*event)

func WithEventName(name string) EventOpt {
	return func(e *event) {
		e.entity.name = name
	}
}

func WithEventDescription(description string) EventOpt {
	return func(e *event) {
		e.entity.description = description
	}
}

func NewEvent(conceptKey Key, code string, opts ...EventOpt) Event {
	result := &event{
		newEntity(conceptKey.Child("event", code)),
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type event struct {
	*entity
}

func (e *event) Id() EventId {
	if e == nil {
		return ""
	}

	return EventId(e.Key().String())
}
