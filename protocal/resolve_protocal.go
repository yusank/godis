package protocal

import (
	"fmt"
	"strconv"
)

// readFunc is generic  func declaration for resolve diffrent protocal types
type readFunc func(startAt int, data []byte) (ele *Element, newStartAt int, err error)

func readBulkString(startAt int, data []byte) (ele *Element, newStartAt int, err error) {
	if data[startAt] == '-' {
		if string(data[startAt:startAt+2+CRLFLen]) != "-1\r\n" {
			return nil, startAt, fmt.Errorf("invalid protocal value")
		}

		return &Element{
			Type:            ElementTypeNil,
			DescriptionType: BulkStringsPrefix,
			Value:           "-1",
		}, startAt + 2 + CRLFLen, nil
	}

	var (
		ln int
	)
	for newStartAt = startAt; data[newStartAt] != '\r'; newStartAt++ {
		ln = (ln * 10) + int(data[newStartAt]-'0')
	}

	// $1\r\na\r\n
	newStartAt += CRLFLen
	ele = &Element{
		Type:            ElementTypeString,
		DescriptionType: BulkStringsPrefix,
		Len:             ln,
	}

	if ln > len(data[newStartAt:]) {
		return nil, startAt, fmt.Errorf("invalid protocal value")
	}

	if ln != 0 {
		ele.Value = string(data[newStartAt : newStartAt+ln])
	}

	newStartAt += CRLFLen
	return
}

func readSimpleString(startAt int, data []byte) (ele *Element, newStartAt int, err error) {
	ele = &Element{
		Type:            ElementTypeString,
		DescriptionType: SimpleStringsPrefix,
	}

	var b []byte
	for newStartAt = startAt; data[newStartAt] != '\r'; newStartAt++ {
		b = append(b, data[newStartAt])
	}

	ele.Value = string(b)
	newStartAt += CRLFLen
	return
}

func readError(startAt int, data []byte) (ele *Element, newStartAt int, err error) {
	ele = &Element{
		Type:            ElementTypeString,
		DescriptionType: ErrorsPrefix,
	}

	var b []byte
	for newStartAt = startAt; data[newStartAt] != '\r'; newStartAt++ {
		b = append(b, data[newStartAt])
	}

	ele.Value = string(b)
	newStartAt += CRLFLen
	return
}

func readInteger(startAt int, data []byte) (ele *Element, newStartAt int, err error) {
	ele = &Element{
		Type:            ElementTypeString,
		DescriptionType: IntegersPrefix,
	}

	var b []byte
	for newStartAt = startAt; data[newStartAt] != '\r'; newStartAt++ {
		b = append(b, data[newStartAt])
	}

	ele.Value = string(b)
	newStartAt += CRLFLen
	return
}

func readArrary(startAt int, data []byte) (ele *Element, newStartAt int, err error) {
	var (
		ln int
	)
	for newStartAt = startAt; data[newStartAt] != '\r'; newStartAt++ {
		ln = (ln * 10) + int(data[newStartAt]-'0')
	}

	ele = &Element{
		Type:            ElementTypeArray,
		DescriptionType: ArrarysPrefix,
		Len:             ln,
		Value:           strconv.Itoa(ln),
	}
	newStartAt += CRLFLen
	return
}