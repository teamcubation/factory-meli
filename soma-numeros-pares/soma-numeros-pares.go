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

	sum := sumEvenNumbers(number)

	fmt.Println(sum)
}

func sumEvenNumbers(number int) int {
	sum := 0

	for i := 0; i <= number; i++ {
		if i%2 == 0 {
			sum += i
		}
	}

	return sum
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
