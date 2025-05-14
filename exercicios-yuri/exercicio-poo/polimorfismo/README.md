
## 4. Polimorfismo

**4.1 Processadores de Pagamento**

* Declare interface `ProcessadorPagamento` com método `Processar(valor float64) error`.
* Implemente `PayPal` e `Stripe` com saídas diferentes no `Processar`.
* Escreva função `FinalizarCompra(p ProcessadorPagamento, valor float64)` que chame `p.Processar(valor)`.

**4.2 Notificações**

* Defina interface `Notificador` com `Notificar(destinatario, mensagem string) error`.
* Implemente `EmailNotificador` e `SMSNotificador` (simule com `fmt.Printf`).
* Crie um slice de `Notificador` e envie uma mensagem para todos em um loop.