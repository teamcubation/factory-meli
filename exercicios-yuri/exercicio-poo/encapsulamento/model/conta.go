package model

import "errors"

type Conta struct {
	titular string
	saldo   float64
}

func (c *Conta) Depositar(valor float64) error {
	if valor < 0 {
		return errors.New("Depositar: não pode usar valores negativo")
	}
	c.saldo += valor
	return nil
}

func (c *Conta) Sacar(valor float64) error {
	if valor < 0 || c.saldo < valor {
		return errors.New("Sacar: saldo insuficiente")
	}
	c.saldo -= valor
	return nil
}

func (c *Conta) Saldo() float64 {
	return c.saldo
}

func (c *Conta) SetTitular(nome string) {
	c.titular = nome
}

func (c *Conta) GetTitular() string {
	return c.titular
}
