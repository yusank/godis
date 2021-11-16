package e2e

import (
	"log"
	"net"

	"github.com/yusank/godis/protocal"
)

func connAndSendMsg(addr string, msgChan chan *protocal.Message) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	for msg := range msgChan {
		if err = msg.Encode(); err != nil {
			return err
		}

		n, err1 := conn.Write(msg.OriginalData)
		if err1 != nil {
			return err1
		}
		log.Println("write len:", n)

		buf := make([]byte, n)
		n, err1 = conn.Read(buf)
		if err1 != nil {
			return err1
		}
		log.Println("read from conn: ", string(buf), n)
	}

	return nil
}
