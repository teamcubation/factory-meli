package main

import (
	"Yuri/exercises-yuri/polimorfismo/interfaces"
	"Yuri/exercises-yuri/polimorfismo/models"
	"fmt"
)

func main() {

	paypal := models.PayPal{}
	stripe := models.Stripe{}

	errP1 := FinalizarCompra(paypal, 30)
	if errP1 != nil {
		fmt.Println(errP1)
	} else {
		fmt.Println("Primeiro pagamento via Paypal foi sucesso")
	}
	errP2 := FinalizarCompra(paypal, 150)
	if errP2 != nil {
		fmt.Println(errP2)
	} else {
		fmt.Println("Segundo pagamento via Paypal foi sucesso")
	}
	errS1 := FinalizarCompra(stripe, 30)
	if errS1 != nil {
		fmt.Println(errS1)
	} else {
		fmt.Println("Primeiro pagamento via Stripe foi sucesso")
	}
	errS2 := FinalizarCompra(stripe, 20)
	if errS2 != nil {
		fmt.Println(errS2)
	} else {
		fmt.Println("Segundo pagamento via Stripe foi sucesso")
	}

	// ----------------------------------------------------------------

	notificadores := []interfaces.INotificador{
		models.EmailNotificador{},
		models.SMSNotificador{},
		models.SMSNotificador{},
		models.SMSNotificador{},
		models.EmailNotificador{},
		models.EmailNotificador{},
	}

	for _, n := range notificadores {
		err := n.Notificar("teste@teste.c", "Teste de notificação")
		if err != nil {
			fmt.Println("Erro ao notificar:", err)
		}
	}

}

func FinalizarCompra(p interfaces.IProcessadorPagamento, valor float64) error {
	return p.Processar(valor)
}
