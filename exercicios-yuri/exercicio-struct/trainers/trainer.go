package trainers

import (
	p "Yuri/exercicios-yuri/pokemon/pokemons"
	"errors"
	"fmt"
)

type Trainer struct {
	Name  string
	Party []p.Pokemon
}

func (t *Trainer) AddToParty(pokemon p.Pokemon) error {
	if len(t.Party) >= 6 {
		return errors.New("AddToParty: seu time está cheio")
	}
	t.Party = append(t.Party, pokemon)
	return nil
}

func (t *Trainer) Infos() {
	fmt.Printf("Nome do Treinador: %s\n", t.Name)
	fmt.Println("Equipe:")
	for i, pokemon := range t.Party {
		fmt.Printf("  %d. Nome: %s | Tipo: %v | Nível: %d\n", i+1, pokemon.Name, pokemon.Type, pokemon.Level)
	}
}
