package models

type Quadrado struct {
	Lado float64
}

func (q Quadrado) Area() float64 {
	return q.Lado * q.Lado
}

func (q Quadrado) Perimetro() float64 {
	return q.Lado * 4
}
