package v1

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitClientProxy(t *testing.T) {
	testCases := []struct {
		name string

		mock    *mockProxy
		service *UserServiceClient

		wantReq     *Request
		wantResp    *GetByIdResp
		wantInitErr error
		wantErr     error
	}{
		// 入参校验+返回值校验
		{
			name: "req and response",
			mock: &mockProxy{
				result: []byte(`{"name": "Tom"}`),
			},
			wantReq: &Request{
				ServiceName: "user-service",
				MethodName:  "GetById",
				Data:        []byte(`{"id": 13}`),
			},

			wantResp: &GetByIdResp{
				Name: "Tom",
			},
			service: &UserServiceClient{},
		},
		// proxy错误
		{
			name: "proxy return errors",
			mock: &mockProxy{
				err: errors.New("mock err"),
			},
			service: &UserServiceClient{},
			wantErr: errors.New("mock err"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			err := InitClientProxy(tc.service, tc.mock)
			assert.Equal(t, tc.wantInitErr, err)
			if err != nil {
				return
			}
			resp, err := tc.service.GetById(context.Background(), &GetByIdReq{Id: 13})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.mock.req, tc.wantReq)
			assert.Equal(t, tc.wantResp, resp)

		})
	}

}

// mockProxy
// @Description: 测试用Proxy实现
type mockProxy struct {
	req    *Request
	result []byte
	err    error
}

func (m *mockProxy) Invoke(ctx context.Context, req *Request) (*Response, error) {
	m.req = req
	return &Response{
		Data: m.result,
	}, m.err
}

// UserServiceClient
// @Description: 测试用Client实现
type UserServiceClient struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

type GetByIdReq struct {
	Id int64
}

type GetByIdResp struct {
	Name string `json:"name"`
}

func (u UserServiceClient) Name() string {
	return "user-service"
}
