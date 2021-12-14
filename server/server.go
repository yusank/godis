package server

import (
	"context"
)

type IServer interface {
	Start(addr string) error
	Stop(ctx context.Context)
}
