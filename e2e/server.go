package e2e

import (
	"context"
	"log"

	"github.com/yusank/godis/conn/handler"
	"github.com/yusank/godis/server"
)

func startServer(ctx context.Context) {
	srv := server.NewServer(":7379", ctx, handler.PrintHandler{})

	if err := srv.Start(); err != nil {
		log.Println("exiting: ", err)
	}
}
