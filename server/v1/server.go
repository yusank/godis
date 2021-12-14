package v1

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/yusank/godis/event"
	"github.com/yusank/godis/protocol"
)

// Server provide tcp server based on go native net package
// Deprecated
type Server struct {
	addr     string
	ctx      context.Context
	listener net.Listener
	wg       *sync.WaitGroup
}

// NewServer return new Server
// Deprecated
func NewServer() *Server {
	return &Server{
		wg: new(sync.WaitGroup),
	}

}

// Start a new Server
// Deprecated
func (s *Server) Start(ctx context.Context, addr string) error {
	// addr
	s.addr = strings.TrimPrefix(addr, "tcp://")

	// ctx
	s.ctx = ctx
	if ctx == nil {
		s.ctx = context.Background()
	}

	// listen
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Println("listen err:", err)
		return err
	}
	log.Println("listen: ", l.Addr())
	s.listener = l

	// wait
	s.handleListener()
	return nil
}

// Stop the Server
// Deprecated
func (s *Server) Stop(_ context.Context) {
	if s.listener != nil {
		_ = s.listener.Close()
	}
	event.Stop()
}

// Deprecated
func (s *Server) handleListener() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("closed")
				break
			}

			log.Println("accept err:", err)
			continue
		}

		//log.Println("new conn from:", conn.RemoteAddr().String())
		s.wg.Add(1)
		go s.handleConn(conn)
	}

	s.wg.Wait()
}

// handle by a new goroutine
// Deprecated
func (s *Server) handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	ar := protocol.ReceiveDataAsync(reader)
loop:
	for {
		// ctx
		select {
		case <-s.ctx.Done():
			break loop
		case <-ar.ErrorChan:
			//log.Println("handle err:", err)
			break loop
		case rec := <-ar.ReceiveChan:
			reply := event.HandleRequest(rec)
			if len(reply) == 0 {
				continue
			}

			_, err := conn.Write(reply)
			if err != nil {
				log.Println("write err:", err)
				break loop
			}
		}
	}

	_ = conn.Close()
	s.wg.Done()
}
