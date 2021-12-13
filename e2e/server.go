package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	v1 "github.com/yusank/godis/server/v1"
)

func startServer(addr string, t *testing.T) *v1.Server {
	srv := v1.NewServer()

	go func() {
		if err := srv.Start(context.Background(), addr); err != nil {
			assert.NoError(t, err)
		}
	}()

	time.Sleep(time.Second)
	return srv
}
