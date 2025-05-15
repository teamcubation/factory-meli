package animais

import "fmt"

type Animal struct {
	Nome string
}

func (animal Animal) Comer() {
	fmt.Println(animal.Nome, "esta comendo")
}

type Cachorro struct {
	Animal
}

func (cachorro Cachorro) Latir() {
	fmt.Println(cachorro.Nome, "esta latindo")
}

type Gato struct {
	Animal
}

func (gato Gato) Miar() {
	fmt.Println(gato.Nome, "esta miando")
}
