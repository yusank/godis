package protocol

const (
	DescriptionSimpleStrings byte = '+'
	DescriptionErrors        byte = '-'
	DescriptionIntegers      byte = ':'
	DescriptionBulkStrings   byte = '$'
	descriptionArrays        byte = '*' // unexpose to external packages
)

const (
	CRLF      = "\r\n"
	CRLFDebug = "\\r\\n" // using when need to print \r\n

	CRLFLen = len(CRLF)
)

type ProtocolDataType int8

const (
	TypeSimpleStrings ProtocolDataType = iota + 1
	TypeErrors
	TypeIntegers
	TypeBulkStrings
	TypeArrays
)
