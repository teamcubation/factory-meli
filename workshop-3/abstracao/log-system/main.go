package main

import "fmt"

type Logger interface {
	Log(nivel, menssagem string)
}

type ConsoleLogger struct{}

func (c ConsoleLogger) Log(nivel, menssagem string) {
	fmt.Printf("[%s] %s\n", nivel, menssagem)
}

type FileLogger struct {
	file string
}

func (f FileLogger) Log(nivel, menssagem string) {
	fmt.Println("[" + f.file + "][" + nivel + "]" + menssagem)
}

func ProcessarTarefa(logger Logger) {
	logger.Log("INFO", "Logger is working")
}

func main() {
	console := ConsoleLogger{}
	file := FileLogger{"test.txt"}

	ProcessarTarefa(console)
	ProcessarTarefa(file)

}
