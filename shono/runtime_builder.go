package shono

type RuntimeBuilder struct {
}

func (b *RuntimeBuilder) WithReaktor(reaktor ...Reaktor) *RuntimeBuilder {
	for _, r := range reaktor {
		b.withReaktor(r)
	}
	return b
}

func (b *RuntimeBuilder) withReaktor(reaktor Reaktor) *RuntimeBuilder {

	return b
}

func (b *RuntimeBuilder) Build() (*Runtime, error) {
	return nil, nil
}
