package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/yusank/godis/protocal"
)

func Test_basic_connection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go startServer(ctx)
	time.Sleep(time.Second)

	msgChan := make(chan *protocal.Message, 1)

	go func() {
		_ = connAndSendMsg(":7379", msgChan)
	}()

	msg, err := protocal.NewMessage(protocal.SimpleString("OK"))
	if err != nil {
		t.Error(err)
		return
	}

	msgChan <- msg
	close(msgChan)
	time.Sleep(time.Second)
}
