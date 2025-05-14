package main

import "fmt"

type Retangulo struct {
	largura float64
	altura  float64
}

func (r *Retangulo) SetLargura(valor float64) error {
	if valor <= 0 {
		return fmt.Errorf("O valor deve ser positivo")
	}

	r.largura = valor

	return nil
}

func (r *Retangulo) SetAltura(valor float64) error {
	if valor <= 0 {
		return fmt.Errorf("O valor deve ser um número positivo")
	}
	r.altura = valor
	return nil
}

func (r *Retangulo) CalcArea() (float64, error) {
	if r.altura == 0 {
		return -1, fmt.Errorf("Retangulo com valor de altura não iniciado")
	}
	if r.largura == 0 {
		return -1, fmt.Errorf("Retangulo com valor de largura não iniciado")
	}
	return r.altura * r.largura, nil
}

func (r *Retangulo) CalcPerimetro() (float64, error) {
	if r.altura == 0 {
		return -1, fmt.Errorf("Retangulo com o valor de altura não iniciado")
	}
	if r.largura == 0 {
		return -1, fmt.Errorf("Retangulo com o valor de largura não iniciado")
	}

	return (r.altura * 2) + (r.largura * 2), nil

}

func RectanguleFactory(altura, largura float64) *Retangulo {
	return &Retangulo{altura, largura}
}

func main() {
	rec := RectanguleFactory(0, 0)
	area, err := rec.CalcArea()

	// Verificação de erro no canculo da área
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("A area do rectângulo é: %f\n", area)

	// Verificação de erro no calculo set do perimetro
	// err = rec.SetAltura(-1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	err = rec.SetAltura(5)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = rec.SetLargura(2)
	if err != nil {
		fmt.Println(err)
		return
	}

	area, err = rec.CalcArea()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("A area do rectângulo é: %f\n", area)

	perimetro, err := rec.CalcPerimetro()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("O perímetro do rectângulo é: %f\n", perimetro)

}
