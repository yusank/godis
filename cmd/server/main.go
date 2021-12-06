package main

import (
	"flag"
	"log"

	"github.com/yusank/godis/handler"
	"github.com/yusank/godis/redis"
	"github.com/yusank/godis/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		addr      string
		debugMode bool
	)
	flag.StringVar(&addr, "addr", ":7379", "server address")
	flag.BoolVar(&debugMode, "debug", false, "debug mode, prepare some data")
	flag.Parse()

	if debugMode {
		debug()
	}

	srv := server.NewServer(addr, nil, &handler.TCPHandler{})

	if err := srv.Start(); err != nil {
		log.Println("exiting: ", err)
	}
}

var prepareData = [][]string{
	{"rpush", "list1", "1", "2", "3", "4", "5", "6"},
	{"set", "key1", "hello"},
	{"set", "key2", "10"},
	{"zadd", "zset", "1", "a", "2", "b", "3", "c", "4", "d", "5", "e"},
}

func debug() {
	for _, datum := range prepareData {
		c := &redis.Command{
			Command: datum[0],
			Values:  datum[1:],
		}

		c.Execute(nil)
	}

	log.Println("data prepared")
}
