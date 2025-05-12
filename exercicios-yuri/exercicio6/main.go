package main

import "fmt"

func main() {
	palavra := "ovo"
	fmt.Println("A palavra", palavra, "é um palíndromo?", palindromo(palavra))
}

func palindromo(palavra string) bool {
	for i, caracter := range palavra {
		if caracter != rune(palavra[len(palavra)-1-i]) {
			return false
		}
	}
	return true
}
