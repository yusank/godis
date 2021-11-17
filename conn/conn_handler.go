package conn

type Reader interface {
	ReadBytes(byte) ([]byte, error)
	Read([]byte) (int, error)
}

type Handler interface {
	// Handle return reply and error
	Handle(r Reader) ([]byte, error)
}
