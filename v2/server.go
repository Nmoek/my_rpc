package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"my_rpc/v2/message"
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

		req := message.DecodeReq(reqMsg)

		resp := &message.Response{
			MessageId: req.MessageId,

			Version:    req.Version,
			Compress:   req.Compress,
			Serializer: req.Serializer,
		}
		fmt.Printf("server req: %v \n", req)

		// 2.业务执行
		// 2.1 找到本地服务
		service, ok := s.services[req.ServiceName]
		if !ok {
			// 2.4 构造错误信息
			resp.CalHeadLength()
			resp.Error = []byte("服务不存在")
			err = SendMsg(conn, resp)
			if err != nil {
				return err
			}
			continue
		}

		// 2.2 传入参数
		// context
		ctx := context.Background()
		data, err := service.invoke(ctx, req.MethodName, req.Data)
		if err != nil {
			resp.CalHeadLength()
			resp.Error = []byte(err.Error())
			err = SendMsg(conn, resp)
			if err != nil {
				return err
			}
			continue
		}

		resp.CalHeadLength()
		resp.Data = data
		resp.BodyLength = uint32(len(data))

		return SendMsg(conn, resp)

	}
}

type reflectionStub struct {
	value reflect.Value
}

func (r *reflectionStub) invoke(ctx context.Context, methodName string, data []byte) ([]byte, error) {
	method := r.value.MethodByName(methodName)
	methodReq := reflect.New(method.Type().In(1))

	res := method.Call([]reflect.Value{reflect.ValueOf(ctx), methodReq.Elem()})
	if len(res) > 1 && !res[1].IsZero() {
		return nil, res[1].Interface().(error)
	}
	return json.Marshal(res[0].Interface())
}
