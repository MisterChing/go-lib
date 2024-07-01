package jsonx

import (
	"encoding/json"
	"testing"
)

type testCase struct {
	name string
	in   string
	out  string
}

func TestString(t *testing.T) {
	testArr := []testCase{
		{"string_1", `{"name":"ching"}`, `{"name":"ching"}`},
		{"string_2", `{"name":123}`, `{"name":"123"}`},
		{"string_3", `{"name":0}`, `{"name":"0"}`},
		{"string_4", `{"name":-1}`, `{"name":"-1"}`},
		{"string_5", `{"name":11.11}`, `{"name":"11.11"}`},
		{"string_6", `{"name":[]}`, `{"name":""}`},
		{"string_7", `{"name":{}}`, `{"name":""}`},
		{"string_8", `{"name":null}`, `{"name":""}`},
	}
	type Out struct {
		Name string `json:"name"`
	}
	for _, v := range testArr {
		t.Run(v.name, func(t *testing.T) {
			var out Out
			err := json.Unmarshal([]byte(v.in), &out)
			if err != nil {
				t.Errorf("Unmarshal error:%+v", err)
			}
			s, err := json.Marshal(out)
			ss := string(s)
			if ss != v.out || err != nil {
				t.Errorf("result error = %v, origin:%+v my:%+v, want:%v", err, v.in, ss, v.out)
			}
			// jsoniter
			err = Unmarshal([]byte(v.in), &out)
			if err != nil {
				t.Errorf("jsoniter Unmarshal error:%+v", err)
			}
			s, err = Marshal(out)
			ss = string(s)
			if ss != v.out || err != nil {
				t.Errorf("jsoniter result error = %v, origin:%+v my:%+v, want:%v", err, v.in, ss, v.out)
			}
		})
	}
}

func TestInt(t *testing.T) {
	testArr := []testCase{
		{"int_1", `{"age":20}`, `{"age":20}`},
		{"int_2", `{"age":"20"}`, `{"age":20}`},
		{"int_3", `{"age":11.11}`, `{"age":11}`},
		{"int_4", `{"age":11.6}`, `{"age":11}`},
		{"int_5", `{"age":[]}`, `{"age":0}`},
		{"int_6", `{"age":{}}`, `{"age":0}`},
		{"int_7", `{"age":null}`, `{"age":0}`},
		{"int_8", `{"age":0.1}`, `{"age":0}`},
		{"int_9", `{"age":-0.1}`, `{"age":0}`},
		{"int_9", `{"age":""}`, `{"age":0}`},
	}
	type Out struct {
		Age int `json:"age"`
	}
	for _, v := range testArr {
		t.Run(v.name, func(t *testing.T) {
			var out Out
			err := json.Unmarshal([]byte(v.in), &out)
			if err != nil {
				t.Errorf("Unmarshal error:%+v", err)
			}
			s, err := json.Marshal(out)
			ss := string(s)
			if ss != v.out || err != nil {
				t.Errorf("result error = %v, origin:%+v my:%+v, want:%v", err, v.in, ss, v.out)
			}
			// jsoniter
			err = Unmarshal([]byte(v.in), &out)
			if err != nil {
				t.Errorf("jsoniter Unmarshal error:%+v", err)
			}
			s, err = Marshal(out)
			ss = string(s)
			if ss != v.out || err != nil {
				t.Errorf("jsoniter result error = %v, origin:%+v my:%+v, want:%v", err, v.in, ss, v.out)
			}
		})
	}
}

func TestBool(t *testing.T) {
	testArr := []testCase{
		{"bool_1", `{"is_new":true}`, `{"is_new":true}`},
		{"bool_2", `{"is_new":false}`, `{"is_new":false}`},
		{"bool_3", `{"is_new":"true"}`, `{"is_new":true}`},
		{"bool_4", `{"is_new":"false"}`, `{"is_new":false}`},
		{"bool_5", `{"is_new":"123"}`, `{"is_new":false}`},
		{"bool_6", `{"is_new":123}`, `{"is_new":false}`},
		{"bool_7", `{"is_new":""}`, `{"is_new":false}`},
		{"bool_8", `{"is_new":null}`, `{"is_new":false}`},
		{"bool_9", `{"is_new":[]}`, `{"is_new":false}`},
		{"bool_10", `{"is_new":{}}`, `{"is_new":false}`},
		{"bool_11", `{"is_new":[1]}`, `{"is_new":false}`},
		{"bool_12", `{"is_new":["ching"]}`, `{"is_new":false}`},
		{"bool_13", `{"is_new":"1"}`, `{"is_new":true}`},
		{"bool_14", `{"is_new":1}`, `{"is_new":true}`},
		{"bool_15", `{"is_new":0}`, `{"is_new":false}`},
		{"bool_16", `{"is_new":"0"}`, `{"is_new":false}`},
	}
	type Out struct {
		IsNew bool `json:"is_new"`
	}
	for _, v := range testArr {
		t.Run(v.name, func(t *testing.T) {
			var out Out
			err := json.Unmarshal([]byte(v.in), &out)
			if err != nil {
				t.Errorf("Unmarshal error:%+v", err)
			}
			s, err := json.Marshal(out)
			ss := string(s)
			if ss != v.out || err != nil {
				t.Errorf("result error = %v, origin:%+v my:%+v, want:%v", err, v.in, ss, v.out)
			}
			// jsoniter
			err = Unmarshal([]byte(v.in), &out)
			if err != nil {
				t.Errorf("jsoniter Unmarshal error:%+v", err)
			}
			s, err = Marshal(out)
			ss = string(s)
			if ss != v.out || err != nil {
				t.Errorf("jsoniter result error = %v, origin:%+v my:%+v, want:%v", err, v.in, ss, v.out)
			}
		})
	}
}
