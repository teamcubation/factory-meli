package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readInput() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Digite a fase: ")

	initialPhrase, _ := reader.ReadString('\n')
	initialPhrase = strings.TrimSpace(initialPhrase)

	return initialPhrase
}

func reversePhrase(phrase string) string {
	words := strings.Fields(phrase)

	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}

	return strings.Join(words, " ")
}

func reverseWords(phrase string) string {
	newPhrase := ""

	for i := len(phrase) - 1; i >= 0; i-- {
		newPhrase += string(phrase[i])
	}

	return newPhrase
}

func reverseWordsRunes(phrase string) string {
	runes := []rune(phrase)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func main() {
	initialPhrase := readInput()

	fmt.Println("Frase inicial: " + initialPhrase)

	fmt.Println("1: " + reversePhrase(initialPhrase))
	fmt.Println("2: " + reverseWords(initialPhrase))
	fmt.Println("3: " + reverseWordsRunes(initialPhrase))
}
