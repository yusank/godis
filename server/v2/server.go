package v2

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"

	"github.com/yusank/godis/event"
	"github.com/yusank/godis/protocol"
)

type Server struct {
	*gnet.EventServer
	pool *goroutine.Pool
	addr string
}

func NewServer(p *goroutine.Pool) *Server {
	return &Server{
		EventServer: nil,
		pool:        p,
	}
}

func (s *Server) Start(addr string) error {
	s.addr = addr
	log.Println("listen: ", addr)
	return gnet.Serve(s, addr, gnet.WithMulticore(true), gnet.WithReusePort(true))
}

func (s *Server) Stop() {
	log.Println("graceful shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := gnet.Stop(ctx, s.addr); err != nil {
		log.Println("stop with err:", err)
	}
}

func (s *Server) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	var data = make([]byte, len(frame))
	copy(data, frame)

	_ = s.pool.Submit(func() {
		buf := bytes.NewBuffer(data)
		rec, err := protocol.DecodeFromReader(buf)
		if err != nil {
			log.Println(err)
			return
		}

		reply := event.HandleRequest(rec)
		if len(reply) == 0 {
			return
		}

		_ = c.AsyncWrite(reply)
	})

	return
}

func (s *Server) OnShutdown(srv gnet.Server) {
	event.Stop()
}
