package conta

import "fmt"

type Conta struct {
	titular string
	saldo   float64
}

func (c *Conta) Depositar(valor float64) error {
	if valor <= 0 {
		return fmt.Errorf("Valor de deposito deve ser maior que 0. Valor informado: %v", valor)
	}

	c.saldo += valor

	return nil
}

func (c *Conta) Sacar(valor float64) error {
	if c.saldo < valor {
		return fmt.Errorf("Voce nao possui valor em conta o suficiente para efetuar esse saque.")
	}

	c.saldo -= valor

	return nil
}

func (c *Conta) Saldo() float64 {
	return c.saldo
}
