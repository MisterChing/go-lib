package extra

import (
	"encoding/json"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

type fuzzyStringDecoder struct {
}

func (decoder *fuzzyStringDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		*((*string)(ptr)) = string(number)
	case jsoniter.StringValue:
		*((*string)(ptr)) = iter.ReadString()
	case jsoniter.NilValue:
		iter.Skip()
		*((*string)(ptr)) = ""
	case jsoniter.ArrayValue:
		iter.Skip()
		*((*string)(ptr)) = ""
	case jsoniter.ObjectValue:
		iter.Skip()
		*((*string)(ptr)) = ""
	default:
		iter.ReportError("fuzzyStringDecoder", "not number or string")
	}
}
