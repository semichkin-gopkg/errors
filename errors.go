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
}

func From(err interface{}) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e
	}

	return &Error{
		Message: fmt.Sprintf("%v", err),
		stack:   callers(),
	}
}

func New(message string) *Error {
	return From(message)
}

func (e *Error) WithMessage(message interface{}) *Error {
	if e == nil {
		return nil
	}

	e.Message = fmt.Sprintf("%v", message)
	return e
}

func (e *Error) WithCode(code string) *Error {
	if e == nil {
		return nil
	}

	e.Code = code
	return e
}

func (e *Error) WithStack() *Error {
	if e == nil {
		return nil
	}

	e.stack = callers()
	return e
}

func (e *Error) WithParent(err interface{}) *Error {
	if e == nil {
		return nil
	}

	e.Parent = From(err)
	return e
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

	if e.Code == "" {
		return false
	}

	appErr, ok := err.(*Error)
	if ok && appErr.Code == e.Code {
		return true
	}

	return false
}

func (e *Error) GetStack() []uintptr {
	return *e.stack
}

// wrapping logic

func (e *Error) WrapWithMessage(message string) *Error {
	err := clone(e)

	if err == nil {
		return nil
	}

	return err.
		WithMessage(message).
		WithParent(e)
}

func (e *Error) WrapWithCode(code string) *Error {
	err := clone(e)

	if err == nil {
		return nil
	}

	return err.
		WithCode(code).
		WithParent(e)
}

func (e *Error) WrapWithPrefix(prefix string) *Error {
	message := e.Message
	if prefix != "" {
		message = prefix + ": " + message
	}

	return e.WrapWithMessage(message)
}

func clone(err *Error) *Error {
	if err == nil {
		return nil
	}

	newErr := *err
	return &newErr
}
