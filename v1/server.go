package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
)

type Server struct {
	services map[string]reflectionStub
}

func NewServer() *Server {
	return &Server{
		services: map[string]reflectionStub{},
	}
}

func (s *Server) Register(service Service) {
	s.services[service.Name()] = reflectionStub{
		value: reflect.ValueOf(service),
	}
}

func (s *Server) Start(network string, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		return err
	}

	for {
		// 接收连接
		conn, err2 := l.Accept()
		if err2 != nil {
			continue
		}

		// 请求解析
		go func() {
			if err3 := s.handleConn(conn); err3 != nil {
				conn.Close()
				return
			}
		}()
	}

}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		// 1.读取请求
		reqMsg, err := ReadMsg(conn)
		if err != nil {
			return err
		}

		req := &Request{}

		err = json.Unmarshal(reqMsg, req)
		if err != nil {
			return err
		}

		fmt.Printf("server req: %v \n", req)

		// 2.业务执行
		// 2.1 找到本地服务
		service, ok := s.services[req.ServiceName]
		if !ok {
			// 2.4 构造错误信息
			return errors.New("服务不存在")
		}

		// 2.2 传入参数
		// context
		ctx := context.Background()
		res := service.invoke(ctx, req.MethodName, req.Data)

		// 3. 序列化响应内容, 写回响应
		// 第一个返回值
		methodResp := res[0].Interface()
		// 第二个返回值
		//methodErr := resp[1].Interface()

		return SendMsg(conn, methodResp)

	}
}

type reflectionStub struct {
	value reflect.Value
}

func (r *reflectionStub) invoke(ctx context.Context, methodName string, data []byte) []reflect.Value {
	method := r.value.MethodByName(methodName)
	methodReq := reflect.New(method.Type().In(1))

	err := json.Unmarshal(data, methodReq.Interface())
	if err != nil {
		return nil
	}

	return method.Call([]reflect.Value{reflect.ValueOf(ctx), methodReq.Elem()})
}
