package extra

import (
	"encoding/json"
	"io"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

type fuzzyIntegerDecoder struct {
	fun func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator)
}

func (decoder *fuzzyIntegerDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoniter.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		str = string(number)
	case jsoniter.StringValue:
		str = iter.ReadString()
	case jsoniter.BoolValue:
		if iter.ReadBool() {
			str = "1"
		} else {
			str = "0"
		}
	case jsoniter.NilValue:
		iter.Skip()
		str = "0"
	case jsoniter.ArrayValue:
		iter.Skip()
		str = "0"
	case jsoniter.ObjectValue:
		iter.Skip()
		str = "0"
	default:
		iter.ReportError("fuzzyIntegerDecoder", "not number or string")
	}
	if len(str) == 0 {
		str = "0"
	}
	newIter := iter.Pool().BorrowIterator([]byte(str))
	defer iter.Pool().ReturnIterator(newIter)
	isFloat := strings.IndexByte(str, '.') != -1
	decoder.fun(isFloat, ptr, newIter)
	if newIter.Error != nil && newIter.Error != io.EOF {
		iter.Error = newIter.Error
	}
}
