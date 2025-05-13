package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	a, b, err := readNumber()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	result := mdc(a, b)

	fmt.Println(result)
}

func readNumber() (int, int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() == false {
		if err := scanner.Err(); err != nil {
			return 0, 0, fmt.Errorf("Entrada encerrada de forma inesperada.")
		}
	}

	input := strings.Fields(scanner.Text())

	if len(input) != 2 {
		return 0, 0, fmt.Errorf("Entrada enviada no formato incorreto. Ex: '1 2'")
	}

	numberA, err := strconv.Atoi(input[0])
	numberB, err := strconv.Atoi(input[1])

	if err != nil {
		return 0, 0, fmt.Errorf("Erro ao converter entrada '%s'. Erro: %w", input, err)
	}

	return numberA, numberB, nil
}

func mdc(a int, b int) int {
	for b != 0 {
		r := a % b
		a = b
		b = r
	}

	return a
}
