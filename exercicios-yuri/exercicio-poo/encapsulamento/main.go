package main

import (
	"Yuri/exercises-yuri/encapsulamento/model"
	"fmt"
)

func main() {

	conta := model.Conta{}
	conta.SetTitular("Yuri")
	conta.Depositar(500)
	fmt.Println(conta)
	errDeposito := conta.Depositar(-10)
	if errDeposito != nil {
		fmt.Println(errDeposito)
	}

	conta.Sacar(100)
	fmt.Println(conta)

	errSacar := conta.Sacar(500)
	if errSacar != nil {
		fmt.Println(errSacar)
	}
	fmt.Println(conta)

}
