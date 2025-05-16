package main

import (
	m "Yuri/exercicios-yuri/pokemon/moves"
	pd "Yuri/exercicios-yuri/pokemon/pokedex"
	p "Yuri/exercicios-yuri/pokemon/pokemons"
	t "Yuri/exercicios-yuri/pokemon/trainers"
	"fmt"
)

func ShowInformations(typeInformation TypeInformation) {
	typeInformation.Infos()
}

type TypeInformation interface {
	Infos()
}

func main() {

	pokemons := createPokemons()
	moves := createMoves()
	trainer := t.Trainer{
		Name:  "Blue",
		Party: []p.Pokemon{*pokemons[0]},
	}
	generateInfos(pokemons)
	trainer.Infos()
	generateInfos(moves)

	pokemons[0].Evolve()
	pokeErr1 := pokemons[1].Evolve()
	if pokeErr1 != nil {
		fmt.Println("Error:", pokeErr1)
	}
	pokeErr2 := pokemons[2].Evolve()
	if pokeErr2 != nil {
		fmt.Println("Error:", pokeErr2)
	}

	trainer.AddToParty(*pokemons[1])
	trainer.AddToParty(*pokemons[2])
	trainer.AddToParty(*pokemons[2])
	trainer.AddToParty(*pokemons[2])
	trainer.AddToParty(*pokemons[2])
	errAddtoParty := trainer.AddToParty(*pokemons[2])

	if errAddtoParty != nil {
		fmt.Println(errAddtoParty)
	}
	trainer.Infos()

	damageFireXGrass := moves[0].CalculateDamage("grass")
	damageWaterXFire := moves[1].CalculateDamage("fire")
	damageGrassXFire := moves[2].CalculateDamage("fire")

	fmt.Printf("Damage to Fire X Grass: %d\n", damageFireXGrass)
	fmt.Printf("Damage to Water X Fire: %d\n", damageWaterXFire)
	fmt.Printf("Damage to Grass X Fire: %d\n", damageGrassXFire)

	var pokedex pd.Pokedex

	pokedex.AddPokemon(1, *pokemons[0])
	pokedex.AddPokemon(2, *pokemons[1])
	pokedex.AddPokemon(3, *pokemons[2])

	fmt.Println(pokedex.ListAll())

	getPokemonToPokedex1, status1, err1 := pokedex.GetPokemon(1)
	getPokemonToPokedex2, status2, err2 := pokedex.GetPokemon(4)

	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println("Pokemon entrontado: ", getPokemonToPokedex1, "encontrou?", status1)

	}

	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println("Pokemon entrontado: ", getPokemonToPokedex2, "encontrou?", status2)
	}

}

func createPokemons() []*p.Pokemon {

	var listPokemons []*p.Pokemon

	pokemon1 := p.Pokemon{
		Name:      "Totodile",
		Type:      []string{"Water"},
		Level:     16,
		EvolvesTo: "Croconaw",
	}

	pokemon2 := p.Pokemon{
		Name:      "Staryu",
		Type:      []string{"Water"},
		Level:     10,
		EvolvesTo: "Starmie",
	}

	pokemon3 := p.Pokemon{
		Name:      "Swampert",
		Type:      []string{"Water", "Ground"},
		Level:     36,
		EvolvesTo: "",
	}
	listPokemons = append(listPokemons, &pokemon1)
	listPokemons = append(listPokemons, &pokemon2)
	listPokemons = append(listPokemons, &pokemon3)
	return listPokemons
}

func generateInfos[T TypeInformation](list []T) {
	for _, item := range list {
		item.Infos()
	}
}

func createMoves() []*m.Move {
	var moves []*m.Move
	move1 := m.Move{
		Name:  "Fireblast",
		Power: 120,
		Type:  "fire",
	}

	move2 := m.Move{
		Name:  "Hydro Bomb",
		Power: 100,
		Type:  "water",
	}

	move3 := m.Move{
		Name:  "Giga drain",
		Power: 90,
		Type:  "grass",
	}

	moves = append(moves, &move1)
	moves = append(moves, &move2)
	moves = append(moves, &move3)
	return moves
}
