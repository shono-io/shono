package graph

func Log(level LogLevel, message Expression) LogLogic {
	return LogLogic{
		Level:   level,
		Message: message,
	}
}

type LogLevel string

type LogLogic struct {
	Level    LogLevel
	Message  Expression
	Mappings []Mapping
}

func (e LogLogic) Kind() string {
	return "log"
}
