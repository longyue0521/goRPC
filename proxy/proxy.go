package proxy

import "context"

type Proxy interface {
	Invoke(ctx context.Context, req *Request) (resp *Response, err error)
}

type Request struct {
	ServiceName string
	MethodName  string
	// todo: ctx is ignored
	Arg any
}

type Response struct {
	Result []byte
}
