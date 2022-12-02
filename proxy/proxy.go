package proxy

import (
	"context"

	"github.com/longyue0521/goRPC/transport"
)

type Proxy interface {
	Invoke(ctx context.Context, req *transport.Request) (resp *transport.Response, err error)
}

