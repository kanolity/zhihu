package xcode

import "strconv"

type XCode interface {
	Error() string
	Code() int
	Message() string
	Details() []interface{}
}

type Code struct {
	code    int
	message string
}

func (c Code) Error() string {
	if len(c.message) > 0 {
		return c.message
	}
	return strconv.Itoa(c.code)
}

func (c Code) Code() int {
	return c.code
}

func (c Code) Message() string {
	return c.message
}
func (c Code) Details() []interface{} {
	return nil
}

func String(s string) Code {
	if len(s) == 0 {
		return OK
	}
	code, err := strconv.Atoi(s)
	if err != nil {
		return ServerErr
	}
	return Code{code: code}
}

func New(code int, message string) Code {
	return Code{code: code, message: message}
}

func add(code int, message string) Code {
	return Code{code: code, message: message}
}
