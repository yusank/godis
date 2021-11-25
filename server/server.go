package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/yusank/godis/conn"
)

type Server struct {
	addr     string
	ctx      context.Context
	cancel   context.CancelFunc
	handler  conn.Handler
	listener net.Listener
}

func NewServer(addr string, ctx context.Context, h conn.Handler) *Server {
	if ctx == nil {
		ctx = context.Background()
	}

	s := &Server{
		addr:    addr,
		handler: h,
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	return s
}

func (s *Server) Start() error {
	defer conn.DestroyAllConn()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// 半阻塞，开始监听端口后异步处理
	l, err := conn.Listen(s.ctx, s.addr, s.handler)
	if err != nil {
		fmt.Println(err)
		return err
	}

	s.listener = l

	select {
	case <-s.ctx.Done():
		return nil
	case sig := <-sigChan:
		s.Stop()
		fmt.Printf("kill by signal:%s", sig.String())
		return nil
	}
}

func (s *Server) Stop() {
	_ = s.listener.Close()
	s.cancel()
}
