package confguard

import (
	kencoding "github.com/go-kratos/kratos/v2/encoding"
	"log"
	"reflect"
	"sync"
)

type Guarder struct {
	child interface{}
	rw    sync.RWMutex
}

func NewGuarder(child interface{}) *Guarder {
	obj := &Guarder{
		child: child,
	}
	return obj
}

func (gdr *Guarder) Get() interface{} {
	gdr.rw.RLock()
	defer gdr.rw.RUnlock()
	return gdr.child
}

func (gdr *Guarder) Lock() {
	log.Println("lock")
	gdr.rw.Lock()
}

func (gdr *Guarder) Unlock() {
	log.Println("unlock")
	gdr.rw.Unlock()
}

//populateWhenPtr  只适用于child实际类型是ptr类型
func (gdr *Guarder) populateWhenPtr(input []byte) error {
	rv := reflect.ValueOf(gdr)

	//childField := rv.Elem().FieldByName("child")
	//childFieldContent := childField.Elem()
	//debugutil.DebugPrintV2("childFiled",
	//	childField,
	//	childField.Kind(),
	//	childField.Type(),
	//	//childField.Pointer(),
	//	//childField.UnsafePointer(),
	//	childField.Addr(),
	//	childField.UnsafeAddr(),
	//	childField.Addr().Pointer(),
	//)
	//debugutil.DebugPrintV2("childFiledContent",
	//	childFieldContent,
	//	childFieldContent.Kind(),
	//	childFieldContent.Type(),
	//	childFieldContent.Pointer(),
	//	childFieldContent.UnsafePointer(),
	//	//childAddressableDebug.Addr(),
	//	//childAddressableDebug.UnsafeAddr(),
	//	//childAddressableDebug.Addr().Pointer(),
	//)

	//根据child类型创建一个copy的ptr   child --> interface{} --> ptr --> struct --> ptr
	childField := rv.Elem().FieldByName("child")
	childFieldContent := childField.Elem()                                   //ptr
	childContentUnderlyingType := childFieldContent.Elem().Type()            //struct
	childContentUnderlyingCopyPtr := reflect.New(childContentUnderlyingType) //ptr of copied struct

	if err := kencoding.GetCodec("json").Unmarshal(input, childContentUnderlyingCopyPtr.Interface()); err != nil {
		return err
	}

	//debugutil.DebugPrintV2("childContentUnderlyingCopyPtr",
	//	childContentUnderlyingCopyPtr,
	//	childContentUnderlyingCopyPtr.Kind(),
	//	childContentUnderlyingCopyPtr.Type(),
	//	childContentUnderlyingCopyPtr.Interface(),
	//	childContentUnderlyingCopyPtr.Pointer(),
	//	childContentUnderlyingCopyPtr.UnsafePointer(),
	//
	//	//childContentUnderlyingCopyPtr.Addr(),
	//	//childContentUnderlyingCopyPtr.UnsafeAddr(),
	//	//unsafe.Pointer(childContentUnderlyingCopyPtr.UnsafeAddr()),
	//)

	rv.MethodByName("Lock").Call(nil) //lock
	//rv.Elem().FieldByName("child").Set(addressableChildCopy.Elem()) //可导出字段可用，非导出字段panic
	//获取非导出字段并转换为可寻址的ptr
	childField2 := rv.Elem().FieldByName("child")
	childAddressablePtr := reflect.NewAt(childField2.Elem().Elem().Type(), childField2.Elem().UnsafePointer())
	//childAddressablePtr := reflect.NewAt(childAddressable.Type(), unsafe.Pointer(childAddressable.UnsafeAddr()))
	//debugutil.DebugPrintV2("childAddressablePtr",
	//	childAddressablePtr,
	//	childAddressablePtr.Kind(),
	//	childAddressablePtr.Type(),
	//	childAddressablePtr.Interface(),
	//	childAddressablePtr.Pointer(),
	//	childAddressablePtr.UnsafePointer(),
	//	//childAddressablePtr.Addr(),
	//	//childAddressablePtr.UnsafeAddr(),
	//	//unsafe.Pointer(childAddressablePtr.UnsafeAddr()),
	//)
	childAddressablePtr.Elem().Set(childContentUnderlyingCopyPtr.Elem())
	rv.MethodByName("Unlock").Call(nil) //unlock
	return nil
}
