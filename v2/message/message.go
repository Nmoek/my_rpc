package message

import (
	"bytes"
	"encoding/binary"
)

const (
	splitter     = '\n'
	pariSplitter = '\r'
)

type Request struct {
	HeadLength uint32
	BodyLength uint32

	MessageId uint32
	Version   byte
	// 压缩算法
	Compress byte
	// 序列化协议
	Serializer byte
	//padding byte // 本地CPU高速缓存字节对齐(不发送到对端)

	ServiceName string
	MethodName  string

	Meta map[string]string

	Data []byte
}

func EncodeReq(req *Request) []byte {
	bs := make([]byte, req.HeadLength+req.BodyLength)
	cur := bs

	binary.BigEndian.PutUint32(cur[:4], req.HeadLength)
	cur = cur[4:]

	binary.BigEndian.PutUint32(cur[:4], req.BodyLength)
	cur = cur[4:]

	binary.BigEndian.PutUint32(cur[:4], req.MessageId)
	cur = cur[4:]

	cur[0] = req.Version
	cur = cur[1:]
	cur[0] = req.Compress
	cur = cur[1:]
	cur[0] = req.Serializer
	cur = cur[1:]

	copy(cur, req.ServiceName)
	cur[len(req.ServiceName)] = splitter
	cur = cur[len(req.ServiceName)+1:]
	copy(cur, req.MethodName)
	cur[len(req.MethodName)] = splitter
	cur = cur[len(req.MethodName)+1:]

	for key, val := range req.Meta {
		copy(cur, key)
		cur[len(key)] = pariSplitter
		cur = cur[len(key)+1:]

		copy(cur, val)
		cur[len(val)] = splitter
		cur = cur[len(val)+1:]
	}

	copy(bs[req.HeadLength:], req.Data)

	return bs
}
func DecodeReq(data []byte) *Request {
	req := &Request{}

	req.HeadLength = binary.BigEndian.Uint32(data[:4])
	req.BodyLength = binary.BigEndian.Uint32(data[4:8])
	req.MessageId = binary.BigEndian.Uint32(data[8:12])

	req.Version = data[12]
	req.Compress = data[13]
	req.Serializer = data[14]

	retainData := data[15:req.HeadLength]
	index := bytes.IndexByte(retainData, splitter)
	req.ServiceName = string(retainData[:index])
	retainData = retainData[index+1:]

	index = bytes.IndexByte(retainData, splitter)
	req.MethodName = string(retainData[:index])
	retainData = retainData[index+1:]

	if len(retainData) > 0 {
		meta := map[string]string{}
		index = bytes.IndexByte(retainData, splitter)
		for index != -1 {
			pairIdx := bytes.IndexByte(retainData, pariSplitter)
			meta[string(retainData[:pairIdx])] = string(retainData[pairIdx+1 : index])
			retainData = retainData[index+1:]
			index = bytes.IndexByte(retainData, splitter)

		}
		req.Meta = meta
	}
	if req.BodyLength > 0 {
		req.Data = data[req.HeadLength:]
	}

	return req
}

func (r *Request) CalHeadLength() {
	// 固定15Byte
	res := 15
	res += len(r.ServiceName)
	res += 1 // 注意加入bit分隔符
	res += len(r.MethodName)
	res += 1 // 注意加入bit分隔符

	for key, val := range r.Meta {
		res += len(key) + 1 + len(val) + 1 // +1区分不同的键值对
	}

	r.HeadLength = uint32(res)

}

type Response struct {
	HeadLength uint32
	BodyLength uint32

	MessageId uint32
	Version   byte
	// 压缩算法
	Compress byte
	// 序列化协议
	Serializer byte

	Error []byte

	Data []byte
}

func EncodeResponse(resp *Response) []byte {
	bs := make([]byte, resp.HeadLength+resp.BodyLength)
	cur := bs

	binary.BigEndian.PutUint32(cur[:4], resp.HeadLength)
	cur = cur[4:]

	binary.BigEndian.PutUint32(cur[:4], resp.BodyLength)
	cur = cur[4:]

	binary.BigEndian.PutUint32(cur[:4], resp.MessageId)
	cur = cur[4:]

	cur[0] = resp.Version
	cur = cur[1:]

	cur[0] = resp.Compress
	cur = cur[1:]

	cur[0] = resp.Serializer
	cur = cur[1:]

	copy(cur, resp.Error)
	cur = cur[len(resp.Error):]

	copy(cur, resp.Data)

	return bs
}

func DecodeResponse(data []byte) *Response {
	resp := &Response{}

	resp.HeadLength = binary.BigEndian.Uint32(data[:4])
	resp.BodyLength = binary.BigEndian.Uint32(data[4:8])
	resp.MessageId = binary.BigEndian.Uint32(data[8:12])

	resp.Version = data[12]
	resp.Compress = data[13]
	resp.Serializer = data[14]

	if len(data[15:resp.HeadLength]) > 0 {
		resp.Error = data[15:resp.HeadLength]
	}
	if resp.BodyLength > 0 {
		resp.Data = data[resp.HeadLength:]
	}

	return resp
}

func (r *Response) CalHeadLength() {
	r.HeadLength = 15 + uint32(len(r.Error))
}
