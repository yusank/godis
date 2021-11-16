package e2e

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/yusank/godis/protocal"
)

func Test_SimpleString(t *testing.T) {
	msgChan, cancel := prepare(":7379", t)
	defer cancel()

	msg := protocal.NewMessage(protocal.SimpleString("OK"))
	msgChan <- msg
	close(msgChan)
	time.Sleep(time.Second)
}

func Test_BulkString(t *testing.T) {
	msgChan, cancel := prepare(":7379", t)
	defer cancel()

	msg := protocal.NewMessage(protocal.BulkString("Hello"))
	msgChan <- msg
	close(msgChan)
	time.Sleep(time.Second)
}

func Test_Array(t *testing.T) {
	msgChan, cancel := prepare(":7379", t)
	defer cancel()

	msg := protocal.NewMessage(protocal.Array(protocal.BulkString("hello"), protocal.SimpleString("world")))
	msgChan <- msg
	close(msgChan)
	time.Sleep(time.Second)
}

func prepare(addr string, t *testing.T) (msgChan chan *protocal.Message, cancel context.CancelFunc) {
	log.SetFlags(log.Lshortfile)
	ctx, cancel := context.WithCancel(context.Background())
	go startServer(addr, ctx)
	time.Sleep(time.Second)

	msgChan = make(chan *protocal.Message, 1)

	go func() {
		_ = connAndSendMsg(addr, msgChan, t)
	}()

	return msgChan, cancel

}
