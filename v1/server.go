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
	services map[string]Service
}

func NewServer() *Server {
	return &Server{
		services: map[string]Service{},
	}
}

func (s *Server) Register(service Service) {
	s.services[service.Name()] = service
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
		// 待执行方法指针
		method := reflect.ValueOf(service).MethodByName(req.MethodName)

		// 2.2 传入参数
		// 第一个入参
		ctx := context.Background()
		// 第二个入参
		methodReq := reflect.New(method.Type().In(1))

		// 反序列化参数值
		err = json.Unmarshal(req.Data, methodReq.Interface())
		if err != nil {
			//TODO: 构造错误信息
			return errors.New("参数值解析错误")
		}

		// 2.3 执行业务逻辑
		resp := method.Call([]reflect.Value{reflect.ValueOf(ctx), methodReq.Elem()})

		// 3. 序列化响应内容, 写回响应
		// 第一个返回值
		methodResp := resp[0].Interface()
		// 第二个返回值
		//methodErr := resp[1].Interface()

		return SendMsg(conn, methodResp)

	}
}
