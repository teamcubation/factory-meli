package utils

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func StringToSliceInt(numString string, hasVirgula bool) ([]int, error) {
	numArray := []int{}
	var numStrs []string

	if hasVirgula {
		numStrs = strings.Split(strings.TrimSpace(numString), ",")
	} else {
		numStrs = strings.Split(strings.TrimSpace(numString), "")
	}

	for _, numStr := range numStrs {
		num, err := strconv.Atoi(strings.TrimSpace(numStr))
		if err != nil {
			return []int{}, errors.New("erro ao converter " + numStr + " para número: " + err.Error() + "\n")
		}

		numArray = append(numArray, num)
	}

	return numArray, nil
}

func StringToInt(numString string) (int, error) {
	numInt, err := strconv.Atoi(strings.TrimSpace(numString))
	if err != nil {
		return 0, errors.New("erro ao converter string para int")
	}

	return numInt, nil
}

func CalcularNumerosPares(numeros []int) int {
	resultado := 0
	for _, num := range numeros {
		if num%2 == 0 {
			resultado += num
		}
	}

	return resultado
}

func CalcularFatorial(numero int) (int, error) {
	if numero < 0 {
		return 0, errors.New("número precisa ser maior que zero")
	}

	if numero == 0 || numero == 1 {
		return 1, nil
	}

	resultado := 1
	for i := numero; i >= 1; i-- {
		resultado *= i
	}

	return resultado, nil
}

func EncontrarMDC(numeros []int) (int, error) {
	if len(numeros) != 2 {
		return 0, errors.New("é necessário fornecer no máximo 2 números")
	}

	if numeros[0] <= 0 || numeros[1] <= 0 {
		return 0, errors.New("ambos os números devem ser maior que 0")
	}

	// Algoritmo de Euclides
	num1, num2 := numeros[0], numeros[1]
	for num2 != 0 {
		num1, num2 = num2, num1%num2
	}

	return num1, nil
}

func GerarSequenciaFibonacci(indiceMaximo int) ([]int, error) {
	if indiceMaximo <= 0 {
		return nil, errors.New("o indice deve ser maior que zero (0)")
	}

	escalaFibonacci := []int{0, 1}
	for i := 2; i < indiceMaximo; i++ {
		resultado := escalaFibonacci[i-1] + escalaFibonacci[i-2]
		escalaFibonacci = append(escalaFibonacci, resultado)
	}

	return escalaFibonacci, nil
}

func SomaDigitos(numero string) (int, error) {
	numeros, err := StringToSliceInt(numero, false)
	if err != nil {
		return 0, err
	}

	resultado := 0
	for _, value := range numeros {
		resultado += value
	}

	return resultado, nil
}

func IsPalindromo(input string) bool {
	var limpo strings.Builder
	limpo.Grow(len(input))

	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			limpo.WriteRune(unicode.ToLower(char))
		}
	}

	runas := []rune(limpo.String())
	tamanho := len(runas)

	for i := range tamanho / 2 {
		if runas[i] != runas[tamanho-1-i] {
			return false
		}
	}

	return true
}
