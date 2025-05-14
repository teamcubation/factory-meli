package models

import "fmt"

type Carro struct {
	Veiculo Veiculo
	Portas  int
}

func (c Carro) AbrirPortas() {
	fmt.Printf("Abriu %d portas\n", c.Portas)
}
