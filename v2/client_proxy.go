package v1

import (
	"context"
	"github.com/silenceper/pool"
	"my_rpc/v2/message"
	"net"
	"time"
)

type Client struct {
	connPool pool.Pool // 客户端连接池
}

func NewClient(network string, addr string) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap: 10,  // 初始容量
		MaxCap:     100, // 最大连接数
		MaxIdle:    50,  //最大空闲数
		Factory: func() (interface{}, error) {
			return net.Dial(network, addr)
		},
		IdleTimeout: time.Minute, //空闲1min就关闭
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		connPool: p,
	}, nil
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {

	// 取出连接
	conObj, err := c.connPool.Get()
	// 注意区分框架err 和 用户err
	if err != nil {
		return nil, err
	}
	conn := conObj.(net.Conn)

	// 协议约定的几种方式：
	// 1. 长度+数据
	// 2. 命令字+结构体
	err = SendMsg(conn, req)
	if err != nil {
		return nil, err
	}

	// 读取响应
	// 1.1 先读长度
	// 1.2 后读数据
	respMsg, err := ReadMsg(conn)
	if err != nil {
		return nil, err
	}

	return message.DecodeResponse(respMsg), nil
}
