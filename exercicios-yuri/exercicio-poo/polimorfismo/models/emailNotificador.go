package models

import (
	"errors"
	"fmt"
)

type EmailNotificador struct{}

func (e EmailNotificador) Notificar(destinatario, mensagem string) error {
	if destinatario == "" || mensagem == "" {
		return errors.New("EmailNotificador - Notificar: destinatário e mensagem estão vazios")
	}
	fmt.Printf("[EmailNotificador][%s] %s\n", destinatario, mensagem)
	return nil
}
