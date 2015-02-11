package helper

import (
	"errors"
)

type ErrorType int

const (
	DefaultError ErrorType = iota
	ExistsError
	ParamsError
	RequestError
	DbError

)

func NewError(msg string, err ...error) error {
	str := msg
	if len(err) > 0 && err[0] != nil {
		str += ": "+err[0].Error()
	}
	return errors.New(str)
}
