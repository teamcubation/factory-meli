package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func obtemTotalNuvens(scanner *bufio.Scanner) (int, string) {
	scanner.Scan()

	totalNuvens, err := strconv.Atoi(scanner.Text())

	if err != nil {
		return -1, "Erro ao ler o numero de nuvens:" + err.Error()
	}

	return totalNuvens, ""
}

func obtemCaminhoNuvens(scanner *bufio.Scanner, numeroNuvens int) ([]int, string) {
	scanner.Scan()

	caminhoString := strings.Fields(scanner.Text())

	if len(caminhoString) != numeroNuvens {
		return nil, "Quantidade de nuvens no caminho nao bate com o numero informado anteriormente"
	}

	caminho := make([]int, numeroNuvens)
	var err error

	for i, x := range caminhoString {
		caminho[i], err = strconv.Atoi(x)

		if err != nil {
			return nil, "Erro ao converter uma das nuvens para int:" + err.Error()
		}
	}

	return caminho, ""
}

func saltarEmNuvens(caminho []int) int {
	saltos := 0
	contador := 0

	for contador < len(caminho)-1 {
		if contador+2 < len(caminho) && caminho[contador+2] == 0 {
			contador += 2
		} else {
			contador += 1
		}

		saltos++
	}

	return saltos
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	totalNuvens, errTotalNuvens := obtemTotalNuvens(scanner)
	if errTotalNuvens != "" {
		fmt.Println(errTotalNuvens)
		return
	}

	caminho, errCaminhoNuvens := obtemCaminhoNuvens(scanner, totalNuvens)
	if errCaminhoNuvens != "" {
		fmt.Println(errCaminhoNuvens)
		return
	}

	totalSaltos := saltarEmNuvens(caminho)
	fmt.Println(totalSaltos)
}
