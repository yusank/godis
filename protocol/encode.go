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

func NewResponse(isArray ...bool) *Response {
	rsp := &Response{
		Items: make([]*ResponseItem, 0),
	}

	if len(isArray) > 0 && isArray[0] {
		rsp.IsArray = true
	}

	return rsp
}

func NewResponseWithSimpleString(str string) *Response {
	resp := NewResponse()
	resp.Items = append(resp.Items, &ResponseItem{
		Value: str,
		Type:  TypeSimpleStrings,
	})

	return resp
}

func NewResponseWithBulkString(str string) *Response {
	resp := NewResponse()
	resp.Items = append(resp.Items, &ResponseItem{
		Value: str,
		Type:  TypeBulkStrings,
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

func NewResponseWithNilBulk() *Response {
	resp := NewResponse(false)
	return resp
}

func NewResponseWithEmptyArray() *Response {
	resp := NewResponse(true)
	return resp
}

func NewResponseWithInteger(i int64) *Response {
	resp := NewResponse()
	resp.Items = append(resp.Items, &ResponseItem{
		Value: i,
		Type:  TypeIntegers,
	})

	return resp
}

func (r *Response) SetIsArray() {
	r.IsArray = true
}

func (r *Response) AppendBulkInterfaces(slice ...interface{}) {
	for _, v := range slice {
		r.Items = append(r.Items, &ResponseItem{
			Value: v,
			Type:  TypeBulkStrings,
		})
	}
}

func (r *Response) AppendBulkStrings(strSlice ...string) {
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
		// empty
		if r.IsArray {
			return []byte("*0\r\n")
		}

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
