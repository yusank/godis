package e2e

import (
	"context"
	"log"

	"github.com/yusank/godis/handler"
	"github.com/yusank/godis/server"
)

func startServer(addr string, ctx context.Context) {
	srv := server.NewServer(addr, ctx, handler.MessageHandler{})

	if err := srv.Start(); err != nil {
		log.Println("exiting: ", err)
	}
}
