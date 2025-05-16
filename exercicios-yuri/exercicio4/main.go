package main

import "fmt"

func main() {
	num := 7
	fib, err := fibonnacci(num)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Os números de Fibonacci até", num, "são:", fib)

}

func fibonnacci(n int) ([]int, error) {
	if n < 0 {
		return nil, fmt.Errorf("fibonacci: Número negativo não pode ser lido como número de Fibonacci")
	}
	if n == 0 {
		return []int{0}, nil
	}
	if n == 1 {
		return []int{0, 1}, nil
	}
	resultado := []int{0, 1}
	for i := 2; i <= n; i++ {
		proximo := resultado[i-1] + resultado[i-2]
		resultado = append(resultado, proximo)
	}
	return resultado, nil
}
