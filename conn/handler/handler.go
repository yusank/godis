package handler

type Handler interface {
	Handle(b []byte) ([]byte, error)
}
