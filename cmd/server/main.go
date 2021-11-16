package main

import (
	"flag"
	"log"

	"github.com/yusank/godis/conn/handler"
	"github.com/yusank/godis/server"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":7379", "server address")
	flag.Parse()

	srv := server.NewServer(addr, nil, &handler.PrintHandler{})

	if err := srv.Start(); err != nil {
		log.Println("exiting: ", err)
	}
}
