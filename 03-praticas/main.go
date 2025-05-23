package main

import (
	utils "03-praticas/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		exibirMenu()

		opcao, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler opção. Tente novamente")
			continue
		}
		opcao = strings.TrimSpace(opcao)

		switch opcao {
		case "0":
			fmt.Println("Saindo... Até logo!")
			return
		case "1":
			fmt.Println("Digite os números separados por vírgulas: ")
			numeros, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				continue
			}
			numeros = strings.TrimSpace(numeros)

			numArray, err := utils.StringToSliceInt(numeros, true)
			if err != nil {
				fmt.Println(err)
				continue
			}

			resSomaPares := utils.CalcularNumerosPares(numArray)
			fmt.Printf("Resultado da soma dos números pares: %d \n", resSomaPares)
		case "2":
			fmt.Println("Digite o número: ")
			numero, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				continue
			}

			numeroInt, err := utils.StringToInt(numero)
			if err != nil {
				fmt.Println(err)
				continue
			}

			resFatorial, err := utils.CalcularFatorial(numeroInt)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Resultado do fatorial: %d \n", resFatorial)
		case "3":
			fmt.Println("Digite os dois números serparados por vírgula: ")
			numeros, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				continue
			}

			numerosSlice, err := utils.StringToSliceInt(numeros, true)
			if err != nil {
				fmt.Println(err)
				continue
			}

			resMDC, err := utils.EncontrarMDC(numerosSlice)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Resultado do MDC %d \n", resMDC)
		case "4":
			fmt.Println("Digite o indice máximo: ")
			numero, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				continue
			}

			numeroInt, err := utils.StringToInt(numero)
			if err != nil {
				fmt.Println(err)
				continue
			}

			resFibonacci, err := utils.GerarSequenciaFibonacci(numeroInt)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Resultado do Fibonacci %d \n", resFibonacci)
		case "5":
			fmt.Println("Digite um número: ")
			numero, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				continue
			}

			resultado, err := utils.SomaDigitos(numero)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Resultado da soma de dígitos: %d \n", resultado)
		case "6":
			fmt.Println("Digite uma palavra ou frase: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				continue
			}

			if utils.IsPalindromo(input) {
				fmt.Printf("A palavra %s é um palíndromo! \n", strings.TrimSpace(input))
			} else {
				fmt.Printf("A palavra %s não é um palíndromo! \n", strings.TrimSpace(input))
			}
		default:
			fmt.Println("Opção não encontrada. Tente novamente.")
		}
	}
}

func exibirMenu() {
	fmt.Println("+----------------------------------------------------------+")
	fmt.Println("|                           MENU                           |")
	fmt.Println("+----------------------------------------------------------+")
	fmt.Println("| [1] Calcular números pares                               |")
	fmt.Println("| [2] Calcular fatorial                                    |")
	fmt.Println("| [3] Encontrar o MDC de dois números                      |")
	fmt.Println("| [4] Gerar sequência Fibonacci                            |")
	fmt.Println("| [5] Somar dígitos de um número                           |")
	fmt.Println("| [6] Verificar palíndromo                                 |")
	fmt.Println("| [0] Sair                                                 |")
	fmt.Println("+----------------------------------------------------------+")
	fmt.Print("Escolha uma opção: ")
}
