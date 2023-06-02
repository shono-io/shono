package benthos

import "github.com/shono-io/shono"

type ReaktorTestOpt func(*reaktorTest)

func WithMock(name string, mock any) ReaktorTestOpt {
	return func(t *reaktorTest) {
		t.mocks[name] = mock
	}
}

func WithEnvironment(key string, value any) ReaktorTestOpt {
	return func(t *reaktorTest) {
		t.environment[key] = value
	}
}

func WithCondition(condition reaktorTestCondition) ReaktorTestOpt {
	return func(t *reaktorTest) {
		t.then = append(t.then, condition)
	}
}

func NewReaktorTest(summary string, when *reaktorTestEvent, opts ...ReaktorTestOpt) shono.ReaktorTest {
	result := &reaktorTest{
		summary:     summary,
		when:        when,
		then:        []reaktorTestCondition{},
		environment: map[string]any{},
		mocks:       map[string]any{},
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type reaktorTest struct {
	summary     string
	environment map[string]any
	mocks       map[string]any
	when        *reaktorTestEvent
	then        []reaktorTestCondition
}

func (r *reaktorTest) Summary() string {
	return r.summary
}

func (r *reaktorTest) Environment() map[string]any {
	return r.environment
}

func (r *reaktorTest) Mocks() map[string]any {
	return r.mocks
}

func (r *reaktorTest) When() shono.ReaktorTestEvent {
	return r.when
}

func (r *reaktorTest) Then() []shono.ReaktorTestCondition {
	// FIXME: I hate doing this, there must be a better way
	result := make([]shono.ReaktorTestCondition, len(r.then))
	for i, condition := range r.then {
		result[i] = condition
	}
	return result
}

func (r *reaktorTest) AsBenthos() map[string]any {
	inputs := []map[string]any{
		r.when.AsBenthos(),
	}

	outputs := []map[string]any{}
	for _, condition := range r.then {
		outputs = append(outputs, condition.AsBenthos())
	}

	return map[string]any{
		"name":        r.summary,
		"environment": r.environment,
		"mocks":       r.mocks,
		"input_batch": inputs,
		"output_batches": [][]map[string]any{
			outputs,
		},
	}
}

func WithStringContent(content string) ReaktorTestEventOpt {
	return func(event *reaktorTestEvent) {
		event.contentType = "content"
		event.content = content
	}
}

func WithJsonContent(content map[string]any) ReaktorTestEventOpt {
	return func(event *reaktorTestEvent) {
		event.contentType = "json_content"
		event.content = content
	}
}

func WithMetadata(key string, value any) ReaktorTestEventOpt {
	return func(event *reaktorTestEvent) {
		event.metadata[key] = value
	}
}

type ReaktorTestEventOpt func(*reaktorTestEvent)

func NewReaktorTestEvent(opts ...ReaktorTestEventOpt) shono.ReaktorTestEvent {
	result := &reaktorTestEvent{
		metadata: map[string]any{},
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type reaktorTestEvent struct {
	metadata    map[string]any
	content     any
	contentType string
}

func (r *reaktorTestEvent) Metadata() map[string]any {
	return r.metadata
}

func (r *reaktorTestEvent) Content() any {
	return r.content
}

func (r *reaktorTestEvent) AsBenthos() map[string]any {
	return map[string]any{
		"metadata":    r.metadata,
		r.contentType: r.content,
	}
}

type reaktorTestCondition struct {
	kind      string
	condition any
}

func (r reaktorTestCondition) Kind() string {
	return r.kind
}

func (r reaktorTestCondition) Condition() any {
	return r.condition
}

func (r reaktorTestCondition) AsBenthos() map[string]any {
	return map[string]any{
		r.kind: r.condition,
	}
}
