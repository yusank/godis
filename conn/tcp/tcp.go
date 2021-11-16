package tcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/yusank/godis/handler"
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
			reply, err := h.Handle(reader)
			if err != nil {
				log.Println("handle err:", err)
				if err == io.EOF {
					log.Println("eof")
				}
				return
			}

			if len(reply) == 0 {
				continue
			}

			_, err = conn.Write(reply)
			if err != nil {
				log.Println("write err:", err)
				return
			}
		}
	}
}
