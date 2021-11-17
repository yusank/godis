package protocol

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/yusank/godis/conn"
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

func NewNilBulkStringElement() *Element {
	return &Element{
		Type:        ElementTypeNil,
		Description: DescriptionBulkStrings,
		Value:       "-1",
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

func NewIntegerElement(is string) *Element {
	return &Element{
		Type:        ElementTypeInt,
		Description: DescriptionIntegers,
		Value:       is,
	}
}

func (e *Element) String() string {
	if e.Type == ElementTypeArray {
		return fmt.Sprintf("type:%s, descType:%s, value:%d\n", ElementTypeMap[e.Type], string(e.Description), e.Len)
	}
	return fmt.Sprintf("type:%s, descType:%s, value:%s\n", ElementTypeMap[e.Type], string(e.Description), e.Value)
}

type Message struct {
	originalData []byte // unserialization data
	Elements     []*Element
}

func NewMessage(opts ...Option) *Message {
	msg := &Message{
		Elements: make([]*Element, 0),
	}

	for _, o := range opts {
		o(msg)
	}

	return msg
}

func NewMessageFromReader(r conn.Reader) (msg *Message, err error) {
	b, err := r.ReadBytes('\n')
	if err != nil {
		log.Println("readBytes err:", err)
		return nil, err
	}

	e, err := initElementFromLine(b)
	if err != nil {
		log.Println("init message err:", err)
		return nil, err
	}

	if e.Description == DescriptionBulkStrings {
		temp, err1 := readBulkStrings(r, e.Len)
		if err1 != nil {
			log.Println("read bulk str err:", err1)
			return nil, err1
		}

		log.Println(string(temp))
		e.Value = string(temp)
	}

	msg = NewMessage(WithElements(e))

	if e.Description == DescriptionArray {
		elements, err1 := readArray(r, e.Len)
		if err1 != nil {
			log.Println("read bulk str err:", err1)
			return nil, err1
		}

		msg = NewMessage(WithElements(elements...))
	}

	return
}

func (m *Message) Bytes() []byte {
	return m.originalData
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

	m.originalData = make([]byte, encodeByte.Len())
	copy(m.originalData, encodeByte.Bytes())

	return nil
}

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
