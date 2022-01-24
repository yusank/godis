package protocol

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

/*
 * put all  decode functions
 */

// Receive 代表客户端请求来的数据
// A client sends the Redis server a RESP Array consisting of just Bulk Strings.
// 所以不考虑太多特殊情况
type Receive struct {
	Elements []string
	OrgStr   string
}

func (r Receive) String() string {
	sb := new(strings.Builder)

	_, _ = sb.WriteString("[ ")
	for i, e := range r.Elements {
		_, _ = sb.WriteString(fmt.Sprintf("%v", e))
		if i < len(r.Elements)-1 {
			_, _ = sb.WriteString(", ")
		}
	}
	_, _ = sb.WriteString(" ]")

	return sb.String()
}

func (r *Receive) append(ele, org string) {
	if r.Elements == nil {
		r.Elements = make([]string, 0)
	}

	r.Elements = append(r.Elements, ele)
	r.OrgStr += org
}

func (r *Receive) addOrg(s string) {
	r.OrgStr += s
}

type AsyncReceive struct {
	ReceiveChan chan *Receive
	ErrorChan   chan error
}

func ReceiveDataAsync(r Reader) *AsyncReceive {
	var ar = &AsyncReceive{
		ReceiveChan: make(chan *Receive, 1),
		ErrorChan:   make(chan error, 1),
	}
	go func() {
		defer func() {
			close(ar.ReceiveChan)
			close(ar.ErrorChan)
		}()

		for {
			rec, err := DecodeFromReader(r)
			if err != nil {
				ar.ErrorChan <- err

				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}
				log.Println(err)
				continue
			}

			ar.ReceiveChan <- rec
		}
	}()

	return ar
}

func DecodeFromReader(r Reader) (rec *Receive, err error) {
	rec = new(Receive)
	b, err := r.ReadBytes('\n')
	if err != nil {
		//log.Println("readBytes err:", err)
		return nil, err
	}

	length, desc, err := rec.decodeSingleLine(b)
	if err != nil {
		log.Println("init message err:", err)
		return nil, err
	}

	if desc == DescriptionBulkStrings {
		err1 := rec.readBulkStrings(r, length)
		if err1 != nil {
			log.Println("read bulk str err:", err1)
			return nil, err1
		}
	}

	if desc == descriptionArrays {
		// won't sava array element
		err1 := rec.readArray(r, length)
		if err1 != nil {
			log.Println("read bulk str err:", err1)
			return nil, err1
		}
	}

	return
}

func (r *Receive) decodeSingleLine(line []byte) (length int, desc byte, err error) {
	if len(line) < 3 {
		return 0, 0, fmt.Errorf("unsupported protocol")
	}

	r.addOrg(string(line))

	desc = line[0]
	switch desc {
	// bulk string
	case DescriptionBulkStrings, descriptionArrays:
		length = readBulkOrArrayLength(line)
	case DescriptionSimpleStrings, DescriptionErrors, DescriptionIntegers:
		r.append(string(line[1:len(line)-CRLFLen]), "")
		//str = string(line[1 : len(line)-CRLFLen])
	default:
		if string(line) == "PING\r\n" {
			r.append("PING", "PING\r\n")
			return 0, DescriptionSimpleStrings, nil
		}
		return 0, 0, fmt.Errorf("unsupport protocol: %s", string(line))
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

func (r *Receive) readBulkStrings(rr Reader, ln int) error {
	if ln < 0 {
		return fmt.Errorf("invalid length")
	}

	val := make([]byte, ln+2)
	_, err := rr.Read(val)
	if err != nil {
		return err
	}

	// trim last \r\n
	r.append(string(val[:ln]), string(val))

	return nil
}

func (r *Receive) readArray(rr Reader, ln int) error {
	for i := 0; i < ln; i++ {
		line, err := rr.ReadBytes('\n')
		if err != nil {
			return err
		}

		if line[0] == descriptionArrays {
			subErr := r.readArray(rr, readBulkOrArrayLength(line))
			if subErr != nil {
				return subErr
			}
		}

		length, desc, err := r.decodeSingleLine(line)
		if err != nil {
			return err
		}

		if desc == DescriptionBulkStrings {
			err = r.readBulkStrings(rr, length)
			if err != nil {
				log.Println("read bulk str err:", err)
				return err
			}
		}
	}

	return nil
}
