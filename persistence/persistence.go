package persistence

import (
	"fmt"

	"github.com/yusank/godis/protocol"
)

type AppendOnlyFiles struct {
}

// WriteCommand write into file async
func (a *AppendOnlyFiles) WriteCommand(r *protocol.Receive) {
	fmt.Println(r.OrgStr)
}

var AOF *AppendOnlyFiles

var opts *options

type options struct {
}

type Option func(o *options)

func Init(opt ...Option) error {
	opts = new(options)

	for _, o := range opt {
		o(opts)
	}

	return nil
}
