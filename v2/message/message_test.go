package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecodeRequest(t *testing.T) {

	testCases := []struct {
		name string

		req *Request
	}{
		// 带有meta数据
		{
			name: "with meta",
			req: &Request{
				MessageId:   1,
				Version:     1,
				Compress:    1,
				Serializer:  1,
				ServiceName: "test",
				MethodName:  "ljk",
				Meta: map[string]string{
					"key1": "val1",
					"key2": "val2",
					"key3": "val3",
				},
				Data: []byte("123456"),
			},
		},
		// 没有meta数据
		{
			name: "no meta",
			req: &Request{
				MessageId:   2,
				Version:     2,
				Compress:    2,
				Serializer:  2,
				ServiceName: "test",
				MethodName:  "ljk2",
				Data:        []byte("123456"),
			},
		},
		// 没有data
		{
			name: "no data",
			req: &Request{
				MessageId:   2,
				Version:     2,
				Compress:    2,
				Serializer:  2,
				ServiceName: "test",
				MethodName:  "ljk3",
				Meta: map[string]string{
					"key1": "val1",
					"key2": "val2",
					"key3": "val3",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.req.CalHeadLength()
			tc.req.BodyLength = uint32(len(tc.req.Data))

			bs := EncodeReq(tc.req)
			req := DecodeReq(bs)
			assert.Equal(t, tc.req, req)
		})
	}

}

func TestEncodeDecodeResponse(t *testing.T) {

	testCases := []struct {
		name string

		resp *Response
	}{
		// err + data
		{
			name: "with err + data",
			resp: &Response{
				MessageId:  1,
				Version:    1,
				Compress:   1,
				Serializer: 1,
				Error:      []byte("test err"),
				Data:       []byte("123456"),
			},
		},
		// err 没有data
		{
			name: "with err",
			resp: &Response{
				MessageId:  1,
				Version:    1,
				Compress:   1,
				Serializer: 1,
				Error:      []byte("test err"),
			},
		},
		// data 没有err
		{
			name: "with data",
			resp: &Response{
				MessageId:  1,
				Version:    1,
				Compress:   1,
				Serializer: 1,
				Data:       []byte("123456"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.resp.CalHeadLength()
			tc.resp.BodyLength = uint32(len(tc.resp.Data))

			bs := EncodeResponse(tc.resp)
			resp := DecodeResponse(bs)
			assert.Equal(t, tc.resp, resp)
		})
	}

}
