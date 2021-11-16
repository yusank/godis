package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yusank/godis/conn/handler"
	"github.com/yusank/godis/conn/tcp"
)

type Server struct {
	addr    string
	ctx     context.Context
	handler handler.Handler
}

func NewServer(addr string, ctx context.Context, h handler.Handler) *Server {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Server{
		addr:    addr,
		ctx:     ctx,
		handler: h,
	}
}

func (s *Server) Start() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		if err := tcp.Listen(s.ctx, s.addr, s.handler); err != nil {
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
