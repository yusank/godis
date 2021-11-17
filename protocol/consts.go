package protocol

const (
	DescriptionSimpleStrings = '+'
	DescriptionErrors        = '-'
	DescriptionIntegers      = ':'
	DescriptionBulkStrings   = '$'
	DescriptionArray         = '*'
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

const (
	ElementTypeString ElementType = iota + 1
	ElementTypeInt
	ElementTypeNil
	ElementTypeArray
)

var ElementTypeMap = map[ElementType]string{
	ElementTypeString: "string",
	ElementTypeInt:    "int",
	ElementTypeNil:    "nil",
	ElementTypeArray:  "array",
}
