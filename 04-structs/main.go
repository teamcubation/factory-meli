package main

import (
	"04-structs/models"
	"fmt"
)

func main() {
	pokemon1 := models.Pokemon{
		Name:      "Bulbassaur",
		Types:     []string{"Grass"},
		Level:     1,
		EvolvesTo: "Ivyssaur",
	}

	pokemon2 := models.Pokemon{
		Name:      "Charmander",
		Types:     []string{"Fire"},
		Level:     16,
		EvolvesTo: "Charmeleon",
	}

	pokemon3 := models.Pokemon{
		Name:      "Squirtle",
		Types:     []string{"Water"},
		Level:     17,
		EvolvesTo: "Wartortle",
	}

	pokemon4 := models.Pokemon{
		Name:      "Pikachu",
		Types:     []string{"Electric"},
		Level:     5,
		EvolvesTo: "Raichu",
	}

	pokemon5 := models.Pokemon{
		Name:      "Jigglypuff",
		Types:     []string{"Fairy", "Normal"},
		Level:     3,
		EvolvesTo: "Wigglytuff",
	}

	pokemon6 := models.Pokemon{
		Name:      "Meowth",
		Types:     []string{"Normal"},
		Level:     2,
		EvolvesTo: "Persian",
	}

	pokemon7 := models.Pokemon{
		Name:      "Zé da Manga",
		Types:     []string{"Manga"},
		Level:     18,
		EvolvesTo: "Mango Joh",
	}

	trainer1 := models.Trainer{
		Name:  "Luiz",
		Party: []*models.Pokemon{},
	}

	if err := trainer1.AddToParty(&pokemon1); err != nil {
		fmt.Printf("Error adding %s to party: %v\n", pokemon1.Name, err)
	}
	if err := trainer1.AddToParty(&pokemon2); err != nil {
		fmt.Printf("Error adding %s to party: %v\n", pokemon2.Name, err)
	}
	if err := trainer1.AddToParty(&pokemon3); err != nil {
		fmt.Printf("Error adding %s to party: %v\n", pokemon3.Name, err)
	}
	if err := trainer1.AddToParty(&pokemon4); err != nil {
		fmt.Printf("Error adding %s to party: %v\n", pokemon4.Name, err)
	}
	if err := trainer1.AddToParty(&pokemon5); err != nil {
		fmt.Printf("Error adding %s to party: %v\n", pokemon5.Name, err)
	}
	if err := trainer1.AddToParty(&pokemon6); err != nil {
		fmt.Printf("Error adding %s to party: %v\n", pokemon6.Name, err)
	}
	if err := trainer1.AddToParty(&pokemon7); err != nil {
		fmt.Printf("Error adding %s to party: %v\n", pokemon7.Name, err)
	}

	fmt.Println(trainer1.String())

	// ------------------------------------------------------------------------

	move1 := models.Move{
		Name:  "Flamethrower",
		Power: 90,
		Type:  "Fogo",
	}

	move2 := models.Move{
		Name:  "Hydro Pump",
		Power: 110,
		Type:  "Água",
	}

	move3 := models.Move{
		Name:  "Solar Beam",
		Power: 120,
		Type:  "Planta",
	}

	fmt.Printf("Dano de %s contra Planta: %d \n", move1.Name, move1.CalculateDamage("Planta"))
	fmt.Printf("Dano de %s contra Fogo: %d \n", move2.Name, move2.CalculateDamage("Fogo"))
	fmt.Printf("Dano de %s contra Água: %d \n", move3.Name, move3.CalculateDamage("Água"))

	// ------------------------------------------------------------------------

	pokedex1 := models.Pokedex{}
	pokedex1.AddPokemon(1, &pokemon1)
	pokedex1.AddPokemon(2, &pokemon2)
	pokedex1.AddPokemon(3, &pokemon3)
	pokedex1.AddPokemon(4, &pokemon4)

	pkmns := pokedex1.ListAll()
	if len(pkmns) <= 0 {
		fmt.Println("A Pokedex está vazia")
	}
	fmt.Println("Lista de Pokemons da Pokedex:")
	for _, v := range pkmns {
		fmt.Println(v.String())
	}

	pkm, exists := pokedex1.GetPokemon(3)
	if !exists {
		fmt.Println("Pokemon inexistente na Pokedex")
	}
	fmt.Printf("Pokemon encontrado: \n %v \n", pkm.String())
}
