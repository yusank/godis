package protocol

import (
	"fmt"
	"log"

	"github.com/yusank/godis/conn"
)

func initElementFromLine(line []byte) (e *Element, err error) {
	if len(line) == 0 {
		return nil, nil
	}

	if len(line) < 3 {
		return nil, fmt.Errorf("unsupported protocal")
	}

	desc := line[0]
	switch desc {
	// bulk string
	case DescriptionBulkStrings:
		e = NewBulkStringElement("")
		e.Len = readBulkOrArrayLength(line)
		if e.Len < 0 {
			e = NewNilBulkStringElement()
		}
	case DescriptionSimpleStrings:
		e = NewSimpleStringElement(string(line[1 : len(line)-2]))
	case DescriptionErrors:
		e = NewErrorElement(string(line[1 : len(line)-2]))
	case DescriptionIntegers:
		e = NewIntegerElement(string(line[1 : len(line)-2]))
	case DescriptionArray:
		e = NewArrayElement(readBulkOrArrayLength(line))
	default:
		return nil, fmt.Errorf("unsupport protocal: %s", string(desc))
	}

	return
}

func readBulkOrArrayLength(line []byte) int {
	if line[0] == '-' {
		return -1
	}

	var (
		ln int
	)
	for i := 1; line[i] != '\r'; i++ {
		ln = (ln * 10) + int(line[i]-'0')
	}

	return ln
}

func readBulkStrings(r conn.Reader, ln int) (val []byte, err error) {
	if ln <= 0 {
		return
	}

	val = make([]byte, ln+2)
	_, err = r.Read(val)

	// trim last \r\n
	val = val[:ln]
	return
}

func readArray(r conn.Reader, ln int) ([]*Element, error) {
	elements := []*Element{NewArrayElement(ln)}
	for i := 0; i < ln; i++ {
		line, err := r.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		if line[0] == DescriptionArray {
			sub, subErr := readArray(r, readBulkOrArrayLength(line))
			if subErr != nil {
				return nil, subErr
			}

			elements = append(elements, sub...)
		}

		e, err := initElementFromLine(line)
		if err != nil {
			return nil, err
		}

		if e.Description == DescriptionBulkStrings {
			temp, err1 := readBulkStrings(r, e.Len)
			if err1 != nil {
				log.Println("read bulk str err:", err1)
				return nil, err1
			}

			e.Value = string(temp)
		}

		elements = append(elements, e)
	}

	return elements, nil
}
