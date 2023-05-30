package shono

type ReaktorOpt func(reaktor *reaktor)

func WithName(name string) ReaktorOpt {
	return func(reaktor *reaktor) {
		reaktor.entity.name = name
	}
}

func WithReaktorDescription(description string) ReaktorOpt {
	return func(reaktor *reaktor) {
		reaktor.entity.description = description
	}
}

func WithOutputEvent(outputEvents ...EventId) ReaktorOpt {
	return func(reaktor *reaktor) {
		reaktor.outputEvents = append(reaktor.outputEvents, outputEvents...)
	}
}

func WithStore(stores ...Store) ReaktorOpt {
	return func(reaktor *reaktor) {
		reaktor.stores = append(reaktor.stores, stores...)
	}
}
