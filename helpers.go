package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
)

func Is(err, target error) bool { return stderrors.Is(err, target) }

func As(err error, target interface{}) bool { return stderrors.As(err, target) }

func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func Describe(err error) string {
	type withStack interface {
		GetStack() []uintptr
	}

	result := fmt.Sprintf("%v", err)

	e, ok := err.(withStack)
	if ok && e != nil {
		frames := runtime.CallersFrames(e.GetStack())
		for {
			frame, more := frames.Next()
			result += fmt.Sprintf("\n%s\n\t%s:%d", frame.Function, frame.File, frame.Line)

			if !more {
				break
			}
		}
	}

	return result
}
