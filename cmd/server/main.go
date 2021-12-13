package main

import (
	"context"
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/yusank/godis/debug"
	"github.com/yusank/godis/redis"
	"github.com/yusank/godis/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var addr string
	flag.StringVar(&addr, "addr", ":7379", "server address")
	flag.Parse()

	if debug.DEBUG {
		insertPreData()
	}

	go func() {
		if !debug.DEBUG {
			return
		}

		for {
			time.Sleep(10 * time.Millisecond)
			if n := runtime.NumGoroutine(); n > 100 {
				log.Println("Goroutine: ", n)
			}
		}
	}()

	srv := server.NewServer(addr, nil)

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

func insertPreData() {
	redis.PrintSupportedCmd()
	for _, datum := range prepareData {
		c := &redis.Command{
			Command: datum[0],
			Values:  datum[1:],
		}

		c.ExecuteWithContext(context.TODO())
	}

	log.Println("data prepared")
}
