package dsl

import "fmt"

func Log(level LogLevel, message string) LogLogicStep {
	return LogLogicStep{
		Level:   level,
		Message: message,
	}
}

func LogWithMapping(level LogLevel, language string, sourcecode string) LogLogicStep {
	return LogLogicStep{
		Level: level,
		Mapping: &Mapping{
			Language:   language,
			Sourcecode: sourcecode,
		},
	}
}

type LogLevel string

var (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type LogLogicStep struct {
	Level   LogLevel `yaml:"level"`
	Message string   `yaml:"message,omitempty"`
	Mapping *Mapping `yaml:"mapping,omitempty"`
}

func (e LogLogicStep) Kind() string {
	return "log"
}

func (e LogLogicStep) Validate() error {
	if e.Level == "" {
		return fmt.Errorf("no log level set")
	}

	if e.Level != LogLevelDebug && e.Level != LogLevelInfo && e.Level != LogLevelWarn && e.Level != LogLevelError {
		return fmt.Errorf("invalid log level; one of debug, info, warn, error expected")
	}

	if e.Message == "" && e.Mapping == nil {
		return fmt.Errorf("neither a message or a mapping has been defined")
	}

	if e.Message != "" && e.Mapping != nil {
		return fmt.Errorf("both a message and mapping have been defined")
	}

	if e.Mapping != nil {
		if err := e.Mapping.Validate(); err != nil {
			return err
		}
	}

	return nil
}
