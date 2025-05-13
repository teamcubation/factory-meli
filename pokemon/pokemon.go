package main

import "fmt"

type Pokemon struct {
	Name      string
	Type      []string
	Level     int
	EvolvesTo string
}

type Trainer struct {
	Name  string
	Party []Pokemon
}

type Move struct {
	Name  string
	Power int
	Type  string
}

type Pokedex map[int]Pokemon

func main() {
	trainer := Trainer{Name: "Ash"}

	pokedex := Pokedex{}

	pokemon1 := Pokemon{Name: "Bulbasaur", Type: []string{"Planta", "Venenoso"}, Level: 17, EvolvesTo: "Ivysaur"}
	pokemon2 := Pokemon{Name: "Charmander", Type: []string{"Fogo"}, Level: 16, EvolvesTo: "Charmeleon"}
	pokemon3 := Pokemon{Name: "Squirtle", Type: []string{"Água"}, Level: 18, EvolvesTo: "Wartortle"}
	pokemon4 := Pokemon{Name: "Pikachu", Type: []string{"Elétrico"}, Level: 22, EvolvesTo: "Raichu"}
	pokemon5 := Pokemon{Name: "Eevee", Type: []string{"Normal"}, Level: 20, EvolvesTo: "Vaporeon"}
	pokemon6 := Pokemon{Name: "Machop", Type: []string{"Lutador"}, Level: 19, EvolvesTo: "Machoke"}

	move1 := Move{Name: "Ember", Power: 40, Type: "Fogo"}
	move2 := Move{Name: "Water Gun", Power: 40, Type: "Água"}
	move3 := Move{Name: "Vine Whip", Power: 45, Type: "Planta"}

	// 1 - Evolve pokemon
	fmt.Println("Evolve:")
	err := Evolve(&pokemon1)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = Evolve(&pokemon2)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 2 - Manage team
	err = trainer.AddToParty(pokemon3)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = trainer.AddToParty(pokemon4)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 3 - Calculate damages
	fmt.Println("\nMoves:")
	fmt.Printf("Move %v of type %v deals %v damage against 'Grass' type.\n", move1.Name, move1.Type, calculateDamage(move1, "Planta"))
	fmt.Printf("Move %v of type %v deals %v damage against 'Fire' type.\n", move2.Name, move2.Type, calculateDamage(move2, "Fogo"))
	fmt.Printf("Move %v of type %v deals %v damage against 'Normal' type.\n", move3.Name, move3.Type, calculateDamage(move3, "Normal"))

	// 4 - Pokedex
	AddPokemon(pokedex, 5, pokemon5)
	AddPokemon(pokedex, 6, pokemon6)

	fmt.Println("\nPokedex:")

	searchedPokemon, found := GetPokemon(pokedex, 5)

	if found {
		fmt.Printf("Name: %v - Type: %v\n", searchedPokemon.Name, searchedPokemon.Type)
	} else {
		fmt.Println("Pokemon not found")
	}

	searchedPokemon, found = GetPokemon(pokedex, 6)

	if found {
		fmt.Printf("Name: %v - Type: %v\n", searchedPokemon.Name, searchedPokemon.Type)
	} else {
		fmt.Println("Pokemon not found")
	}

	allPokemon := ListAll(pokedex)

	fmt.Println("\nList of Pokemon:")
	for _, pokemon := range allPokemon {
		fmt.Printf("- %v\n", pokemon.Name)
	}
}

func Evolve(pokemon *Pokemon) error {
	if pokemon.EvolvesTo == "" {
		return fmt.Errorf("Pokemon nao possui uma evolução")
	}

	if pokemon.Level < 16 {
		return fmt.Errorf("Pokemon precisa estar level 16 ou mais para evoluir. Level atual: %d", pokemon.Level)
	}

	fmt.Printf("%s evolui para %s\n", pokemon.Name, pokemon.EvolvesTo)
	pokemon.Name = pokemon.EvolvesTo
	pokemon.EvolvesTo = ""

	return nil
}

func (trainer Trainer) AddToParty(pokemon Pokemon) error {
	if len(trainer.Party) >= 6 {
		return fmt.Errorf("Time ja possui 6 pokemons.")
	}

	trainer.Party = append(trainer.Party, pokemon)

	return nil
}

func calculateDamage(move Move, targetType string) int {
	if move.Type == "Fogo" && targetType == "Planta" {
		return move.Power * 2
	} else if move.Type == "Água" && targetType == "Fogo" {
		return move.Power * 2
	} else {
		return move.Power
	}
}

func AddPokemon(pokedex Pokedex, id int, pokemon Pokemon) {
	pokedex[id] = pokemon
}

func GetPokemon(pokedex Pokedex, id int) (Pokemon, bool) {
	pokemon, exists := pokedex[id]

	return pokemon, exists
}

func ListAll(pokedex Pokedex) []Pokemon {
	list := []Pokemon{}

	for _, pokemon := range pokedex {
		list = append(list, pokemon)
	}

	return list
}
