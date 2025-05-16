package main

import (
	"errors"
	"fmt"
)

func main() {
	num := -3
	fatorial, err := fatorial(num)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("O fatorial de", num, "é:", fatorial)
}

func fatorial(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("fatorial: Número negativo não pode ser lido como fatorial")
	}
	if n == 0 {
		return 0, nil
	}
	fatorial := 1
	for i := 1; i <= n; i++ {
		fatorial *= i
	}
	return fatorial, nil
}
