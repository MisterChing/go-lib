package extra

import (
	"encoding/json"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

type fuzzyBoolDecoder struct {
}

func (decoder *fuzzyBoolDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.BoolValue:
		*((*bool)(ptr)) = iter.ReadBool()
	case jsoniter.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		if number == "1" {
			*((*bool)(ptr)) = true
		} else if number == "0" {
			*((*bool)(ptr)) = false
		} else {
			*((*bool)(ptr)) = false
		}
	case jsoniter.StringValue:
		str := iter.ReadString()
		if str == "true" || str == "1" {
			*((*bool)(ptr)) = true
		} else if str == "false" || str == "0" {
			*((*bool)(ptr)) = false
		} else {
			*((*bool)(ptr)) = false
		}
	case jsoniter.NilValue:
		iter.Skip()
		*((*bool)(ptr)) = false
	case jsoniter.ArrayValue:
		iter.Skip()
		*((*bool)(ptr)) = false
	case jsoniter.ObjectValue:
		iter.Skip()
		*((*bool)(ptr)) = false
	default:
		iter.ReportError("fuzzyBoolDecoder", "convert bool failed")
	}
}
