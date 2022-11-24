package client

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/longyue0521/goRPC/proxy"
)

var (
	ErrInvalidArgument = errors.New("client: invalid argument")
)

type Service interface {
	Name() string
}

func Init(service Service, p proxy.Proxy) error {

	val := reflect.ValueOf(service)
	typ := reflect.TypeOf(service)

	// service == nil时，reflect.ValueOf(service).IsValid() == false
	// service == (*struct)(nil) reflect.ValueOf(service).IsValid() == true reflect.ValueOf(service).IsNil() == true
	if service == nil || typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct || val.IsNil() {
		return fmt.Errorf("%w: service must be a pointer of struct", ErrInvalidArgument)
	}

	val = val.Elem()
	typ = typ.Elem()

	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		if fieldType.Type.Kind() != reflect.Func {
			continue
		}

		fn := reflect.MakeFunc(fieldType.Type, func(args []reflect.Value) (results []reflect.Value) {

			req := &proxy.Request{
				ServiceName: service.Name(),
				MethodName:  fieldType.Name,
				Arg:         args[1].Interface(),
			}

			for j := 0; j < fieldType.Type.NumOut(); j++ {
				results = append(results, reflect.New(fieldValue.Type().Out(j)).Elem())
			}
			log.Println(req)
			return
		})
		fieldValue.Set(fn)
	}

	return nil
}
