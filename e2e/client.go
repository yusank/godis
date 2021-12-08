package e2e

import (
	"net"

	"github.com/yusank/godis/protocol"
)

func connAndSendMsg(addr string, rsp *protocol.Response) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	n, err1 := conn.Write(rsp.Encode())
	if err1 != nil {
		return err1
	}
	// log.Println("write from client:", string(msg.Bytes()), n)

	buf := make([]byte, n)
	_, err1 = conn.Read(buf)
	if err1 != nil {
		return err1
	}
	// log.Println("read from conn: ", string(buf), n)
	return nil
}
