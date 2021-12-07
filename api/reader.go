package api

// Reader define tcp message read api
type Reader interface {
	ReadBytes(byte) ([]byte, error)
	Read([]byte) (int, error)
}
