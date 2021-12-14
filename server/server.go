package server

type IServer interface {
	Start(addr string) error
	Stop()
}
