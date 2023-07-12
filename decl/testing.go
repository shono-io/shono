package decl

type TestSpec struct {
	Summary string         `yaml:"summary"`
	Given   *TestGivenSpec `yaml:"given,omitempty"`
	When    TestWhenSpec   `yaml:"when"`
	Then    TestThenSpec   `yaml:"then"`
}

type TestGivenSpec struct {
	Environment *map[string]any `yaml:"environment,omitempty"`
}

type TestWhenSpec struct {
	Raw   *TestWhenRawSpec `yaml:"raw,omitempty"`
	Event *EventRef        `yaml:"event,omitempty"`
}

type TestWhenRawSpec struct {
	Content    *string            `yaml:"content,omitempty"`
	Structured *map[string]any    `yaml:"structured,omitempty"`
	Metadata   *map[string]string `yaml:"metadata,omitempty"`
}

type TestThenSpec struct {
	MetadataEquals   *map[string]string `yaml:"metadataEquals,omitempty"`
	MetadataContains *map[string]string `yaml:"metadataContains,omitempty"`
	ContentEquals    *map[string]any    `yaml:"contentEquals,omitempty"`
	ContentContains  *map[string]any    `yaml:"contentContains,omitempty"`
}
