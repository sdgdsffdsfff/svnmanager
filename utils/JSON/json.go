package JSON

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"king/utils"
	"reflect"
)

type Type map[string]interface{}
type TypeSlice []map[string]interface{}

//序列化对象为字符串
func Stringify(data interface{}) string {
	if result := Compile(data); result != nil {
		return string(result)
	}
	return ""
}

//对象转换为byte切片
func Compile(data interface{}) []byte {
	str, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Compille Struct Error: ", err)
		return nil
	}
	return str
}

//转换为JSON.Type
func Parse(data interface{}) Type {
	var output Type

	real := reflect.Indirect(reflect.ValueOf(data))
	t := reflect.TypeOf(real)

	switch t.Kind() {
	case reflect.Struct, reflect.String, reflect.Slice:
		ParseToStruct(data, &output)
	case reflect.Map:
		if t.Name() == "Type" {
			return data.(Type)
		}
		output := make(Type)
		x := reflect.ValueOf(data)
		for _, v := range x.MapKeys() {
			var key string
			switch v.Type().Kind() {
			case reflect.String:
				key = v.String()
			case reflect.Struct:
				if m := utils.CallMethod(v.Interface(), "String"); m != nil {
					key = m.(string)
				}
			}
			if key != "" {
				output[key] = x.MapIndex(v).Interface()
			}
		}
		return output
	}
	return output
}

func ParseToStruct(data interface{}, m interface{}) error {
	real := reflect.Indirect(reflect.ValueOf(data))
	t := reflect.TypeOf(real.Interface())
	switch t.Kind() {
	case reflect.String:
		return json.Unmarshal([]byte(data.(string)), m)
	case reflect.Map, reflect.Struct:
		return json.Unmarshal(Compile(data), m)
	case reflect.Slice:
		if real.Type().Kind() == reflect.Uint8 {
			return json.Unmarshal(data.([]byte), m)
		}
	}
	return nil
}

//字符串转换为map数组
func ParseStringToSlice(data string) []Type {
	return ParseByteToSlice([]byte(data))
}

//byte切片转换为map数组
func ParseByteToSlice(data []byte) []Type {
	var m []Type
	err := json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return m
}

func ParseBlob(data []uint8) Type {
	return Parse(string(data))
}

//扩展map, 相同键值对后面的覆盖前面的
func Extend(src interface{}, dst ...interface{}) Type {
	var output = make(Type)
	org := Parse(src)

	for k, _ := range org {
		output[k] = org[k]
	}
	for _, d := range dst {
		if dd := Parse(d); dd != nil {
			for k, _ := range dd {
				output[k] = dd[k]
			}
		}
	}
	return output
}

//io.Reader转换为map
func FormRequest(body io.Reader) Type {
	json, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}
	return Parse(string(json))
}

func ConvertMap(data Type) map[string]interface{} {
	return map[string]interface{}(data)
}

func ConvertSlice(data []Type) []map[string]interface{} {
	d := TypeSlice{}
	for _, item := range data {
		d = append(d, ConvertMap(item))
	}
	return d
}

func GetKeys(data Type, fn ...func(str string) string) []string {
	keys := []string{}
	callback := func(str string) string {
		return str
	}

	if len(fn) > 0 {
		callback = fn[0]
	}

	for key, _ := range data {
		keys = append(keys, callback(key))
	}

	return keys
}
