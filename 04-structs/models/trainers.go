package models

import (
	"errors"
	"fmt"
)

type Trainer struct {
	Name  string
	Party []*Pokemon
}

const MaxPartySize int = 6

func (t *Trainer) AddToParty(p *Pokemon) error {
	if len(t.Party) >= MaxPartySize {
		return errors.New("o time do treinador já está cheio")
	}

	if t.Party == nil {
		t.Party = make([]*Pokemon, 0, MaxPartySize)
	}

	t.Party = append(t.Party, p)

	return nil
}

func (t Trainer) String() string {
	result := fmt.Sprintf("Treinador: %s\n Equipe: \n", t.Name)

	for i, p := range t.Party {
		result += fmt.Sprintf("%d - %s\n", i+1, p.String())
	}

	return result
}
