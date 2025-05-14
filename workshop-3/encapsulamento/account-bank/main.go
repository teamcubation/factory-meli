package main

import "fmt"

type Account struct {
	titular string
	saldo   float64
}

func (a *Account) Depositar(valor float64) error {
	if valor <= 0 {
		return fmt.Errorf("O valor depositado deve ser positivo.")
	}
	a.saldo += valor
	return nil
}

func (a *Account) Sacar(valor float64) error {
	if valor <= 0 {
		return fmt.Errorf("O valor a ser sacado dever ser positivo.")
	}
	if valor > a.saldo {
		return fmt.Errorf("O valor requerido no saque é maior que o saldo em conta.")
	}
	a.saldo -= valor
	return nil
}

func AccountFactory(titular string, saldo float64) *Account {
	return &Account{titular, saldo}
}

func (a *Account) Saldo() float64 {
	return a.saldo
}

func main() {

	// account := AccountFactory("Rafael Medeiros", 0)
	// fmt.Printf("Conta criada para o usuário: %s, com o valor inicial de: %f\n", account.titular, account.saldo)

	// //Deposita um novo valor
	// err := account.Depositar(10.50)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(account.Saldo())

	// err = account.Sacar(2)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(account.Saldo())

	// err = account.Sacar(5)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(account.Saldo())

}
