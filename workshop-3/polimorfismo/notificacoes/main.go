package main

import "fmt"

type Notificador interface {
	Notificar(destinatario, menssagem string)
}

type EmailNotificador struct{}

func (e *EmailNotificador) Notificar(destinatario, menssagem string) {
	fmt.Printf("Messagem: %s, enviada para %s, via Email.\n", menssagem, destinatario)
}

func EmailNotificadorFactory() *EmailNotificador {
	return &EmailNotificador{}
}

type SMSNotificador struct{}

func (s *SMSNotificador) Notificar(destinatario, menssagem string) {
	fmt.Printf("Messagem: %s, enviada para %s, via SMS.\n", destinatario, menssagem)
}

func SMSNotificarFactory() *SMSNotificador {
	return &SMSNotificador{}
}

func main() {

	notificadores := []Notificador{EmailNotificadorFactory(), SMSNotificarFactory()}
	for _, notificador := range notificadores {
		notificador.Notificar("Rafael Medeiros", "Menssagem teste")
	}
}
