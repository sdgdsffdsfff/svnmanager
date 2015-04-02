package helper

import (
	"errors"
)

type ErrorType int

const (
	NoError ErrorType = iota
	DefaultError
	ExistsError
	ParamsError
	RequestError
	DbError
	EmptyError
	OfflineError
	BusyError
)

func NewError(msg string, err ...error) error {
	str := msg
	if len(err) > 0 && err[0] != nil {
		str += ": " + err[0].Error()
	}
	return errors.New(str)
}
