package main

import (
	"fmt"
)

type Pokemon struct {
	Name      string
	Type      []string
	Level     int
	EvolvesTo string
}

func Evolve(p *Pokemon) error {
	if p.Level < 16 {
		return fmt.Errorf("O %s ainda não está no nível 16 para poder evoluir!", p.Name)
	}
	fmt.Printf("¡%s evolui para %s!\n", p.Name, p.EvolvesTo)
	p.Name = p.EvolvesTo
	return nil
}

func PokemonFactory(name string, types []string, level int, evolvesTo string) *Pokemon {
	return &Pokemon{name, types, level, evolvesTo}
}

type Trainer struct {
	Name  string
	Party []Pokemon
}

func (t *Trainer) AddToParty(p Pokemon) error {
	if len(t.Party) >= 6 {
		return fmt.Errorf("o time já está cheio. Não foi possível adicionar mais um pokemon")
	}
	t.Party = append(t.Party, p)
	return nil
}

func TrainerFactory(name string) *Trainer {
	return &Trainer{name, []Pokemon{}}
}

type Move struct {
	Name  string
	Power int
	Type  string
}

func CalculateDamage(m Move, targetType string) int {
	damage := m.Power
	if (m.Type == "Fogo" && targetType == "Planta") ||
		(m.Type == "Água" && targetType == "Fogo") {
		damage *= 2
	}
	return damage
}

type Pokedex map[int]Pokemon

func AddPokemon(pokedex Pokedex, id int, p Pokemon) error {
	if _, ok := pokedex[id]; ok {
		return fmt.Errorf("O Pokémon com ID %d já está na Pokedex", id)
	}
	pokedex[id] = p
	return nil
}

func GetPokemon(pokedex Pokedex, id int) (Pokemon, bool) {
	pokemon, exists := pokedex[id]
	return pokemon, exists
}

func ListAll(pokedex Pokedex) []Pokemon {
	pokemons := []Pokemon{}
	for _, pokemon := range pokedex {
		pokemons = append(pokemons, pokemon)
	}
	return pokemons
}

func main() {
	// Criando Pokémons
	charmander := PokemonFactory("Charmander", []string{"Fogo"}, 16, "Charmeleon")
	bulbasaur := PokemonFactory("Bulbasaur", []string{"Planta"}, 10, "Ivysaur")
	squirtle := PokemonFactory("Squirtle", []string{"Água"}, 16, "Wartortle")

	// Testando evolução
	if err := Evolve(charmander); err != nil {
		fmt.Println("Erro:", err)
	}
	if err := Evolve(bulbasaur); err != nil {
		fmt.Println("Erro:", err)
	}

	// Criando treinador
	ash := TrainerFactory("Ash")

	// Adicionando à Party
	for _, p := range []Pokemon{*charmander, *bulbasaur, *squirtle, *squirtle, *squirtle, *squirtle, *squirtle} {
		err := ash.AddToParty(p)
		if err != nil {
			fmt.Println("Erro:", err)
		}
	}
	fmt.Println("Party de", ash.Name, ":", ash.Party)

	// Criando movimentos
	ember := Move{Name: "Ember", Power: 40, Type: "Fogo"}
	waterGun := Move{Name: "Water Gun", Power: 40, Type: "Água"}

	fmt.Println("Dano Ember contra Planta:", CalculateDamage(ember, "Planta"))
	fmt.Println("Dano Water Gun contra Fogo:", CalculateDamage(waterGun, "Fogo"))
	fmt.Println("Dano Ember contra Água:", CalculateDamage(ember, "Água"))

	// Pokédex
	pokedex := make(Pokedex)
	AddPokemon(pokedex, 1, *charmander)
	AddPokemon(pokedex, 2, *bulbasaur)

	// Tentando adicionar duplicado
	if err := AddPokemon(pokedex, 1, *squirtle); err != nil {
		fmt.Println("Erro:", err)
	}

	// Buscar Pokémon
	if poke, ok := GetPokemon(pokedex, 2); ok {
		fmt.Println("Encontrado na Pokédex:", poke)
	} else {
		fmt.Println("Pokémon não encontrado na Pokédex")
	}

	// Listar todos
	fmt.Println("Todos os Pokémons na Pokédex:")
	for _, p := range ListAll(pokedex) {
		fmt.Println("-", p.Name)
	}
}
