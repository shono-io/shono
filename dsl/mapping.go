package dsl

import "fmt"

type Mapping struct {
	Language   string
	Sourcecode string
}

func (m Mapping) Validate() error {
	if m.Language == "" {
		return fmt.Errorf("transform logic must have a language")
	}

	if m.Sourcecode == "" {
		return fmt.Errorf("transform logic must have sourcecode")
	}

	return nil
}
