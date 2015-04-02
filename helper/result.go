package helper

import (
	"king/utils"
	"king/utils/JSON"
	"reflect"
)

func Error(args ...interface{}) JSON.Type {
	res := make(JSON.Type)
	res["code"] = "error"

	length := len(args)

	Map(args, func(key, msg interface{}) bool {
		index := key.(int)
		if msg == nil {
			return false
		}
		t := reflect.Indirect(reflect.ValueOf(msg)).Interface()
		kind := reflect.TypeOf(t).Kind()

		if index == 0 && reflect.TypeOf(msg).Name() == "ErrorType" {
			res["type"] = msg
			return false
		}

		switch kind {
		case reflect.Slice, reflect.Map:
			if index == 1 || length == 1 {
				res["result"] = msg
			}
		case reflect.String:
			if index == 0 {
				res["message"] = msg
			}
		case reflect.Ptr, reflect.Struct:
			if index == 1 || length == 1 {
				if d := utils.CallMethod(msg, "Error").(string); d != "" {
					res["message"] = d
				} else {
					res["result"] = msg
				}
			} else {
				res["result"] = msg
			}
		}
		return false
	})
	return res
}

func Success(data ...interface{}) JSON.Type {
	res := make(JSON.Type)
	res["code"] = "success"
	if len(data) > 0 {
		switch reflect.TypeOf(data[0]).Kind() {
		case reflect.String:
			res["message"] = data[0].(string)
		default:
			res["result"] = data[0]
			if len(data) > 1 {
				res = JSON.Extend(res, data[1:]...)
			}
		}
	}
	return res
}
