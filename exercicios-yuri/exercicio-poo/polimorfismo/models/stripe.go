package models

import "errors"

type Stripe struct{}

func (s Stripe) Processar(valor float64) error {
	if valor > 20 {
		return errors.New("Stripe - Processar: valor excedeu o limite")
	}
	return nil
}
