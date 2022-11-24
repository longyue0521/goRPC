package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/longyue0521/goRPC/proxy"
)

var (
	ErrInvalidArgument        = errors.New("client: invalid argument")
	ErrFailedToDecodeResponse = errors.New("client: failed to decode response")
)

type Service interface {
	Name() string
}

func Init(service Service, p proxy.Proxy) error {

	if p == nil || reflect.ValueOf(p).IsNil() {
		return fmt.Errorf("%w: proxy shoud not be nil", ErrInvalidArgument)
	}

	val := reflect.ValueOf(service)
	typ := reflect.TypeOf(service)

	// service == nil时，reflect.ValueOf(service).IsValid() == false
	// service == (*struct)(nil) reflect.ValueOf(service).IsValid() == true reflect.ValueOf(service).IsNil() == true
	if service == nil || typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct || val.IsNil() {
		return fmt.Errorf("%w: service should be a pointer of struct", ErrInvalidArgument)
	}

	val = val.Elem()
	typ = typ.Elem()

	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanSet() {
			// todo: return error
			continue
		}

		if fieldType.Type.Kind() != reflect.Func {
			// todo: return error
			continue
		}

		// todo: check signature of func
		// fieldType.Type.NumIn() == 2 && first is context.Context, second is pointer of Struct
		// fieldType.Type.NumOut() ==2 && first is pointer of Struct type, second is error type

		fn := reflect.MakeFunc(fieldType.Type, func(args []reflect.Value) (results []reflect.Value) {

			serviceRespType := fieldType.Type.Out(0)
			
			arg := args[1].Interface()
			ctx, _ := args[0].Interface().(context.Context)

			req := &proxy.Request{
				ServiceName: service.Name(),
				MethodName:  fieldType.Name,
				Arg:         arg,
			}

			resp, err := p.Invoke(ctx, req)
			if err != nil {
				results = append(results, reflect.New(serviceRespType).Elem())
				results = append(results, reflect.ValueOf(err))
				return
			}

			// convert resp to XXXResp
			serviceResp := reflect.New(serviceRespType).Interface()
			err = json.Unmarshal(resp.Result, serviceResp)
			if err != nil {
				results = append(results, reflect.New(serviceRespType).Elem())
				results = append(results, reflect.ValueOf(ErrFailedToDecodeResponse))
				return
			}

			results = append(results, reflect.ValueOf(serviceResp).Elem())
			results = append(results, reflect.ValueOf(new(error)).Elem())
			return
		})
		fieldValue.Set(fn)
	}

	return nil
}
