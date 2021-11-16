package protocal

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
)

type ElementType uint8

type Element struct {
	Type        ElementType // string int nil array
	Description byte        // + - $ : *
	Value       string
	Len         int // valid for bulk string and array
}

func NewSimpleStringElement(str string) *Element {
	return &Element{
		Type:        ElementTypeString,
		Description: DescriptionSimpleStrings,
		Value:       str,
	}
}

func NewErrorElement(e string) *Element {
	return &Element{
		Type:        ElementTypeString,
		Description: DescriptionErrors,
		Value:       e,
	}
}

func NewBulkStringElement(str string) *Element {
	return &Element{
		Type:        ElementTypeString,
		Description: DescriptionBulkStrings,
		Value:       str,
		Len:         len(str),
	}
}

func NewArrayElement(ln int) *Element {
	return &Element{
		Type:        ElementTypeArray,
		Description: DescriptionArray,
		Value:       strconv.Itoa(ln),
		Len:         ln,
	}
}

func NewIntegerElement(i int) *Element {
	return &Element{
		Type:        ElementTypeInt,
		Description: DescriptionIntegers,
		Value:       strconv.Itoa(i),
	}
}

func (e *Element) String() string {
	if e.Type == ElementTypeArray {
		return fmt.Sprintf("type:%s, descType:%s, value:%d\n", ElementTypeMap[e.Type], string(e.Description), e.Len)
	}
	return fmt.Sprintf("type:%s, descType:%s, value:%s\n", ElementTypeMap[e.Type], string(e.Description), e.Value)
}

type Message struct {
	OriginalData []byte // unserialization data
	Elements     []*Element
}

func NewMessageFromBytes(data []byte) (msg *Message, err error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	msg = &Message{
		OriginalData: data,
		Elements:     make([]*Element, 0),
	}

	if err = msg.Decode(); err != nil {
		return nil, fmt.Errorf("decode data fail. err: %w", err)
	}

	if err = validArray(msg.Elements); err != nil {
		return nil, err
	}

	return msg, nil
}

func NewMessage(opts ...Option) (msg *Message, err error) {
	msg = &Message{
		Elements: make([]*Element, 0),
	}

	for _, o := range opts {
		o(msg)
	}

	return msg, nil
}

func (m *Message) Decode() error {
	var (
		startAt int
		err     error
	)

	for startAt < len(m.OriginalData) && m.OriginalData[startAt] != 0 {
		startAt, err = m.decodeSingle(startAt)
		if err != nil {
			return err
		}
	}

	m.OriginalData = m.OriginalData[:startAt]
	return err
}

func (m *Message) decodeSingle(startAt int) (n int, err error) {
	if len(m.OriginalData) == 0 || startAt >= len(m.OriginalData) {
		return 0, nil
	}

	var f readFunc
	desc := m.OriginalData[startAt]
	log.Println("desc: ", desc, string(desc))
	switch desc {
	// bulk string
	case DescriptionBulkStrings:
		f = readBulkString
	case DescriptionSimpleStrings:
		f = readSimpleString
	case DescriptionErrors:
		f = readError
	case DescriptionIntegers:
		f = readInteger
	case DescriptionArray:
		f = readArray
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
		encodeByte.WriteByte(e.Description)
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

/*
*3\r\n
$3\r\n
foo\r\n
$-1\r\n
$3\r\n
bar\r\n

["foo", nil, "bar"]
*/
