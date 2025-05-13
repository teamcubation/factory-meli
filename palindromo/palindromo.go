package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	word, err := readInput()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if isPalindrome(word) {
		fmt.Printf("%s é um palíndromo", word)
	} else {
		fmt.Printf("%s não é um palíndromo", word)
	}
}

func readInput() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() == false {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("Entrada encerrada de forma inesperada.")
		}
	}

	input := strings.Trim(scanner.Text(), " ")

	return input, nil
}

func isPalindrome(word string) bool {
	reversedWord := reverseWord(word)
	return word == reversedWord
}

func reverseWord(word string) string {
	runes := []rune(word)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
