package e2e

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/yusank/godis/protocal"
)

func Test_basic_connection(t *testing.T) {
	msgChan, cancel := prepare(":7379", t)
	defer cancel()

	msg := protocal.NewMessage(protocal.SimpleString("OK"))
	msgChan <- msg
	close(msgChan)
	time.Sleep(time.Second)
}

func Test_bulk_msg(t *testing.T) {
	msgChan, cancel := prepare(":7379", t)
	defer cancel()

	msg := protocal.NewMessage(protocal.BulkString("Hello"))
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
