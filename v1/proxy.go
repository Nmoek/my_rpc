package v1

import "context"

type Proxy interface {
	Invoke(ctx context.Context, req *Request) (*Response, error)
}

type Request struct {
	ServiceName string
	MethodName  string
	Arg         any
}
type Response struct {
	Data []byte
}
