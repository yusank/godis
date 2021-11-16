package protocal

const (
	DescriptionSimpleStrings = '+'
	DescriptionErrors        = '-'
	DescriptionIntegers      = ':'
	DescriptionBulkStrings  = '$'
	DescriptionArray          = '*'
)

const (
	CRLF      = "\r\n"
	CRLFDebug = "\\r\\n" // using when need to print \r\n

	CRLFLen = len(CRLF)
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
