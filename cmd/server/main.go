package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/panjf2000/gnet/pkg/pool/goroutine"

	"github.com/yusank/godis/debug"
	"github.com/yusank/godis/redis"
	"github.com/yusank/godis/server"
	v2 "github.com/yusank/godis/server/v2"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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

	// srv
	srv, err := getSrv(server.Version)
	if err != nil {
		panic(err)
	}

	var (
		wg      = new(sync.WaitGroup)
		sigChan = make(chan os.Signal, 1)
		errChan = make(chan error, 1)
	)

	// signal
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = srv.Start(addr); err != nil {
			log.Println("service got err: ", err)
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-sigChan:
			srv.Stop()
		case <-errChan:
			srv.Stop()
		}
	}()

	// wait for service shutdown graceful
	wg.Wait()
}

var deferred = make([]func(), 0)

func getSrv(version string) (server.IServer, error) {
	var srv server.IServer
	switch version {
	case "v1":
		//srv = v1.NewServer()
		return nil, errors.New("v1 already deprecated")
	case "v2":
		p := goroutine.Default()
		deferred = append(deferred, p.Release)
		srv = v2.NewServer(p)
	default:
		return nil, errors.New("unknown server version")
	}

	return srv, nil
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

		c.ExecuteAsync()
	}

	log.Println("data prepared")
}
