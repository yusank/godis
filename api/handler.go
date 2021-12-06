package api

// Handler defines tcp connection handler interface
type Handler interface {
	// Handle return reply and error
	// TODO: should support timeout
	Handle(r Reader) ([]byte, error)
}
