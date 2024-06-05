package runtimex

import (
	"context"
	"reflect"
)

func Call(ctx context.Context, fn interface{}, params ...interface{}) (result []reflect.Value) {
	f := reflect.ValueOf(fn)
	in := make([]reflect.Value, len(params)+1)
	in[0] = reflect.ValueOf(ctx)
	for i, param := range params {
		in[i+1] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
