package errors

import (
	"fmt"
)

type stack []uintptr

type Error struct {
	Message string
	Code    string
	Parent  *Error
	stack   *stack

	original error
}

func From(err interface{}) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e
	}

	if e, ok := err.(error); ok {
		return &Error{
			Message:  fmt.Sprintf("%v", err),
			stack:    callers(),
			original: e,
		}
	}

	return &Error{
		Message: fmt.Sprintf("%v", err),
		stack:   callers(),
	}
}

func New(message string) *Error {
	return From(message)
}

func (e *Error) SetParent(parent interface{}) {
	if e != nil {
		e.Parent = From(parent)
	}
}

func (e *Error) WithParent(parent interface{}) *Error {
	if err := clone(e); err != nil {
		err.SetParent(parent)
		return err
	}

	return nil
}

func (e *Error) SetMessage(message interface{}) {
	if e != nil {
		e.Message = fmt.Sprintf("%v", message)
	}
}

func (e *Error) WithMessage(message interface{}) *Error {
	if err := clone(e); err != nil {
		err.SetMessage(message)
		return err.WithParent(e)
	}

	return nil
}

func (e *Error) SetCode(code interface{}) {
	if e != nil {
		e.Code = fmt.Sprintf("%v", code)
	}
}

func (e *Error) WithCode(code interface{}) *Error {
	if err := clone(e); err != nil {
		err.SetCode(code)
		return err.WithParent(e)
	}

	return nil
}

func (e *Error) ResetStack() {
	if e != nil {
		e.stack = callers()
	}
}

func (e *Error) WithStack() *Error {
	if err := clone(e); err != nil {
		err.ResetStack()
		return err.WithParent(e)
	}

	return nil
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	if e == nil || e.Parent == nil {
		return nil
	}

	return e.Parent
}

func (e *Error) Is(err error) bool {
	if e == nil {
		return false
	}

	if e.original != nil && e.original == err {
		return true
	}

	return e.Code != "" && e.Code == From(err).Code
}

func (e *Error) GetStack() []uintptr {
	return *e.stack
}

func clone(err *Error) *Error {
	if err == nil {
		return nil
	}

	newErr := *err
	return &newErr
}
