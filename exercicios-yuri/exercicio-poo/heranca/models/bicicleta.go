package models

import "fmt"

type Bicicleta struct {
	Veiculo      Veiculo
	TemCampainha bool
}

func (b Bicicleta) TocarCampanhinha() {
	if b.TemCampainha {
		fmt.Println("Tocou a campanhinha")
	} else {
		fmt.Println("Não tem campanhinha")
	}
}
