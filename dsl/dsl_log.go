package dsl

import "github.com/shono-io/shono/graph"

func Log(level LogLevel, message graph.Expression) LogLogic {
	return LogLogic{
		Level:   level,
		Message: message,
	}
}

type LogLevel string

type LogLogic struct {
	Level    LogLevel         `yaml:"level"`
	Message  graph.Expression `yaml:"message"`
	Mappings []Mapping        `yaml:"mappings"`
}

func (e LogLogic) Kind() string {
	return "log"
}
