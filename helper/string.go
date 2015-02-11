package helper

import (
	"strconv"
	"crypto/rand"
	"strings"
	"fmt"
	"reflect"
)

func RandString(length int) string {
	alphaNum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphaNum[b%byte(len(alphaNum))]
	}
	return string(bytes)
}

func Random(length int) int {
	alphaNum := "0123456789"
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphaNum[b%byte(len(alphaNum))]
	}
	n, _ := strconv.Atoi(string(bytes))
	return n
}

func AppendString(args ...interface{}) string{
	argLen := len(args)
	if argLen == 0 {
		return ""
	}

	var format string
	i := 0
	for i < argLen {
		i++
		format+= "%v"
	}
	str := fmt.Sprintf(format, args...)
	return str
}

func Trim(s string) string {
	s = strings.Trim(s, " ")
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\\s")
	s = strings.Trim(s, "\t")
	return s
}

//todo real value
func Num(arg interface{}) int {
	if reflect.TypeOf(arg).Kind() == reflect.String {
		if num, err := strconv.Atoi(arg.(string)); err == nil {
			return num
		}
	}
	return 0
}

func UpperCaseFirstLetter(str string) string {
	return strings.ToUpper(string(str[0])) + str[1:]
}
