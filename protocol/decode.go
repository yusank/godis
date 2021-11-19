package protocol

import (
	"fmt"
	"log"
	"strconv"

	"github.com/yusank/godis/conn"
)

/*
 * put all messages and elements decode funtions
 */

func newSimpleStringElement(str string) *Element {
	return &Element{
		Description: DescriptionSimpleStrings,
		Value:       str,
	}
}

func newErrorElement(e string) *Element {
	return &Element{
		Description: DescriptionErrors,
		Value:       e,
	}
}

func newBulkStringElement(str string) *Element {
	return &Element{
		Description: DescriptionBulkStrings,
		Value:       str,
	}
}

func newNilBulkStringElement() *Element {
	return &Element{
		Description: DescriptionBulkStrings,
		Value:       "-1",
	}
}

func newIntegerElement(is string) *Element {
	return &Element{
		Description: DescriptionIntegers,
		Value:       is,
	}
}

// only use when decode protocal data to msg won't store in elements slice
func newArrayElement(ln int) *Element {
	return &Element{
		Description: descriptionArray,
		Value:       strconv.Itoa(ln),
		Len:         ln,
	}
}

func initElementFromLine(line []byte) (e *Element, err error) {
	if len(line) == 0 {
		return nil, nil
	}

	if len(line) < 3 {
		return nil, fmt.Errorf("unsupported protocol")
	}

	desc := line[0]
	switch desc {
	// bulk string
	case DescriptionBulkStrings:
		e = newBulkStringElement("")
		e.Len = readBulkOrArrayLength(line)
		if e.Len < 0 {
			e = newNilBulkStringElement()
		}
	case DescriptionSimpleStrings:
		e = newSimpleStringElement(string(line[1 : len(line)-CRLFLen]))
	case DescriptionErrors:
		e = newErrorElement(string(line[1 : len(line)-CRLFLen]))
	case DescriptionIntegers:
		e = newIntegerElement(string(line[1 : len(line)-CRLFLen]))
	case descriptionArray:
		e = newArrayElement(readBulkOrArrayLength(line))
	default:
		return nil, fmt.Errorf("unsupport protocol: %s", string(desc))
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
	elements := []*Element{}
	for i := 0; i < ln; i++ {
		line, err := r.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		if line[0] == descriptionArray {
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
