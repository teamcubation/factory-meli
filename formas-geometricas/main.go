package main

import (
	"example.com/formas-geometricas/formas"
)

func main() {
	circulo := formas.Circulo{Raio: 5}
	quadrado := formas.Quadrado{Lado: 10}

	formas.ImprimeMetrics(circulo)
	formas.ImprimeMetrics(quadrado)
}
