package conn

type Reader interface {
	ReadBytes(byte) ([]byte, error)
	Read([]byte) (int, error)
}
