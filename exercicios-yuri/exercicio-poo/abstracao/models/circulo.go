package models

const PI = 3.14

type Circulo struct {
	Raio float64
}

func (c Circulo) Area() float64 {
	return PI * (c.Raio * c.Raio)
}

func (c Circulo) Perimetro() float64 {
	return 2 * PI * c.Raio
}
