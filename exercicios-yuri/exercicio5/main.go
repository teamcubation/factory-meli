package main

import (
	"errors"
	"fmt"
)

func main() {
	num := 12345
	soma, err := somaDigitos(num)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("A soma dos dígitos de", num, "é:", soma)

}

func somaDigitos(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("somaDigitos: Número negativo não pode ser lido como soma de dígitos")
	}
	soma := 0
	for n > 0 {
		soma += n % 10
		fmt.Println("Soma:", soma)
		n /= 10
		fmt.Println("N:", n)

	}
	return soma, nil
}
