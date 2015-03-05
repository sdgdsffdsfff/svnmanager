package helper

import (
	"reflect"
)

func Map(list interface{}, do func(int) bool) {
	t := reflect.TypeOf(list)
	kind := t.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		data := reflect.ValueOf(list)
		for count, length := 0, data.Len(); count < length; count++ {
			if do(count) == true {
				return
			}
		}
	}
}

func AsyncMap(list interface{}, do func(int) bool) {
	t := reflect.TypeOf(list)
	kind := t.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		data := reflect.ValueOf(list)
		length := data.Len()
		if length == 0 {
			return
		}
		end := make(chan int)
		count := 0
		for i := 0; i < length; i++ {
			go func(i int) {
				//break
				if do(i) == true {
					end <-1
					return
				}
				count++
				if count == length {
					end <-1
				}
			}(i)
		}
		<-end
	}
}

//TODO
//添加所有类型
//目前指支持int, string
func ExtendStruct(d, s interface{}, fields ...string ) error {

	dst := reflect.ValueOf(d).Elem()
	src := reflect.ValueOf(s).Elem()

	for _, key := range fields {
		srcField := src.FieldByName(key)
		dstField := dst.FieldByName(key)

		switch srcField.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dstField.SetInt(srcField.Int())
			break
		case reflect.String:
			dstField.SetString(srcField.String())
		}
	}
	return nil
}
