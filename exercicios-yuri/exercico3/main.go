package main

import (
	"fmt"
	"slices"
)

func main() {
	num1 := 17
	num2 := 2
	mdc, err := mdc(num1, num2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("O MDC de", num1, "e", num2, "é:", mdc)
}

func mdc(a, b int) (int, error) {
	if a == 0 && b == 0 {
		return 0, fmt.Errorf("mdc: Não foi possível encontrar o MDC entre %d e %d", a, b)

	}
	divisoresA := divisores(a)
	divisoresB := divisores(b)
	for i := len(divisoresA) - 1; i >= 0; i-- {
		if slices.Contains(divisoresB, divisoresA[i]) {
			return divisoresA[i], nil
		}
	}
	return maiorNumero(a, b), nil
}

func divisores(n int) []int {
	divisores := []int{}
	for i := 1; i <= n; i++ {
		if n%i == 0 {
			divisores = append(divisores, i)
		}
	}
	return divisores
}

func maiorNumero(a, b int) int {
	if a > b {
		return a
	}
	return b
}
