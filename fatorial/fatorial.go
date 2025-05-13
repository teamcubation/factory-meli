package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	number, err := readNumber()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	result := factorial(number)

	fmt.Println(result)
}

func factorial(number int) int {
	result := 1

	for i := 1; i < number; i++ {
		result *= (i + 1)
	}

	return result
}

func readNumber() (int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() == false {
		if err := scanner.Err(); err != nil {
			return 0, fmt.Errorf("Entrada encerrada de forma inesperada.")
		}
	}

	input := scanner.Text()
	number, err := strconv.Atoi(input)

	if err != nil {
		return 0, fmt.Errorf("Erro ao converter entrada '%s'. Erro: %w", input, err)
	}

	return number, nil
}
