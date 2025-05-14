package models

import "fmt"

type ConsoleLogger struct {
}

func (c ConsoleLogger) Log(nivel, mensagem string) {
	fmt.Printf("[%s] %s\n", nivel, mensagem)
}
