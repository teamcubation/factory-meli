package main

import (
	"Yuri/exercises-yuri/abstracao/interfaces"
	"Yuri/exercises-yuri/abstracao/models"
	"fmt"
)

func main() {

	quadrado := models.Quadrado{Lado: 5}
	circulo := models.Circulo{Raio: 50}
	fmt.Println(quadrado)
	fmt.Println(circulo)
	ImprimirMetrics(quadrado)
	ImprimirMetrics(circulo)

	// --------------------------------

	consoleLogger := models.ConsoleLogger{}
	fileLogger := models.FileLogger{}

	ProcessarTarefa(consoleLogger)
	ProcessarTarefa(fileLogger)

}

func ImprimirMetrics(f interfaces.IForma) {
	fmt.Println("Area: ", f.Area())
	fmt.Println("Perimetro: ", f.Perimetro())
}

func ProcessarTarefa(logger interfaces.ILogger) {
	logger.Log("TESTE", "teste de tarefa")
}
