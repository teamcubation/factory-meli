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

	result := fibonacci(number)

	fmt.Println(result)
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

	if number < 0 {
		return 0, fmt.Errorf("Numero precisa ser MAIOR que 0.")
	}

	return number, nil
}

func fibonacci(number int) int {
	if number == 0 {
		return 0
	} else if number == 1 {
		return 1
	} else {
		return fibonacci(number-1) + fibonacci(number-2)
	}
}
