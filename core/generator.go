package core

type ReaktorGenerator interface {
	// Generate create a new reaktor artifact for all reactors for a given concept.
	Generate(env Environment, conceptRef Reference) (Artifact, error)
}

type InjectorGenerator interface {
	// Generate create a new injector artifact for the given injector.
	Generate(env Environment, injectorRef Reference) (Artifact, error)
}

type ExtractorGenerator interface {
	// Generate create a new extractor artifact for the given extractor.
	Generate(env Environment, extractorRef Reference) (Artifact, error)
}
