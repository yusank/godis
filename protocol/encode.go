package protocol

import (
	"bytes"
	"strconv"

	"github.com/yusank/godis/util"
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
	resp.Items = []*ResponseItem{
		{
			Value: str,
			Type:  TypeSimpleStrings,
		},
	}

	return resp
}

func NewResponseWithBulkString(str string) *Response {
	resp := NewResponse()
	resp.Items = []*ResponseItem{
		{
			Value: str,
			Type:  TypeBulkStrings,
		},
	}

	return resp
}

func NewResponseWithError(e error) *Response {
	resp := NewResponse()
	resp.Items = []*ResponseItem{
		{
			Value: e,
			Type:  TypeErrors,
		},
	}

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
	resp.Items = []*ResponseItem{
		{
			Value: i,
			Type:  TypeIntegers,
		},
	}

	return resp
}

func (r *Response) SetIsArray() {
	r.IsArray = true
}

func (r *Response) AppendBulkInterfaces(slice ...interface{}) *Response {
	r.Items = make([]*ResponseItem, len(slice))
	for i, v := range slice {
		r.Items[i] = &ResponseItem{
			Value: v,
			Type:  TypeBulkStrings,
		}
	}

	return r
}

func (r *Response) AppendBulkStrings(strSlice ...string) *Response {
	r.Items = make([]*ResponseItem, len(strSlice))
	for i, str := range strSlice {
		r.Items[i] = &ResponseItem{
			Value: str,
			Type:  TypeBulkStrings,
		}
	}

	return r
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
	lnStr := strconv.Itoa(len(r.Items))
	buf.WriteString(util.StringConcat(len(lnStr)+3, "*", lnStr, CRLF))
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
		str := ri.Value.(string)
		_, _ = buf.WriteString(util.StringConcat(len(str)+3, "+", str, CRLF))
	case TypeErrors:
		str := ri.Value.(error).Error()
		_, _ = buf.WriteString(util.StringConcat(len(str)+3, "-", str, CRLF))
	case TypeIntegers:
		str := strconv.FormatInt(ri.Value.(int64), 10)
		_, _ = buf.WriteString(util.StringConcat(len(str)+3, ":", str, CRLF))
	case TypeBulkStrings:
		if ri.Value == nil {
			_, _ = buf.WriteString("$-1\r\n")
			break
		}
		str := ri.Value.(string)
		lnStr := strconv.Itoa(len(str))
		_, _ = buf.WriteString(util.StringConcat(len(str)+len(lnStr)+5, "$", lnStr, CRLF, str, CRLF))
	}
}
