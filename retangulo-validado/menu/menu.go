package menu

import "fmt"

type MenuOptions int

const (
	SetAltura MenuOptions = iota
	SetLargura
	Area
	Perimetro
	Sair = 10
)

func RenderMenu() {
	fmt.Println("---------------------")
	fmt.Print("Menu:\n\n")

	for _, option := range AllOptions() {
		fmt.Printf("%d - %s\n", option, option)
	}
	fmt.Println("---------------------")
	fmt.Print("Digite uma opcao: ")
}

func (option MenuOptions) String() string {
	switch option {
	case SetAltura:
		return "Setar Altura"
	case SetLargura:
		return "Setar largura"
	case Area:
		return "Area do retângulo"
	case Perimetro:
		return "Perimetro do retângulo"
	case Sair:
		return "Sair"
	default:
		return ""
	}
}

func AllOptions() []MenuOptions {
	return []MenuOptions{SetAltura, SetLargura, Area, Perimetro, Sair}
}
