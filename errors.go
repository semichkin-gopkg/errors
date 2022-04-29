package errors

type stack []uintptr

type Error struct {
	Message string
	Code    string
	parent  error
	stack   *stack
}

func New(message string) *Error {
	return &Error{
		Message: message,
		stack:   callers(),
	}
}

func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}

func (e *Error) WithCode(code string) *Error {
	e.Code = code
	return e
}

func (e *Error) WithStack() *Error {
	e.stack = callers()
	return e
}

func (e *Error) WithParent(err error) *Error {
	e.parent = err
	return e
}

func (e *Error) Error() string { return e.Message }

func (e *Error) Unwrap() error { return e.parent }

func (e *Error) Is(err error) bool {
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
