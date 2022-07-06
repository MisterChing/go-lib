package interval_manager

import (
	"reflect"
	"unsafe"

	"github.com/gin-gonic/gin"
)

var fakeEngine = &gin.Engine{}

func genGinContext() *gin.Context {
	var ctx = &gin.Context{}
	v := reflect.ValueOf(ctx).Elem().FieldByName("engine")
	rv := reflect.ValueOf(fakeEngine)
	reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Set(rv)
	return ctx
}

func FetchMessage(q <-chan interface{}, batchSize int) []interface{} {
	var (
		ret = make([]interface{}, 0, batchSize)
	)
	for i := 0; i < batchSize; i++ {
		if msg := readChanNoBlock(q); msg != nil {
			ret = append(ret, msg)
		}
	}
	return ret
}

func readChanNoBlock(q <-chan interface{}) interface{} {
	select {
	case v, ok := <-q:
		if ok {
			return v
		}
		return nil
	default:
		return nil
	}
}
