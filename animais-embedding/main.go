package main

import (
	"example.com/animais-embedding/animais"
)

func main() {
	cachorro := animais.Cachorro{Animal: animais.Animal{Nome: "Cookies"}}
	gato := animais.Gato{Animal: animais.Animal{Nome: "Nyx"}}

	cachorro.Latir()
	cachorro.Comer()

	gato.Miar()
	gato.Comer()
}
