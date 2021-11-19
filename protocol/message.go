package protocol

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/yusank/godis/conn"
)

type Element struct {
	Description byte // + - $ : *
	Value       string
	Len         int // valid for bulk string and array
}

func (e *Element) String() string {
	return fmt.Sprintf("descType:%s, value:%s\n", string(e.Description), e.Value)
}

type Message struct {
	originalData *bytes.Buffer // unserialization data
	Elements     []*Element    // only use for transfer to command
}

func NewMessage(opts ...Option) *Message {
	msg := &Message{
		originalData: new(bytes.Buffer),
		Elements:     make([]*Element, 0),
	}

	for _, o := range opts {
		o(msg)
	}

	return msg
}

func NewMessageFromEncodeData(encodeData ...string) *Message {
	msg := &Message{
		originalData: new(bytes.Buffer),
	}

	for _, str := range encodeData {
		_, _ = msg.originalData.WriteString(str)
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

	if e.Description == descriptionArray {
		// won't sava array element
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
	return m.originalData.Bytes()
}

func (m *Message) String() string {
	sb := new(strings.Builder)

	_, _ = sb.WriteString("[ ")
	for i, e := range m.Elements {
		_, _ = sb.WriteString(e.Value)
		if i < len(m.Elements)-1 {
			_, _ = sb.WriteString(", ")
		}
	}
	_, _ = sb.WriteString(" ]")

	return sb.String()
}
