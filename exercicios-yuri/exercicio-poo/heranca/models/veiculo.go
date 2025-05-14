package models

import "fmt"

type Veiculo struct {
	Marca      string
	Velocidade int
}

func (v Veiculo) Ligar() {
	fmt.Println("Ligou o veiculo")
}
