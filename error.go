package misc

import "fmt"

//----------------------------------------------------------------------------------------------------------------------------//

// Error --
type Error struct {
	code int
	msg  string
}

// SetCode --
func (me *Error) SetCode(code int) {
	me.code = code
}

// SetMessage --
func (me *Error) SetMessage(format string, options ...any) {
	me.msg = fmt.Sprintf(format, options...)
}

// Error --
func (me *Error) Error() string {
	return me.msg
}

// Code --
func (me *Error) Code() int {
	return me.code
}

// MakeError --
func MakeError(code int, format string, options ...any) *Error {
	e := &Error{}
	e.SetCode(code)
	e.SetMessage(format, options...)
	return e
}

//----------------------------------------------------------------------------------------------------------------------------//
