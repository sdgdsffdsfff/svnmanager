package helper

import (
	"reflect"
)

func Map(list interface{}, do func(key, value interface{}) bool) {
	t := reflect.TypeOf(list)
	v := reflect.ValueOf(list)
	kind := t.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		len := v.Len()
		if len == 0 {
			return
		}
		for i := 0; i < len; i++ {
			if do(i, v.Index(i).Interface()) == true {
				return
			}
		}
	} else if kind == reflect.Map {
		keys := v.MapKeys()
		len := len(keys)

		if len == 0 {
			return
		}
		for _, key := range keys {
			if do(key.Interface(), v.MapIndex(key).Interface()) == true {
				return
			}
		}
	}
}

//并发调用slice或者map
func AsyncMap(list interface{}, do func(key, value interface{}) bool) {
	t := reflect.TypeOf(list)
	v := reflect.ValueOf(list)
	kind := t.Kind()
	count := 0
	if kind == reflect.Array || kind == reflect.Slice {
		len := v.Len()
		if len == 0 {
			return
		}
		end := make(chan int)
		for i := 0; i < len; i++ {
			go func(i int) {
				if do(i, v.Index(i).Interface()) == true {
					end <- 1
					return
				}
				count++
				if count == len {
					end <- 1
				}
			}(i)
		}
		<-end
	} else if kind == reflect.Map {
		keys := v.MapKeys()
		len := len(keys)

		if len == 0 {
			return
		}
		end := make(chan int)
		for _, key := range keys {
			go func(key reflect.Value) {
				if do(key.Interface(), v.MapIndex(key).Interface()) == true {
					end <- 1
					return
				}
				count++
				if count == len {
					end <- 1
				}
			}(key)
		}
		<-end
	}
}

//TODO
//添加所有类型
//目前指支持int, string
func ExtendStruct(d, s interface{}, fields ...string) error {

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

func Cap(d interface{}) int {
	cap := 0
	t := reflect.TypeOf(d)
	v := reflect.ValueOf(d)
	k := t.Kind()

	if k == reflect.Slice || k == reflect.Array || k == reflect.Chan {
		cap = v.Len()
	} else if k == reflect.Map {
		cap = len(v.MapKeys())
	}
	return cap
}
