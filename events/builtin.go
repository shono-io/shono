package events

type OperationFailed struct {
	Reason string `json:"reason"`
}
