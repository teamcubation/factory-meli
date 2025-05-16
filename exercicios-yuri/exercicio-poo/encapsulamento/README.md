## 1. Encapsulamento

**1.1 Conta Bancária**

* Crie um `struct` `Conta` com campos não exportados `titular string` e `saldo float64`.
* Implemente métodos públicos:

  * `Depositar(valor float64)` — só adiciona se `valor > 0`.
  * `Sacar(valor float64) error` — retorna erro se não houver saldo suficiente.
  * `Saldo() float64` — retorna o saldo atual.
* Verifique que de outro pacote você **não** consegue acessar `conta.saldo` diretamente.

**1.2 Retângulo Validado**

* Defina um `struct` `Retangulo` com campos privados `largura, altura float64`.
* Crie setters públicos:

  * `SetLargura(w float64) error`
  * `SetAltura(h float64) error`
    que só aceitam valores positivos.
* Adicione métodos:

  * `Area() float64`
  * `Perimetro() float64`
* Teste criando um retângulo com valores negativos e assegure-se de receber erro no setter.

---