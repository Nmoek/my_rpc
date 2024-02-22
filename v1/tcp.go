package v1

import (
	"encoding/binary"
	"errors"
	"math"
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

func DecodeMsg(data []byte) ([]byte, error) {
	if float64(len(data)) >= math.Pow(2, 8) {
		return nil, errors.New("data too long")
	}

	resp := make([]byte, 0, len(data)+lengthByte)
	l := len(data)

	binary.BigEndian.PutUint64(resp, uint64(l))
	copy(resp[lengthByte:], data)

	return resp, nil
}
