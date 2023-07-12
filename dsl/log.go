package dsl

import (
	"fmt"
	"github.com/shono-io/shono/inventory"
	"strings"
)

type LogBuilder struct {
	result *LogLogicStep
}

func (lb *LogBuilder) Label(label string) *LogBuilder {
	lb.result.label = label
	return lb
}

func (lb *LogBuilder) Message(message string) *LogBuilder {
	lb.result.Message = message
	return lb
}

func (lb *LogBuilder) Mapping(mapping string) *LogBuilder {
	lb.result.Mapping = mapping
	return lb
}

func (lb *LogBuilder) Build() inventory.LogicStep {
	return *lb.result
}

func Log(level LogLevel) *LogBuilder {
	return &LogBuilder{result: &LogLogicStep{Level: level}}
}

type LogLevel string

var (
	LogLevelTrace LogLevel = "trace"
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type LogLogicStep struct {
	label   string
	Level   LogLevel
	Message string
	Mapping string
}

func (e LogLogicStep) Label() string {
	return e.label
}

func (e LogLogicStep) MarshalBenthos(trace string) (map[string]any, error) {
	trace = fmt.Sprintf("%s/%s", trace, e.Kind())

	if err := e.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", trace, err)
	}

	l := map[string]any{
		"level": strings.ToLower(string(e.Level)),
	}

	if e.Message != "" {
		l["message"] = e.Message
	}

	if e.Mapping != "" {
		l["fields_mapping"] = CleanMultilineStringWhitespace(e.Mapping)
	}

	result := map[string]any{
		"log": l,
	}

	if e.label != "" {
		result["label"] = e.label
	}

	return result, nil
}

func (e LogLogicStep) Kind() string {
	return "log"
}

func (e LogLogicStep) Validate() error {
	if e.Level == "" {
		return fmt.Errorf("no log level set")
	}

	lvl := LogLevel(strings.ToLower(string(e.Level)))
	if lvl != LogLevelTrace &&
		lvl != LogLevelDebug &&
		lvl != LogLevelInfo &&
		lvl != LogLevelWarn &&
		lvl != LogLevelError {
		return fmt.Errorf("invalid log level %q; one of trace, debug, info, warn, error expected", e.Level)
	}

	if e.Message == "" && e.Mapping == "" {
		return fmt.Errorf("neither a message or a mapping has been defined")
	}

	if e.Message != "" && e.Mapping != "" {
		return fmt.Errorf("both a message and mapping have been defined")
	}

	return nil
}
