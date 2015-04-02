package helper

import (
	"crypto/rand"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func AppendString(args ...interface{}) string {
	argLen := len(args)
	if argLen == 0 {
		return ""
	}

	var format string
	i := 0
	for i < argLen {
		i++
		format += "%v"
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

func Int64(arg interface{}) int64 {
	return int64(Num(arg))
}

func UpperCaseFirstLetter(str string) string {
	return strings.ToUpper(string(str[0])) + str[1:]
}

func Itoa64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func GetCMDOutputWithComplete(output []byte, err error) error {
	if err != nil {
		return err
	}
	lines := strings.Split(Trim(string(output)), "\n")
	lastLine := lines[len(lines)-1]
	if lastLine != "complete" {
		return NewError(lastLine)
	}
	return nil
}
