package e2e

import (
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yusank/godis/protocal"
)

func connAndSendMsg(addr string, msgChan chan *protocal.Message, t *testing.T) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	for msg := range msgChan {
		if err = msg.Encode(); err != nil {
			return err
		}

		n, err1 := conn.Write(msg.Bytes())
		if err1 != nil {
			return err1
		}
		log.Println("write from client:", string(msg.Bytes()), n)

		buf := make([]byte, n)
		n, err1 = conn.Read(buf)
		if err1 != nil {
			return err1
		}
		log.Println("read from conn: ", string(buf), n)

		assert.Equal(t, msg.Bytes(), buf)
	}

	return nil
}
