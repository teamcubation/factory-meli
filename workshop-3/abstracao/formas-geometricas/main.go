package main

import (
	"fmt"
	"math"
)

type Forma interface {
	Area() float64
	Perimetro() float64
}

func ImprimirMetrics(f Forma) {
	fmt.Printf("A forma possui uma area de: %f, e um perimetro de %f \n", f.Area(), f.Perimetro())
}

type Circulo struct {
	raio float64
}

func (c Circulo) Area() float64 {
	return math.Pi * math.Pow(c.raio, 2)
}

func (c Circulo) Perimetro() float64 {
	return 2 * math.Pi * c.raio
}

func CirculoFactory(raio float64) *Circulo {
	return &Circulo{raio}
}

func (q Quadrado) Area() float64 {
	return q.lado * q.lado
}

func (q Quadrado) Perimetro() float64 {
	return 4 * q.lado
}

type Quadrado struct {
	lado float64
}

func QuadradoFactory(lado float64) *Quadrado {
	return &Quadrado{lado}
}

func main() {

	quadrado := QuadradoFactory(5)
	circulo := CirculoFactory(3)

	ImprimirMetrics(quadrado)
	ImprimirMetrics(circulo)

}
