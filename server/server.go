package server

import (
	"context"
)

type IServer interface {
	Start(ctx context.Context, addr string) error
	Stop()
}
