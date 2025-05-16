package models

import "errors"

type PayPal struct{}

func (p PayPal) Processar(valor float64) error {
	if valor > 100 {
		return errors.New("PayPal - Processar: valor excedeu o limite")
	}
	return nil
}
