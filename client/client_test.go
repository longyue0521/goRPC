package client_test

import (
	"context"
	"testing"

	"github.com/longyue0521/goRPC/client"
	"github.com/longyue0521/goRPC/proxy"
	"github.com/stretchr/testify/assert"
)

func TestClient_Init_VerifyParameters(t *testing.T) {
	testCases := map[string]struct {
		service client.Service
		pxy     proxy.Proxy
		wantErr error
	}{
		"nil": {
			service: nil,
			wantErr: client.ErrInvalidArgument,
		},
		"non pointer": {
			service: Integer(1),
			wantErr: client.ErrInvalidArgument,
		},
		"non pointer of struct": {
			service: func() client.Service { a := Integer(1); return &a }(),
			wantErr: client.ErrInvalidArgument,
		},
		"(*struct)(nil)": {
			service: (*UserService)(nil),
			wantErr: client.ErrInvalidArgument,
		},
		"service, with nil Proxy": {
			service: &UserService{},
			pxy:     nil,
			wantErr: client.ErrInvalidArgument,
		},
		"service, with typed nil Proxy": {
			service: &UserService{},
			pxy:     (*mockProxy)(nil),
			wantErr: client.ErrInvalidArgument,
		},
		"service, proxy": {
			service: &UserService{},
			pxy:     &mockProxy{},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := client.Init(tc.service, tc.pxy)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

type Integer int

func (n Integer) Name() string {
	return "Integer"
}

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u *UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
}

type GetByIdResp struct {
}

type mockProxy struct {
}

func (m *mockProxy) Invoke(ctx context.Context, req *proxy.Request) (resp *proxy.Response, err error) {
	// TODO implement me
	panic("implement me")
}
