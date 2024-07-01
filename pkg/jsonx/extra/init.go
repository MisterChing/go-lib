package extra

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1

func RegisterCustomFuzzyDecoders() {
	jsoniter.RegisterTypeDecoder("string", &fuzzyStringDecoder{})
	jsoniter.RegisterTypeDecoder("bool", &fuzzyBoolDecoder{})
	jsoniter.RegisterTypeDecoder("int", &fuzzyIntegerDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(maxInt) || val < float64(minInt) {
				iter.ReportError("fuzzy decode int", "exceed range")
				return
			}
			*((*int)(ptr)) = int(val)
		} else {
			*((*int)(ptr)) = iter.ReadInt()
		}
	}})
}
