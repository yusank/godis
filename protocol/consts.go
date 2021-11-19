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
