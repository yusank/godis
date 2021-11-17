package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yusank/godis/conn"
)

type Server struct {
	addr    string
	ctx     context.Context
	cancel  context.CancelFunc
	handler conn.Handler
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
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		if err := conn.Listen(s.ctx, s.addr, s.handler); err != nil {
			fmt.Println(err)
			return
		}
	}()

	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	case sig := <-sigChan:
		return fmt.Errorf("kill by signal:%s", sig.String())
	}
}

func (s *Server) Stop() {
	s.cancel()
}
