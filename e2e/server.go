package e2e

import (
	"testing"
	"time"

	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"github.com/stretchr/testify/assert"

	"github.com/yusank/godis/server"
	v2 "github.com/yusank/godis/server/v2"
)

func startServer(addr string, t *testing.T) server.IServer {
	srv := v2.NewServer(goroutine.Default())
	go func() {
		if err := srv.Start(addr); err != nil {
			assert.NoError(t, err)
		}
	}()

	time.Sleep(time.Second)
	return srv
}
