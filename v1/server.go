package v1

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Server struct {
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
		// 1.1 先读长度
		// 1.2 后读数据
		lengthBytes := make([]byte, lengthByte)
		_, err := conn.Read(lengthBytes)
		if err != nil {
			return err
		}

		dataLen := binary.BigEndian.Uint64(lengthBytes)

		reqMsg := make([]byte, dataLen)
		_, err = conn.Read(reqMsg)
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

		// 3.写回响应

	}
}
