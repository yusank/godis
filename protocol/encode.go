package protocol

import (
	"bytes"
	"fmt"
)

/*
 * put encode functions
 */

type Response struct {
	Items   []*ResponseItem
	IsArray bool // if there is only one item and want response as array then should set this field as true.
}

type ResponseItem struct {
	Value interface{}
	Type  ProtocolDataType
}

func NewResponse() *Response {
	return &Response{
		Items: make([]*ResponseItem, 0),
	}
}

func NewResponseWithSimpleString(str string) *Response {
	resp := NewResponse()
	resp.Items = append(resp.Items, &ResponseItem{
		Value: str,
		Type:  TypeSimpleStrings,
	})

	return resp
}

func NewResponseWithError(e error) *Response {
	resp := NewResponse()
	resp.Items = append(resp.Items, &ResponseItem{
		Value: e,
		Type:  TypeErrors,
	})

	return resp
}

func (r *Response) AppendBulkStrings(strSlice ...interface{}) {
	for _, str := range strSlice {
		r.Items = append(r.Items, &ResponseItem{
			Value: str,
			Type:  TypeBulkStrings,
		})
	}
}

func (r *Response) Encode() []byte {
	buf := new(bytes.Buffer)
	if len(r.Items) == 0 {
		return []byte("$-1\r\n")
	}

	if !r.IsArray {
		r.Items[0].encode(buf)
		return buf.Bytes()
	}

	// r.IsArray is true
	buf.WriteString(fmt.Sprintf("*%d\r\n", len(r.Items)))
	for _, item := range r.Items {
		item.encode(buf)
	}

	return buf.Bytes()
}

func (ri *ResponseItem) encode(buf *bytes.Buffer) {
	if buf == nil {
		return
	}

	switch ri.Type {
	case TypeSimpleStrings:
		_, _ = buf.WriteString(fmt.Sprintf("+%s\r\n", ri.Value.(string)))
	case TypeErrors:
		_, _ = buf.WriteString(fmt.Sprintf("-%s\r\n", ri.Value.(error).Error()))
	case TypeIntegers:
		_, _ = buf.WriteString(fmt.Sprintf(":%d\r\n", ri.Value.(int64)))
	case TypeBulkStrings:
		if ri.Value == nil {
			_, _ = buf.WriteString(fmt.Sprintf("$-1\r\n"))
			break
		}
		str := ri.Value.(string)
		_, _ = buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(str), str))
	}

	return
}
