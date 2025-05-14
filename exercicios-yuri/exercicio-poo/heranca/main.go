package main

import (
	"Yuri/exercises-yuri/heranca/models"
	"fmt"
)

func main() {
	cao := models.Cachorro{Embed: models.Animal{Nome: "Scooby"}}
	gato := models.Gato{Embed: models.Animal{Nome: "Mika"}}

	fmt.Println("Cachorro", cao)
	fmt.Println("Gato", gato)

	cao.Latir()
	gato.Miar()

	cao.Embed.Comer()
	gato.Embed.Comer()

	// ---------------------------

	carro := models.Carro{
		Veiculo: models.Veiculo{
			Marca: "Nissan", Velocidade: 120,
		},
		Portas: 4,
	}

	bicicleta := models.Bicicleta{
		Veiculo: models.Veiculo{
			Marca: "Caloi", Velocidade: 40,
		},
		TemCampainha: true,
	}

	carro.Veiculo.Ligar()
	carro.AbrirPortas()
	bicicleta.Veiculo.Ligar()
	bicicleta.TocarCampanhinha()

}
