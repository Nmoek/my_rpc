//go:build e2e

package v1

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServer_Start(t *testing.T) {

	s := NewServer()
	s.Register(&UserServiceServer{})
	err := s.Start("tcp", ":8888")
	assert.NoError(t, err)
}

type UserServiceServer struct {
}

func (u *UserServiceServer) Name() string {
	return "user-service"
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	return &GetByIdResp{
		Name: "ljk",
	}, nil
}
