package pokedex

import (
	"Yuri/exercicios-yuri/pokemon/pokemons"
	"errors"
)

type Pokedex struct {
	pokemons map[int]pokemons.Pokemon
}

func (p *Pokedex) AddPokemon(id int, poke pokemons.Pokemon) {
	if p.pokemons == nil {
		p.pokemons = make(map[int]pokemons.Pokemon)
	}
	p.pokemons[id] = poke
}

func (p *Pokedex) GetPokemon(id int) (pokemons.Pokemon, bool, error) {
	var response pokemons.Pokemon
	for i, pokemon := range p.pokemons {
		if i == id {
			response = pokemon
			return response, true, nil
		}
	}
	return response, false, errors.New("GetPokemon: Nenhum pokemon encontrado")
}

func (p *Pokedex) ListAll() []pokemons.Pokemon {
	var response []pokemons.Pokemon
	for _, pokemon := range p.pokemons {
		response = append(response, pokemon)
	}
	return response
}
