package conn

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
)

func Listen(ctx context.Context, address string, h Handler) (l net.Listener, err error) {
	l, err = net.Listen("tcp", address)
	if err != nil {
		return
	}

	log.Println("start listen ", address)
	go handleListner(ctx, l, h)
	return
}

func handleListner(ctx context.Context, l net.Listener, h Handler) {
	defer l.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := l.Accept()
			if err != nil {
				log.Println("accept err:", err)
				return
			}

			log.Println("new conn from:", conn.RemoteAddr().String())
			go handle(ctx, conn, h)
		}
	}
}

// handle by a new goroutine
func handle(ctx context.Context, conn net.Conn, h Handler) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			reply, err := h.Handle(reader)
			if err == io.EOF {
				return
			}

			if err != nil {
				log.Println("handle err:", err)
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
