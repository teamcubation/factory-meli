package main

import (
	"fmt"
)

type ProcessadorPagamento interface {
	Processar(valor float64) error
}

func FinalizarCompra(p ProcessadorPagamento, valor float64) error {
	err := p.Processar(valor)
	if err != nil {
		return err
	}
	return nil
}

type PayPal struct{}

func (p *PayPal) Processar(valor float64) error {
	if valor <= 0 {
		return fmt.Errorf("Valor do pamento deve ser positivo!")
	}

	fmt.Println("Operação processado pelo PayPal")
	return nil
}

func PayPalFactory() *PayPal {
	return &PayPal{}
}

type Stripe struct{}

func (s *Stripe) Processar(valor float64) error {
	if valor <= 0 {
		return fmt.Errorf("O valor deve ser positivo para o pagamento")
	}
	fmt.Println("Operação processada pelo Stripe")
	return nil
}

func StripeFactory() *Stripe {
	return &Stripe{}
}

func main() {
	paypal := PayPalFactory()
	stripe := StripeFactory()

	FinalizarCompra(paypal, 20)
	FinalizarCompra(stripe, 100)

}
