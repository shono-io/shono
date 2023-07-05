import (
	"shono.io/spec/v1:commons"
	"shono.io/spec/v1:dsl"
)

#Spec: {
	scope?: #Scope
	concept?: #Concept
	event?: #Event
}

#Scope: {
	code: string
	summary: string
	status: commons.#Status
	docs?: string
}

#Concept: {
	scope: string
	code: string
	summary: string
	status: commons.#Status
	docs?: string
	stored: bool | *false
}

#Event: {
	scope: string
	concept: string
	code: string
	summary: string
	status: commons.#Status
	docs?: string
}

#Reactor: {
	scope: string
	concept: string
	code: string
	summary: string
	status: commons.#Status
	docs?: string
	logic: [...dsl.#Step]
}