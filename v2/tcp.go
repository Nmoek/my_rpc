package v1

import (
	"encoding/binary"
	"errors"
	"my_rpc/v2/message"
	"net"
)

const lengthBytes = 8

func SendMsg(conn net.Conn, msg any) error {
	var data []byte

	switch msg.(type) {
	case *message.Request:
		req := msg.(*message.Request)
		data = message.EncodeReq(req)
	case *message.Response:
		resp := msg.(*message.Response)
		data = message.EncodeResponse(resp)
	default:
		return errors.New("msg type invalid")
	}

	_, err := conn.Write(data)

	return err
}

func ReadMsg(conn net.Conn) ([]byte, error) {
	msgLenBytes := make([]byte, lengthBytes)
	_, err := conn.Read(msgLenBytes)
	if err != nil {
		return nil, err
	}
	headLength := binary.BigEndian.Uint32(msgLenBytes[:4])
	bodyLength := binary.BigEndian.Uint32(msgLenBytes[4:8])

	msg := make([]byte, headLength+bodyLength)
	_, err = conn.Read(msg[lengthBytes:])
	if err != nil {
		return nil, err
	}

	copy(msg, msgLenBytes)

	return msg, nil
}
