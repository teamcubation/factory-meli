package logger

import "fmt"

type LogLevel int

const (
	Trace LogLevel = iota
	Debug
	Info
	Error
)

func (option LogLevel) String() string {
	switch option {
	case Debug:
		return "Debug"
	case Error:
		return "Error"
	case Info:
		return "Info"
	case Trace:
		return "Trace"
	default:
		return ""
	}
}

type Logger interface {
	Log(nivel LogLevel, mensagem string)
}

type ConsoleLogger struct{}

func (c ConsoleLogger) Log(nivel LogLevel, mensagem string) {
	fmt.Printf("%s - %s\n", nivel, mensagem)
}

type FileLogger struct{}

func (f FileLogger) Log(nivel LogLevel, mensagem string) {
	fmt.Printf("[FILE] - %s - %s\n", nivel, mensagem)
}

func ProcessarTarefa(logger Logger) {
	logger.Log(Info, "Hello World")
}
