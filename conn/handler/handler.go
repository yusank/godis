package handler

type Handler interface {
	Handle(b []byte) error
}
