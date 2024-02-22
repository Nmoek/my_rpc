package v1

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
)

func InitClientProxy(service Service, p Proxy) error {

	// 获取字段值
	val := reflect.ValueOf(service).Elem()

	// 获取字段类型(指针类型)
	typ := reflect.TypeOf(service).Elem()

	// 获取字段数量
	numField := val.NumField()

	for i := 0; i < numField; i++ {

		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		// 如果怎么不为可以赋值类型就跳过
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
			ctx, ok := args[0].Interface().(context.Context)
			if !ok {
				results = append(results, reflect.ValueOf(errors.New("args[0] isn't context")))
				return
			}
			arg := args[1].Interface()

			req := &Request{
				ServiceName: service.Name(),
				// 字段名称字符串
				MethodName: fieldType.Name,
				// 从第二个参数开始取
				Arg: arg,
			}

			// 发送请求到服务端
			resp, err := p.Invoke(ctx, req)
			// 获取第一个返回值对象
			outType := fieldType.Type.Out(0)
			if err != nil {
				results = append(results, reflect.Zero(outType))
				results = append(results, reflect.ValueOf(err))
				return
			}

			// 拿到第一个返回值类型 *GetByIdResp
			first := reflect.New(outType).Interface()

			// 填充返回值列表
			err = json.Unmarshal(resp.Data, first)
			// 需要完成resp.Data --> first
			results = append(results, reflect.ValueOf(first).Elem())

			if err != nil {
				results = append(results, reflect.ValueOf(err))
			} else {
				// 注意: 这里写法意思为 类型为error的nil值
				results = append(results, reflect.Zero(reflect.TypeOf(new(error)).Elem()))
			}

			//numOut := fieldType.Type.NumOut()
			//for j := 0; j < numOut; j++ {
			//	results = append(results, reflect.New(fieldType.Type.Out(j)).Elem())
			//}

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

		fieldValue.Set(fn)

	}

	return nil
}

type Service interface {
	Name() string
}
