package tcp

import (
	"bufio"
	"context"
	"fmt"
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
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, _, err := reader.ReadLine()
			if err != nil {
				fmt.Println("readAll err:", err)
				return
			}

			if err = h.Handle(data); err != nil {
				return
			}

			_, err = conn.Write(append(data, '\n'))
			if err != nil {
				fmt.Println("write err:", err)
				return
			}
		}
	}
}
