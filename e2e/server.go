package e2e

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yusank/godis/handler"
	"github.com/yusank/godis/server"
)

func startServer(addr string, t *testing.T) *server.Server {
	srv := server.NewServer(addr, nil, handler.TCPHandler{})

	go func() {
		if err := srv.Start(); err != nil {
			assert.NoError(t, err)
		}
	}()

	time.Sleep(time.Second)
	return srv
}
