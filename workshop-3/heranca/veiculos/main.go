package main

import "fmt"

type Veiculo struct {
	Marca      string
	Velocidade int
}

func (v Veiculo) Ligar() {
	fmt.Println("O veículo está ligado e pronto para partir")
}

type Carro struct {
	Veiculo
	Portas int
}

func CarroFactory(marca string, velocidade, portas int) *Carro {
	return &Carro{Veiculo{marca, velocidade}, portas}
}

type Bicicleta struct {
	Veiculo
	TemCampainha bool
}

func BicicletaFactory(marca string, velocidade int, temCampainha bool) *Bicicleta {
	return &Bicicleta{Veiculo{marca, velocidade}, temCampainha}
}

func main() {

	carro := CarroFactory("Mercedes", 300, 2)
	carro.Ligar()

	bicicleta := BicicletaFactory("Caloi", 60, false)
	bicicleta.Ligar()
}
