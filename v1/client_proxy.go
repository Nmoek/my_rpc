package v1

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/silenceper/pool"
	"net"
	"time"
)

type Client struct {
	connPool pool.Pool
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

func (c *Client) Invoke(ctx context.Context, req *Request) (*Response, error) {

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

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

	// 网络大端编码
	data, err = EncodeMsg(data)
	if err != nil {
		return nil, err
	}

	// 发送请求
	i, err := conn.Write(data)
	if err != nil {
		return nil, err
	}
	if i != len(data) {
		return nil, errors.New("数据未全全部写入完成")
	}

	// 读取响应
	// TODO: 如何知道应该读取多长的响应?

	return nil, errors.New("todo client")
}
