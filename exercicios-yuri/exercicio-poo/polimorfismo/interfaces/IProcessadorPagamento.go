package interfaces

type IProcessadorPagamento interface {
	Processar(valor float64) error
}
