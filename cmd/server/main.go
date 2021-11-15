package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yusank/godis/conn/handler"
	"github.com/yusank/godis/conn/tcp"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":7379", "server address")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		if err := tcp.Listen(ctx, addr, &handler.PrintHandler{}); err != nil {
			fmt.Println(err)
			return
		}
	}()

	select {
	case <-ctx.Done():
		return
	case <-sig:
		cancel()
	}

	fmt.Println("exit")
}
