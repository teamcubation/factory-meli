package models

import "fmt"

type Gato struct {
	Embed Animal
}

func (g Gato) Miar() {
	fmt.Printf("O gato está miando\n")
}
