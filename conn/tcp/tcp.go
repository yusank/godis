package tcp

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/yusank/godis/conn/handler"
)

func Listen(ctx context.Context, address string, h handler.Handler) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer l.Close()

	fmt.Println("start listen ", address)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			conn, err := l.Accept()
			if err != nil {
				return err
			}

			fmt.Println("new conn from:", conn.RemoteAddr().String())

			go handle(ctx, conn, h)
		}
	}
}

// handle by a new goroutine
func handle(ctx context.Context, conn net.Conn, h handler.Handler) {
	defer conn.Close()
	//reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			//data, err := io.ReadAll(conn)
			//data, _, err := reader.ReadLine()
			if err != nil {
				log.Println("readAll err:", err, " len:  ", n)
				return
			}

			reply, err := h.Handle(buf)
			if err != nil {
				log.Println("handle err:", err)
				return
			}

			_, err = conn.Write(reply)
			if err != nil {
				log.Println("write err:", err)
				return
			}
		}
	}
}
