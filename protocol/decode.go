package protocol

import (
	"fmt"
	"log"
	"strings"

	"github.com/yusank/godis/api"
)

/*
 * put all  decode functions
 */

// Receive 代表客户端请求来的数据
// A client sends the Redis server a RESP Array consisting of just Bulk Strings.
// 所以不考虑太多特殊情况
type Receive []string

func (r Receive) String() string {
	sb := new(strings.Builder)

	_, _ = sb.WriteString("[ ")
	for i, e := range r {
		_, _ = sb.WriteString(fmt.Sprintf("%v", e))
		if i < len(r)-1 {
			_, _ = sb.WriteString(", ")
		}
	}
	_, _ = sb.WriteString(" ]")

	return sb.String()
}

func DecodeFromReader(r api.Reader) (rec Receive, err error) {
	rec = make([]string, 0)
	b, err := r.ReadBytes('\n')
	if err != nil {
		log.Println("readBytes err:", err)
		return nil, err
	}

	str, length, desc, err := decodeSingleLine(b)
	if err != nil {
		log.Println("init message err:", err)
		return nil, err
	}

	if desc == DescriptionBulkStrings {
		temp, err1 := readBulkStrings(r, length)
		if err1 != nil {
			log.Println("read bulk str err:", err1)
			return nil, err1
		}

		rec = append(rec, string(temp))
		return
	}

	if desc == descriptionArrays {
		// won't sava array element
		items, err1 := readArray(r, length)
		if err1 != nil {
			log.Println("read bulk str err:", err1)
			return nil, err1
		}

		rec = append(rec, items...)
		return
	}

	rec = append(rec, str)
	return

}

func decodeSingleLine(line []byte) (str string, length int, desc byte, err error) {
	if len(line) < 3 {
		return "", 0, 0, fmt.Errorf("unsupported protocol")
	}

	desc = line[0]
	switch desc {
	// bulk string
	case DescriptionBulkStrings, descriptionArrays:
		length = readBulkOrArrayLength(line)
	case DescriptionSimpleStrings, DescriptionErrors, DescriptionIntegers:
		str = string(line[1 : len(line)-CRLFLen])
	default:
		return "", 0, 0, fmt.Errorf("unsupport protocol: %s", string(desc))
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

func readBulkStrings(r api.Reader, ln int) (val []byte, err error) {
	if ln < 0 {
		return
	}

	val = make([]byte, ln+2)
	_, err = r.Read(val)

	// trim last \r\n
	val = val[:ln]
	return
}

func readArray(r api.Reader, ln int) ([]string, error) {
	var items []string
	for i := 0; i < ln; i++ {
		line, err := r.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		if line[0] == descriptionArrays {
			sub, subErr := readArray(r, readBulkOrArrayLength(line))
			if subErr != nil {
				return nil, subErr
			}

			items = append(items, sub...)
		}

		str, length, desc, err := decodeSingleLine(line)
		if err != nil {
			return nil, err
		}

		if desc == DescriptionBulkStrings {
			temp, err1 := readBulkStrings(r, length)
			if err1 != nil {
				log.Println("read bulk str err:", err1)
				return nil, err1
			}

			str = string(temp)
		}

		items = append(items, str)
	}

	return items, nil
}
