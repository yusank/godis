package conn

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"sync"
)

func Listen(ctx context.Context, address string, h Handler) (l net.Listener, err error) {
	l, err = net.Listen("tcp", address)
	if err != nil {
		return
	}

	log.Println("start listen ", address)
	go handleListener(ctx, l, h)
	return
}

func handleListener(ctx context.Context, l net.Listener, h Handler) {
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
	addConn(conn)
	defer destroyConn(conn)

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

var defConnPool = newConnPool()

type connPool struct {
	pool sync.Map
}

func newConnPool() *connPool {
	return &connPool{pool: sync.Map{}}
}

func addConn(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	v, load := defConnPool.pool.LoadOrStore(addr, conn)
	if load {
		old := v.(net.Conn)
		_ = old.Close()
		defConnPool.pool.Store(addr, conn)
	}

	return
}

func destroyConn(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	v, load := defConnPool.pool.LoadAndDelete(addr)
	if load {
		old := v.(net.Conn)
		_ = old.Close()
	}

	_ = conn.Close()
}

func DestroyAllConn() {
	defConnPool.pool.Range(func(key, value interface{}) bool {
		conn := value.(net.Conn)
		_ = conn.Close()
		return true
	})
}
