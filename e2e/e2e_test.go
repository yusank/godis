package e2e

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/yusank/godis/protocol"
)

func Test_SimpleString(t *testing.T) {
	msgChan, cancel := prepare(":7379", t)
	defer cancel()

	msg := protocol.NewMessage(protocol.SimpleString("OK"))
	msgChan <- msg
	close(msgChan)
	time.Sleep(time.Second)
}

func Test_BulkString(t *testing.T) {
	msgChan, cancel := prepare(":7379", t)
	defer cancel()

	msg := protocol.NewMessage(protocol.BulkString("Hello"))
	msgChan <- msg
	close(msgChan)
	time.Sleep(time.Second)
}

func Test_Array(t *testing.T) {
	msgChan, cancel := prepare(":7379", t)
	defer cancel()

	msg := protocol.NewMessage(protocol.Array(protocol.BulkString("hello"), protocol.SimpleString("world")))
	msgChan <- msg
	close(msgChan)
	time.Sleep(time.Second)
}

func prepare(addr string, t *testing.T) (msgChan chan *protocol.Message, cancel context.CancelFunc) {
	log.SetFlags(log.Lshortfile)
	ctx, cancel := context.WithCancel(context.Background())
	go startServer(addr, ctx)
	time.Sleep(time.Second)

	msgChan = make(chan *protocol.Message, 1)

	go func() {
		_ = connAndSendMsg(addr, msgChan, t)
	}()

	return msgChan, cancel

}
