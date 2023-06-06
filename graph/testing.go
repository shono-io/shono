package graph

func NewReaktorTest(summary string) ReaktorTest {
	return ReaktorTest{
		Summary:     summary,
		Event:       nil,
		Conditions:  []ReaktorTestCondition{},
		Environment: map[string]any{},
		Mocks:       map[string]any{},
	}
}

type ReaktorTest struct {
	Summary     string
	Environment map[string]any
	Mocks       map[string]any
	Event       *ReaktorTestEvent
	Conditions  []ReaktorTestCondition
}

func (t *ReaktorTest) Given(environment map[string]any) *ReaktorTest {
	t.Environment = environment
	return t
}

func (t *ReaktorTest) When(event ReaktorTestEvent) *ReaktorTest {
	t.Event = &event
	return t
}

func (t *ReaktorTest) Then(conditions ...ReaktorTestCondition) *ReaktorTest {
	for _, condition := range conditions {
		t.Conditions = append(t.Conditions, condition)
	}
	return t
}

func NewReaktorTestEvent(content map[string]any, metadata map[string]string) ReaktorTestEvent {
	return ReaktorTestEvent{
		Metadata: metadata,
		Content:  content,
	}
}

type ReaktorTestEvent struct {
	Metadata map[string]string
	Content  map[string]any
}

type ReaktorTestCondition interface {
	ConditionType() string
}

func Evaluates(expression string) ReaktorTestCondition {
	return BloblangReaktorTestCondition{
		Expression: expression,
	}
}

type BloblangReaktorTestCondition struct {
	Expression string
}

func (c BloblangReaktorTestCondition) ConditionType() string {
	return "bloblang"
}

func HasMetadata(expected map[string]string) ReaktorTestCondition {
	return MetadataReaktorTestCondition{
		Values: expected,
		Strict: false,
	}
}

func HasStrictMetadata(expected map[string]string) ReaktorTestCondition {
	return MetadataReaktorTestCondition{
		Values: expected,
		Strict: true,
	}
}

type MetadataReaktorTestCondition struct {
	Values map[string]string
	Strict bool
}

func (c MetadataReaktorTestCondition) ConditionType() string {
	return "metadata"
}

func HasPayload(values map[string]interface{}) PayloadReaktorTestCondition {
	return PayloadReaktorTestCondition{
		Values: values,
		Strict: false,
	}
}

func HasStrictPayload(values map[string]interface{}) PayloadReaktorTestCondition {
	return PayloadReaktorTestCondition{
		Values: values,
		Strict: true,
	}
}

type PayloadReaktorTestCondition struct {
	Values map[string]interface{}
	Strict bool
}

func (c PayloadReaktorTestCondition) ConditionType() string {
	return "payload"
}
