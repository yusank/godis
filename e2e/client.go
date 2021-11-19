package e2e

import (
	"net"

	"github.com/yusank/godis/protocol"
)

func connAndSendMsg(addr string, msg *protocol.Message) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	n, err1 := conn.Write(msg.Bytes())
	if err1 != nil {
		return err1
	}
	// log.Println("write from client:", string(msg.Bytes()), n)

	buf := make([]byte, n)
	n, err1 = conn.Read(buf)
	if err1 != nil {
		return err1
	}
	// log.Println("read from conn: ", string(buf), n)
	return nil
}
