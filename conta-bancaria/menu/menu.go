package menu

import "fmt"

type MenuOptions int

const (
	Deposito MenuOptions = iota
	Saldo
	Saque
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
	case Deposito:
		return "Deposito"
	case Saldo:
		return "Saldo"
	case Saque:
		return "Saque"
	case Sair:
		return "Sair"
	default:
		return ""
	}
}

func AllOptions() []MenuOptions {
	return []MenuOptions{Deposito, Saldo, Saque, Sair}
}
