package event

type Catalog interface {
	GetEventType(id EventId) (*EventType, error)

	RegisterEventType(event *EventType) error
}

func NewMemoryCatalog() *MemoryCatalog {
	return &MemoryCatalog{
		EventTypes: map[EventId]*EventType{},
	}
}

type MemoryCatalog struct {
	EventTypes map[EventId]*EventType
}

func (m *MemoryCatalog) GetEventType(id EventId) (*EventType, error) {
	if et, ok := m.EventTypes[id]; ok {
		return et, nil
	}

	return nil, nil
}

func (m *MemoryCatalog) RegisterEventType(event *EventType) error {
	m.EventTypes[event.EventId] = event

	return nil
}
