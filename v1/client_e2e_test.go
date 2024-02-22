package v1

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Start(t *testing.T) {

	c, err := NewClient("tcp", ":8888")
	assert.NoError(t, err)

	us := &UserService{}
	err = InitClientProxy(us, c)
	assert.NoError(t, err)

	resp, err := us.GetById(context.Background(), &GetByIdReq{
		Id: 100,
	})
	fmt.Printf("resp:%v \n", resp)
	assert.NoError(t, err)

}
