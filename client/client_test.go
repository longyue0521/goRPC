package client_test

import (
	"context"
	"testing"

	"github.com/longyue0521/goRPC/client"
	"github.com/longyue0521/goRPC/proxy"
	"github.com/stretchr/testify/assert"
)

func TestClient_Init(t *testing.T) {
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
		"user service": {
			service: &UserService{},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := client.Init(tc.service, tc.pxy)
			assert.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr != nil {
				return
			}
			us, _ := tc.service.(*UserService)
			_, _ = us.GetById(context.Background(), &GetByIdReq{})
		})
	}
}

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

type Integer int

func (n Integer) Name() string {
	return "Integer"
}

func (u *UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
}

type GetByIdResp struct {
}
