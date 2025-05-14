package main

import "fmt"

type Animal struct {
	Nome string
}

func (a Animal) Comer() {
	fmt.Printf("%s acabou de se alimentar! \n", a.Nome)
}

type Gato struct {
	Animal
}

func (g Gato) Miar() {
	fmt.Println("Meowwwwwwn")
}

func GatoFactory(nome string) *Gato {
	return &Gato{Animal{nome}}
}

type Cachorro struct {
	Animal
}

func (c Cachorro) Latir() {
	fmt.Println("Au Au au!")
}

func CachorroFactory(nome string) *Cachorro {
	return &Cachorro{Animal{nome}}
}

func main() {
	cachorro := CachorroFactory("caramelo")
	cachorro.Comer()
	cachorro.Latir()

	gato := GatoFactory("tom")
	gato.Comer()
	gato.Miar()
}
