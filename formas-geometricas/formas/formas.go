package formas

import (
	"fmt"
	"math"
)

type Forma interface {
	Area() float64
	Perimetro() float64
}

type Circulo struct {
	Raio float64
}

func (c Circulo) Area() float64 {
	return math.Pi * math.Pow(c.Raio, 2)
}

func (c Circulo) Perimetro() float64 {
	return 2 * math.Pi * c.Raio
}

type Quadrado struct {
	Lado float64
}

func (q Quadrado) Area() float64 {
	return q.Lado * q.Lado
}

func (q Quadrado) Perimetro() float64 {
	return 4 * q.Lado
}

func ImprimeMetrics(f Forma) {
	fmt.Printf("Area: %.2f\n", f.Area())
	fmt.Printf("Perimetro: %.2f\n", f.Perimetro())
}
