package protocol

const (
	DescriptionSimpleStrings byte = '+'
	DescriptionErrors        byte = '-'
	DescriptionIntegers      byte = ':'
	DescriptionBulkStrings   byte = '$'
	descriptionArray         byte = '*' // unexpose to external packages
)

const (
	CRLF      = "\r\n"
	CRLFDebug = "\\r\\n" // using when need to print \r\n

	CRLFLen = len(CRLF)
)

var (
	OK      = []byte("*1\r\n$2\r\nOK\r\n")
	Pong    = []byte("*1\r\n$4\r\nPONG\r\n")
	Command = []byte("*1\r\n$7\r\nCOMMAND\r\n")
)
