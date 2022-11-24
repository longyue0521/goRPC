package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/longyue0521/goRPC/client"
	"github.com/longyue0521/goRPC/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Init_VerifyParameters(t *testing.T) {
	testCases := map[string]struct {
		service client.Service
		pxy     proxy.Proxy
		wantErr error
	}{
		"service is nil": {
			service: nil,
			pxy:     &mockProxy{},
			wantErr: client.ErrInvalidArgument,
		},
		"service is non pointer": {
			service: Integer(1),
			pxy:     &mockProxy{},
			wantErr: client.ErrInvalidArgument,
		},
		"service is non pointer of struct": {
			service: func() client.Service { a := Integer(1); return &a }(),
			pxy:     &mockProxy{},
			wantErr: client.ErrInvalidArgument,
		},
		"service is (*struct)(nil)": {
			service: (*UserService)(nil),
			wantErr: client.ErrInvalidArgument,
		},
		"proxy is nil": {
			service: &UserService{},
			pxy:     nil,
			wantErr: client.ErrInvalidArgument,
		},
		"proxy is (*struct)(nil)": {
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

func TestClient_Init_VerifyService(t *testing.T) {
	testCases := map[string]struct {
		service *UserService
		req     *GetByIdReq
		resp    *GetByIdResp
		proxy   *mockProxy
		wantErr error
	}{
		"failed to invoke RPC": {
			service: &UserService{},
			req:     &GetByIdReq{Id: 13},
			resp:    &GetByIdResp{Name: "Go"},
			proxy:   &mockProxy{err: errors.New("failed to invoke RPC")},
			wantErr: errors.New("failed to invoke RPC"),
		},
		"success to invoke RPC and decode response": {
			service: &UserService{},
			req:     &GetByIdReq{Id: 13},
			resp:    &GetByIdResp{Name: "Go"},
			proxy: &mockProxy{resp: &proxy.Response{
				Result: func() []byte {
					b, err := json.Marshal(&GetByIdResp{
						Name: "Go",
					})
					require.NoError(t, err)
					return b
				}(),
			}},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := client.Init(tc.service, tc.proxy)
			require.NoError(t, err)

			resp, err := tc.service.GetById(context.Background(), tc.req)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.resp, resp)
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
	Id int64
}

type GetByIdResp struct {
	Name string `json:"Name"`
}

type PlayerService struct {
	GetById func(ctx context.Context, req *GetByPlayerIdReq) (*GetByPlayerIdResp, error)
}

func (u *PlayerService) Name() string {
	return "player-service"
}

type GetByPlayerIdReq struct {
	Id int64
}

type GetByPlayerIdResp struct {
	Name string
}

type mockProxy struct {
	req  *proxy.Request
	resp *proxy.Response
	err  error
}

func (m *mockProxy) Invoke(ctx context.Context, req *proxy.Request) (resp *proxy.Response, err error) {
	m.req = req
	return m.resp, m.err
}
