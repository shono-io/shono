package inventory

func AssertMetadataEquals(expected map[string]string) TestAssertion {
	return NewMetadataTestAssertion(expected, true)
}

func AssertMetadataContains(key, value string) TestAssertion {
	return NewMetadataTestAssertion(map[string]string{
		key: value,
	}, false)
}

func AssertContentEquals(expected map[string]interface{}) TestAssertion {
	return NewPayloadTestAssertion(expected, true)
}

func AssertContentContains(key string, value interface{}) TestAssertion {
	return NewPayloadTestAssertion(map[string]interface{}{
		key: value,
	}, false)
}

func NewMetadataTestAssertion(values map[string]string, strict bool) TestAssertion {
	return testAssertion{
		metadata: values,
		strict:   strict,
	}
}

func NewPayloadTestAssertion(values map[string]interface{}, strict bool) TestAssertion {
	return testAssertion{
		payload: values,
		strict:  strict,
	}
}

type testAssertion struct {
	metadata map[string]string
	payload  map[string]interface{}
	strict   bool
}

func (a testAssertion) Metadata() map[string]string {
	return a.metadata
}

func (a testAssertion) Payload() map[string]interface{} {
	return a.payload
}

func (a testAssertion) Strict() bool {
	return a.strict
}
