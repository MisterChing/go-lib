package jsonx

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func init() {
	// RegisterFuzzyDecoders decode input from PHP with tolerance.
	//  It will handle string/number auto conversation, and treat empty [] as empty struct.
	extra.RegisterFuzzyDecoders()
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}
func MarshalToString(obj interface{}) (string, error) {
	return json.MarshalToString(obj)
}
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, &v)
}
func UnmarshalFromString(data string, v interface{}) error {
	return json.UnmarshalFromString(data, &v)
}
