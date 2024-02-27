package v1

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"net"
)

// 使用多少字节表达有效数据的长度
const lengthByte = 8

func EncodeMsg(data []byte) ([]byte, error) {
	if float64(len(data)) >= math.Pow(2, 8) {
		return nil, errors.New("data too long")
	}

	resp := make([]byte, len(data)+lengthByte)
	l := len(data)

	binary.BigEndian.PutUint64(resp, uint64(l))
	copy(resp[lengthByte:], data)

	return resp, nil
}

func SendMsg(conn net.Conn, msg any) error {

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	data, err = EncodeMsg(data)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)

	return err
}

func ReadMsg(conn net.Conn) ([]byte, error) {
	lengthBytes := make([]byte, lengthByte)
	_, err := conn.Read(lengthBytes)
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint64(lengthBytes)

	msg := make([]byte, dataLen)
	_, err = conn.Read(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
