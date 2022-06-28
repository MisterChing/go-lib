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
