package main

import (
	"errors"
	"fmt"
)

func main() {

	num := 193
	err, res := isPrimeNumber(num)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Número", num, "é primo?", res)
}

func isPrimeNumber(number int) (error, bool) {
	if number == 0 {
		return errors.New("isPrimeNumber: número 0 não pode ser lido como número primo"), false
	}
	if number == 1 {
		return nil, false
	}
	for i := 2; i*i <= number; i++ {
		if number%i == 0 {
			return nil, false
		}
	}
	return nil, true
}
