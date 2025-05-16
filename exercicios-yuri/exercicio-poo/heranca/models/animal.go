package models

import "fmt"

type Animal struct {
	Nome string
}

func (a Animal) Comer() {
	fmt.Printf("O animal %s está comendo\n", a.Nome)
}
