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

func newMessageFromOption(opts ...option) *Message {
	msg := &Message{
		originalData: new(bytes.Buffer),
		Elements:     make([]*Element, 0),
	}

	for _, o := range opts {
		o(msg)
	}

	return msg
}

// NewMessageFromBulkStrings build message form native string slice
// For example strSlice = ["GET", "Key"]  ==>
//*3\r\n
// $3\r\n
// GET\r\n
// $3\r\n
// key\r\n
// new line for readable
func NewMessageFromBulkStrings(strSlice ...string) *Message {
	msg := &Message{
		originalData: new(bytes.Buffer),
	}

	_, _ = msg.originalData.WriteString(encodeBulkStrings(strSlice...))
	return msg
}

func NewMessageFromResults(results []interface{}) *Message {
	msg := &Message{
		originalData: new(bytes.Buffer),
	}

	if len(results) == 0 {
		return nil
	}

	_, _ = msg.originalData.WriteString(fmt.Sprintf("*%d\r\n", len(results)))
	for _, result := range results {
		if result == nil {
			_, _ = msg.originalData.WriteString(encodeNilString())
			continue
		}

		switch v := result.(type) {
		case uint, uint64, int, int64, float64:
			_, _ = msg.originalData.WriteString(encodeInteger(fmt.Sprintf("%v", v)))
		case string:
			_, _ = msg.originalData.WriteString(encodeBulkString(v))
		case error:
			_, _ = msg.originalData.WriteString(encodeError(v))
		}
	}

	return msg
}

func NewMessageFromSimpleStrings(str string) *Message {
	msg := &Message{
		originalData: new(bytes.Buffer),
	}

	_, _ = msg.originalData.WriteString(encodeSimpleString(str))
	return msg
}

func NewMessageFromError(err error) *Message {
	msg := &Message{
		originalData: new(bytes.Buffer),
	}

	_, _ = msg.originalData.WriteString(encodeError(err))
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

	msg = newMessageFromOption(withElements(e))

	if e.Description == descriptionArray {
		// won't sava array element
		elements, err1 := readArray(r, e.Len)
		if err1 != nil {
			log.Println("read bulk str err:", err1)
			return nil, err1
		}

		msg = newMessageFromOption(withElements(elements...))
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
