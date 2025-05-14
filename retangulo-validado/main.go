package main

import (
	"fmt"
	"os"

	"example.com/retangulo-validado/menu"
	retangulo "example.com/retangulo-validado/retangulo"
)

func main() {
	retangulo := retangulo.Retangulo{}

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
		case menu.Area:
			fmt.Printf("Area do retângulo: %.2f\n", retangulo.Area())
			continue

		case menu.Perimetro:
			fmt.Printf("Perimetro do retângulo: %.2f\n", retangulo.Perimetro())
			continue

		case menu.SetAltura:
			var altura float64
			fmt.Print("Digite a altura: ")
			fmt.Scan(&altura)
			err := retangulo.SetAltura(altura)

			if err != nil {
				fmt.Fprintln(os.Stderr, "Erro ao setar a altura do retângulo. Motivo:", err.Error())
			}
			continue

		case menu.SetLargura:
			var largura float64
			fmt.Print("Digite a largura: ")
			fmt.Scan(&largura)
			err := retangulo.SetLargura(largura)

			if err != nil {
				fmt.Fprintln(os.Stderr, "Erro ao setar a largura do retângulo. Motivo:", err.Error())
			}
			continue

		default:
			fmt.Printf("Opção '%d' nao existe.\n", userInput)
			continue
		}
	}
}
