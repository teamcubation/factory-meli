package models

import "fmt"

type FileLogger struct{}

func (f FileLogger) Log(nivel, mensagem string) {
	fmt.Printf("[FILE][%s] %s\n", nivel, mensagem)
}
