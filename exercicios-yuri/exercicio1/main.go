package main

func main() {
	num := 15
	soma := somaPares(num)
	println("A soma dos números pares de 0 a", num, "é:", soma)
}

func somaPares(n int) int {
	soma := 0
	for i := 0; i <= n; i++ {
		if i%2 == 0 {
			soma += i
		}
	}
	return soma
}
