package v1

import (
	"context"
	"my_rpc/v2/message"
)

type Proxy interface {
	Invoke(ctx context.Context, req *message.Request) (*message.Response, error)
}
