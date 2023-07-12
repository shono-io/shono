package dsl

#Step: {
	addToStore?: #AddToStore
	asSuccessEvent?: #AsSuccessEvent
	asFailureEvent?: #AsFailureEvent
	catch?: #Catch
	getFromStore?: #GetFromStore
	listFromStore?: #ListFromStore
	log?: #Log
	raw?: #Raw
	removeFromStore?: #RemoveFromStore
	setInStore?: #SetInStore
	switch?: #Switch
}

#AddToStore: {
	scope: string
	concept: string
	key: string
}

#AsSuccessEvent: {
	scope: string
	concept: string
	event: string
	status: int
}

#AsFailureEvent: {
	scope: string
	concept: string
	event: string
	errorCode: int
	reason: string
}

#Catch: [...#Step]

#Filter: {
	field: string
	value: string | int | float | bool
}

#GetFromStore: {
	scope: string
	concept: string
	key: string
}

#ListFromStore: {
	scope: string
	concept: string
	filters: [...#Filter]
}

#Log: {
	level: "debug" | *"info" | "warn" | "error"
	message: string
}

#Raw: _

#RemoveFromStore: {
	scope: string
	concept: string
	key: string
}

#SetInStore: {
	scope: string
	concept: string
	key: string
}

#Switch: [...#SwitchCase]

#SwitchCase: {
	condition?: string
	steps: [...#Step]
}