package v1

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/yusank/godis/event"
	"github.com/yusank/godis/protocol"
)

type Server struct {
	addr     string
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener
	wg       *sync.WaitGroup
}

func NewServer() *Server {
	return &Server{
		wg: new(sync.WaitGroup),
	}

}

func (s *Server) Start(ctx context.Context, addr string) error {
	// init
	// ctx
	if ctx == nil {
		ctx = context.Background()
	}
	s.ctx, s.cancel = context.WithCancel(ctx)

	// addr
	s.addr = addr

	// signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// listen
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Println("listen err:", err)
		return err
	}
	log.Println("listen: ", l.Addr())
	s.listener = l

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		select {
		case <-s.ctx.Done():
			//s.Stop()
			log.Println("kill by ctx")
			return
		case sig := <-sigChan:
			s.Stop()
			log.Printf("kill by signal:%s", sig.String())
			return
		}
	}()

	s.handleListener()
	return nil
}

func (s *Server) Stop() {
	s.cancel()
	_ = s.listener.Close()
	event.Stop()
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
