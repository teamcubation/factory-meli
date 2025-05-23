package models

import (
	"errors"
	"fmt"
)

type Pokemon struct {
	Name      string
	Types     []string
	Level     int
	EvolvesTo string
}

func (p *Pokemon) Evolve() error {
	if p.Level < 16 {
		return errors.New("pokemon ainda não pode evoluir")
	}

	fmt.Printf("¡%s evolui para %s!", p.Name, p.EvolvesTo)
	p.Name = p.EvolvesTo

	return nil
}

func (p Pokemon) String() string {
	return fmt.Sprintf("Pokemon: %s\n Informações: \n Nível: %d \n Tipos: %v \n Próxima evolução: %s \n", p.Name, p.Level, p.Types, p.EvolvesTo)
}
