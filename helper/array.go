package helper

import "reflect"

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
