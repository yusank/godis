//go:build v1
// +build v1

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
type Server struct {
	addr     string
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener
	wg       *sync.WaitGroup
}

// NewServer return new Server
func NewServer() *Server {
	return &Server{
		wg: new(sync.WaitGroup),
	}

}

// Start a new Server
func (s *Server) Start(addr string) error {
	// addr
	s.addr = strings.TrimPrefix(addr, "tcp://")

	// ctx
	s.ctx, s.cancel = context.WithCancel(context.Background())

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
func (s *Server) Stop() {
	// stop accept new connection
	if s.listener != nil {
		_ = s.listener.Close()
	}
	// stop handle new request from old connection
	event.Stop()
	// quit old connection loop
	s.cancel()
}

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
