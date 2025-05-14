## 2. Abstração

**2.1 Formas Geométricas**

* Declare a interface `Forma` com métodos `Area() float64` e `Perimetro() float64`.
* Implemente tipos `Circulo{Raio float64}` e `Quadrado{Lado float64}` que satisfaçam `Forma`.
* Escreva uma função `ImprimirMetrics(f Forma)` que exiba área e perímetro de qualquer `Forma`.

**2.2 Sistema de Logging**

* Defina a interface `Logger` com método `Log(nivel, mensagem string)`.
* Implemente `ConsoleLogger` (imprime no stdout) e `FileLogger` (simula escrita em arquivo com `fmt.Println` prefixado por `[FILE]`).
* Crie uma função `ProcessarTarefa(logger Logger)` que use o logger sem saber se é de console ou arquivo.

---
