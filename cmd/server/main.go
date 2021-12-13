package main

import (
	"context"
	"errors"
	"flag"
	"log"

	"github.com/panjf2000/gnet/pkg/pool/goroutine"

	"github.com/yusank/godis/debug"
	"github.com/yusank/godis/redis"
	"github.com/yusank/godis/server"
	v1 "github.com/yusank/godis/server/v1"
	v2 "github.com/yusank/godis/server/v2"
)

func main() {
	defer func() {
		for _, f := range deferred {
			f()
		}
	}()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var addr string
	flag.StringVar(&addr, "addr", "tcp://:7379", "server address")
	flag.Parse()

	log.Printf("addr:%s, debug:%t, version:%s\n", addr, debug.DEBUG, server.Version)

	if debug.DEBUG {
		insertPreData()
	}

	if err := run(addr, server.Version); err != nil {
		log.Println("exiting: ", err)
	}
}

var deferred = make([]func(), 0)

func run(addr string, version string) error {
	var srv server.IServer
	switch version {
	case "v1":
		srv = v1.NewServer()
	case "v2":
		p := goroutine.Default()
		deferred = append(deferred, p.Release)
		srv = v2.NewServer(p)
	default:
		return errors.New("unknown server version")
	}

	return srv.Start(context.Background(), addr)
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
