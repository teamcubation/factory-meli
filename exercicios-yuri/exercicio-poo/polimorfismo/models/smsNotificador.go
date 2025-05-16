package models

import (
	"errors"
	"fmt"
)

type SMSNotificador struct{}

func (s SMSNotificador) Notificar(destinatario, mensagem string) error {
	if destinatario == "" || mensagem == "" {
		return errors.New("SMSNotificador - Notificar: destinatário e mensagem estão vazios")
	}
	fmt.Printf("[SMSNotificador][%s] %s\n", destinatario, mensagem)
	return nil
}
