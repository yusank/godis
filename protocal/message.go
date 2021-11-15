package protocal

import (
	"bytes"
	"fmt"
	"strconv"
)

type ElementType uint8

type Element struct {
	Type            ElementType // string int nil array
	DescriptionType byte        // + - $ : *
	Value           string
	Len             int
}

func (e *Element) String() string {
	if e.Type == ElementTypeArray {
		return fmt.Sprintf("type:%s, descType:%s, value:%d\n", ElementTypeMap[e.Type], string(e.DescriptionType), e.Len)
	}
	return fmt.Sprintf("type:%s, descType:%s, value:%s\n", ElementTypeMap[e.Type], string(e.DescriptionType), e.Value)
}

type Message struct {
	OriginalData []byte // unserialization data
	Elements     []*Element
}

func NewMessage(data []byte) (msg *Message, err error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	msg = &Message{
		OriginalData: data,
	}

	if err = msg.Decode(); err != nil {
		return nil, fmt.Errorf("decode data fail. err: %w", err)
	}

	if err = validArray(msg.Elements); err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Message) Decode() error {
	var (
		startAt int
		err     error
	)

	for startAt < len(m.OriginalData) {
		startAt, err = m.decodeSingle(startAt)
		if err != nil {
			return err
		}
	}

	// todo check array lenth is valid
	return err
}

func (m *Message) decodeSingle(startAt int) (n int, err error) {
	if len(m.OriginalData) == 0 || startAt >= len(m.OriginalData) {
		return 0, nil
	}

	var f readFunc
	prefix := m.OriginalData[startAt]
	switch prefix {
	// bulk string
	case BulkStringsPrefix:
		f = readBulkString
	case SimpleStringsPrefix:
		f = readSimpleString
	case ErrorsPrefix:
		f = readError
	case IntegersPrefix:
		f = readInteger
	case ArrarysPrefix:
		f = readArrary
	default:
		return 0, fmt.Errorf("unsupport protocal")
	}

	element, st, err := f(startAt+1, m.OriginalData)
	if err != nil {
		return 0, err
	}

	m.Elements = append(m.Elements, element)

	return st, nil
}

func (m *Message) Encode() error {
	var encodeByte = new(bytes.Buffer)

	for _, e := range m.Elements {
		encodeByte.WriteByte(e.DescriptionType)
		if e.Len > 0 {
			encodeByte.WriteString(strconv.Itoa(e.Len))
			encodeByte.Write([]byte(CRLF))
		}

		if e.Type == ElementTypeArray {
			continue
		}

		encodeByte.WriteString(e.Value)
		encodeByte.Write([]byte(CRLF))
	}

	m.OriginalData = make([]byte, encodeByte.Len())
	copy(m.OriginalData, encodeByte.Bytes())

	return nil
}

// todo check array value
func validArray(elements []*Element) error {
	var checkFunc func(int, int) (int, error)

	checkFunc = func(startAt, require int) (int, error) {
		var (
			i        = startAt
			matchCnt int
		)

		for matchCnt < require && i < len(elements) {
			if i >= len(elements) {
				return -1, fmt.Errorf("invalid array")
			}

			if elements[i].Type == ElementTypeArray {
				offset, err := checkFunc(i+1, elements[i].Len)
				if err != nil {
					return -1, err
				}
				i = offset
			} else {
				i++
			}

			matchCnt++
		}

		if matchCnt != require {
			return -1, fmt.Errorf("array len and value not match")
		}

		return i, nil
	}

	if len(elements) < 2 {
		return nil
	}

	if elements[0].Type != ElementTypeArray {
		return nil
	}

	ln, err := checkFunc(1, elements[0].Len)
	if err != nil {
		return err
	}

	if ln != len(elements) {
		return fmt.Errorf("want:%d, got:%d", ln, len(elements))
	}

	return nil
}

// todo: add check orginal data is valid

/*
*3\r\n
$3\r\n
foo\r\n
$-1\r\n
$3\r\n
bar\r\n

["foo", nil, "bar"]
*/
