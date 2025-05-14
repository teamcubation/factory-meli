## 3. Herança (Embedding)

**3.1 Animais com Embedding**

* Crie um `struct` `Animal` com campo `Nome string` e método `Comer()`.
* Defina `Cachorro` e `Gato` que embutem (`embed`) `Animal` e adicionem método próprio `Latir()` ou `Miar()`.
* No `main()`, instancie `Cachorro` e `Gato` e chame `Comer()`, `Latir()`/`Miar()`.

**3.2 Veículos com Embedding**

* Defina `struct` `Veiculo` com campos `Marca string` e `Velocidade int` e método `Ligar()`.
* Crie `Carro` e `Bicicleta` que embutem `Veiculo`, adicionando respectivamente `Portas int` e `TemCampainha bool`.