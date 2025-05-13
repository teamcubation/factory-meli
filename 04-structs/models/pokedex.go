package models

import "fmt"

type Pokedex struct {
	Pokemons map[int]*Pokemon
}

func (pkdx *Pokedex) AddPokemon(id int, pkm *Pokemon) {
	if pkdx.Pokemons == nil {
		pkdx.Pokemons = make(map[int]*Pokemon)
	}

	pkdx.Pokemons[id] = pkm
}

func (pkdx *Pokedex) GetPokemon(id int) (*Pokemon, bool) {
	pkm, exists := pkdx.Pokemons[id]

	return pkm, exists
}

func (pkdx *Pokedex) ListAll() []*Pokemon {
	pokemons := []*Pokemon{}

	for _, pkm := range pkdx.Pokemons {
		pokemons = append(pokemons, pkm)
	}

	return pokemons
}

func (pkdx Pokedex) String() string {
	pokemonStrings := ""
	for _, v := range pkdx.Pokemons {
		pokemonStrings += v.String() + "\n"
	}

	return fmt.Sprintf("Pokedex: \n Pokemons: \n %s", pokemonStrings)
}
