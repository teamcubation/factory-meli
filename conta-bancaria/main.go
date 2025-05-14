package main

import (
	"fmt"
	"os"

	"example.com/conta-bancaria/conta"
	"example.com/conta-bancaria/menu"
)

func main() {
	contaBancaria := conta.Conta{}

	for {
		var userInput menu.MenuOptions

		menu.RenderMenu()
		_, err := fmt.Scan(&userInput)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Erro na leitura do item escolhido.")
			break
		}

		if userInput == menu.Sair {
			break
		}

		switch userInput {
		case menu.Deposito:
			var valor float64
			fmt.Print("Digite o valor que deseja depositar: ")
			fmt.Scan(&valor)

			err := contaBancaria.Depositar(valor)

			if err != nil {
				fmt.Fprintln(os.Stderr, "Erro ao realizar deposito. Motivo: "+err.Error())
			}
			continue

		case menu.Saldo:
			saldo := contaBancaria.Saldo()
			fmt.Printf("Saldo atual: %.2f\n", saldo)
			continue

		case menu.Saque:
			var valor float64
			fmt.Print("Digite o valor que deseja sacar: ")
			fmt.Scan(&valor)

			err := contaBancaria.Sacar(valor)

			if err != nil {
				fmt.Fprintln(os.Stderr, "Erro ao realizar o saque. Motivo: "+err.Error())
			}
			continue

		default:
			fmt.Printf("Opção '%d' nao existe.\n", userInput)
			continue
		}

	}

	fmt.Print("Programa encerrado")
}
