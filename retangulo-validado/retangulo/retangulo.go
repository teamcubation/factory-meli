package retangulo

import "fmt"

type Retangulo struct {
	largura float64
	altura  float64
}

func (r *Retangulo) SetLargura(largura float64) error {
	if largura < 0 {
		return fmt.Errorf("Apenas valores positivos sao aceitos. Valor informado: %.2f", largura)
	}

	r.largura = largura

	return nil
}

func (r *Retangulo) SetAltura(altura float64) error {
	if altura < 0 {
		return fmt.Errorf("Apenas valores positivos sao aceitos. Valor informado: %.2f", altura)
	}

	r.altura = altura

	return nil
}

func (r Retangulo) Area() float64 {
	return r.altura * r.largura
}

func (r Retangulo) Perimetro() float64 {
	return 2 * (r.altura + r.largura)
}
