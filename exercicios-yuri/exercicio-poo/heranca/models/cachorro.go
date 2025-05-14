package models

import "fmt"

type Cachorro struct {
	Embed Animal
}

func (c Cachorro) Latir() {
	fmt.Printf("O cachorro está latindo\n")
}
