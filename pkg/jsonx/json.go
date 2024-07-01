package jsonx

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func init() {
	// RegisterFuzzyDecoders decode input from PHP with tolerance.
	//  It will handle string/number auto conversation, and treat empty [] as empty struct.
	extra.RegisterFuzzyDecoders()
	////自定义扩展
	//extra2.RegisterCustomFuzzyDecoders()
}

var newJson = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	return newJson.Marshal(v)
}
func MarshalToString(v interface{}) (string, error) {
	return newJson.MarshalToString(v)
}
func Unmarshal(data []byte, v interface{}) error {
	return newJson.Unmarshal(data, &v)
}
func UnmarshalFromString(data string, v interface{}) error {
	return newJson.UnmarshalFromString(data, &v)
}
