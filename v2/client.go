package v1

import (
	"context"
	"encoding/json"
	"errors"
	"my_rpc/v2/message"
	"reflect"
	"sync/atomic"
)

var messageId uint32 = 0

type Service interface {
	Name() string
}

func InitClientProxy(service Service, p Proxy) error {

	// 获取字段类型(指针类型)
	typ := reflect.TypeOf(service).Elem()

	// 获取字段值
	val := reflect.ValueOf(service).Elem()

	// 获取字段数量
	numField := val.NumField()

	for i := 0; i < numField; i++ {

		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		// 如果不为可以赋值类型就跳过
		if !fieldValue.CanSet() {
			continue
		}

		// 字段类型是否是函数指针类型
		if fieldType.Type.Kind() != reflect.Func {
			continue
		}

		// 在这里替换新的方法实现
		fn := reflect.MakeFunc(fieldType.Type, func(args []reflect.Value) (results []reflect.Value) {
			// 在这里需要拼接调用信息:
			// 服务名、方法名、参数值
			// 没有严格的参数类型匹配，要小心一些

			/*1. 解析+重组调用信息、入参信息*/
			// 获取第一个参数
			ctx, ok := args[0].Interface().(context.Context)
			if !ok {
				results = append(results, reflect.ValueOf(errors.New("args[0] isn't context")))
				return
			}
			// 获取第一个返回值对象
			outType := fieldType.Type.Out(0)
			// 获取第二个参数
			arg := args[1].Interface()

			data, err := json.Marshal(arg)
			if err != nil {
				results = append(results, reflect.Zero(outType))
				results = append(results, reflect.ValueOf(err))
				return
			}
			msgId := atomic.AddUint32(&messageId, 1)

			req := &message.Request{

				BodyLength: uint32(len(data)),

				Version:    0,
				Compress:   0,
				Serializer: 0,
				MessageId:  msgId,
				// 服务名
				ServiceName: service.Name(),
				// 调用方法名 字段名称字符串
				MethodName: fieldType.Name,
				// 参数值
				Data: data,
			}
			// 计算得到头部长度
			req.CalHeadLength()

			/*2. 发送请求到服务端，并接收响应*/
			resp, err := p.Invoke(ctx, req)
			if err != nil {
				results = append(results, reflect.Zero(outType))
				results = append(results, reflect.ValueOf(err))
				return
			}

			/*3. 解析+重组响应内容，返回真正的返回值*/
			// 拿到第一个返回值类型 *GetByIdResp
			first := reflect.New(outType).Interface()

			// 需要完成resp.Data --> first
			err = json.Unmarshal(resp.Data, first)

			// 填充返回值列表
			results = append(results, reflect.ValueOf(first).Elem())

			if err != nil {
				results = append(results, reflect.ValueOf(err))
			} else {
				// 注意: 这里写法意思为 类型为error的nil值
				results = append(results, reflect.Zero(reflect.TypeOf(new(error)).Elem()))
			}

			// 以下为测试结果探针
			//str, err := json.Marshal(req)
			//if err != nil {
			//	results = append(results, reflect.ValueOf(err))
			//	return
			//}
			//results = append(results, reflect.New(fieldType.Type.Out(0).Elem()))
			//
			//results = append(results, reflect.ValueOf(errors.New(string(str))))
			return
		})

		// 将实现好的函数替换回函数指针中
		fieldValue.Set(fn)

	}

	return nil
}
