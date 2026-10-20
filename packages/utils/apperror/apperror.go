package apperror

import "fmt"

type Error struct {
	Code    Code
	Message string
	Err     error
}

func New(code Code, msg string) *Error {
	return Wrap(nil, code, msg)
}

func Wrap(err error, code Code, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
