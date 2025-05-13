package pokemons

import (
	"errors"
	"fmt"
)

type Pokemon struct {
	Name      string
	Type      []string
	Level     int
	EvolvesTo string
}

func (p *Pokemon) Infos() {
	fmt.Printf("Nome: %s\nTipo: %v\nNível: %d\nEvolução: %s\n", p.Name, p.Type, p.Level, p.EvolvesTo)
}

func (p *Pokemon) Evolve() error {
	if p.EvolvesTo == "" {
		return errors.New("Evolve: Pokemon não possui evolução")
	}
	if p.Level >= 16 {
		fmt.Printf("%s evoluiu para %s\n", p.Name, p.EvolvesTo)
		p.Name = p.EvolvesTo
		p.EvolvesTo = ""
		return nil
	}
	return errors.New("Evolve: Pokemon não pode evoluir ainda")
}

// TODO: FAZER A TERCEIRA EVOLUÇÃO
// func (p *Pokemon) AddEvolve(evolution string){
// 	p.EvolvesTo = evolution
// }
