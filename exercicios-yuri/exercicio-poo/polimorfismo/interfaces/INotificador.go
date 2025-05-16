package interfaces

type INotificador interface {
	Notificar(destinatario, mensagem string) error
}
