package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	number, err := readNumber()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if isPrime(number) {
		fmt.Printf("%d is a prime number", number)
	} else {
		fmt.Printf("%d is not a prime number", number)
	}
}

func isPrime(number int) bool {
	if number < 2 {
		return false
	}

	if number == 2 {
		return true
	}

	if number%2 == 0 {
		return false
	}

	sqrt := int(math.Sqrt(float64(number)))

	for i := 3; i <= sqrt; i++ {
		if number%i == 0 {
			return false
		}
	}

	return true
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
