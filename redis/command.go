package redis

import (
	"context"
	"log"
	"strings"

	"github.com/yusank/godis/protocol"
)

type Command struct {
	Ctx     context.Context
	Command string
	Values  []string
}

// ExecuteFunc define command execute, returns slice of string as result and error if any error occur .
type ExecuteFunc func(*Command) (*protocol.Response, error)

var implementedCommands = map[string]ExecuteFunc{}

func NewCommandFromReceive(ctx context.Context, rec protocol.Receive) *Command {
	if len(rec) == 0 {
		return nil
	}

	c := &Command{
		Command: strings.ToLower(rec[0]),
		Values:  rec[1:],
	}

	return c
}

// TODO: 尝试用一个 channel 控制所有的命令 从而避免并发带来的一些问题 https://github.com/yusank/godis/issues/10
// 	借鉴 Redis 的时间处理器,并发接受 tcp 请求,并放入一个有序事件处理池子里, 而只有一个 worker 去处理这些事件
// 	如果性能差太多,那考虑从加一个 key 级别的轻量级锁(原子操作),一个 key 在任何时刻只有一个worker 去读写

// ExecuteAsync is async method which return a result channel and run command ina new go routine.
// if got any error when execution will transfer protocol bytes
func (c *Command) ExecuteAsync() chan *protocol.Response {
	var (
		rspChan = make(chan *protocol.Response)
	)

	go func() {
		defer close(rspChan)
		//defer func() {
		//	if recover() != nil {
		//		log.Println(c.Command, c.Values)
		//	}
		//}()
		f, ok := implementedCommands[c.Command]
		if !ok {
			c.setRsp(nil, ErrUnknownCommand, rspChan)
			return
		}

		result, err := f(c)
		c.setRsp(result, err, rspChan)
	}()

	return rspChan
}

func (c *Command) setRsp(result *protocol.Response, err error, ch chan *protocol.Response) {
	// if ctx has error won't put
	if c.Ctx != nil && c.Ctx.Err() != nil {
		return
	}

	if err != nil {
		ch <- protocol.NewResponseWithError(err)
		return
	}

	if result == nil {
		ch <- protocol.NewResponseWithNilBulk()
		return
	}

	ch <- result
}

// PrintSupportedCmd call on debug mode
func PrintSupportedCmd() {
	log.Println("[redis/command] supported cmd count: ", len(implementedCommands))
	for c := range implementedCommands {
		log.Println("[redis/command] support: ", c)
	}
}
