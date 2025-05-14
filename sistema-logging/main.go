package main

import (
	"example.com/sistema-logging/logger"
)

func main() {
	consoleLogger := logger.ConsoleLogger{}
	fileLogger := logger.FileLogger{}

	logger.ProcessarTarefa(consoleLogger)
	logger.ProcessarTarefa(fileLogger)
}
