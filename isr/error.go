package isr

import (
	"errors"
	"fmt"
	"io"

	interr "github.com/rjeczalik/bigstruct/internal/errors"
)

type Error struct {
	Type string
	Op   string
	Key  string
	Err  error
	Next *Error

	stack *interr.Stack
}

var _ error = (*Error)(nil)

func (e *Error) Error() string {
	return e.Type + `: failed to ` + e.Op + ` "` + e.Key + `": ` + e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') && e.stack != nil {
			fmt.Fprintf(s, "%v", e.Unwrap())
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *Error) As(target interface{}) bool {
	switch err := target.(type) {
	case *Error:
		*err = *e
		return true
	default:
		for ; e != nil; e = e.Next {
			if errors.As(e.Err, target) {
				return true
			}
		}
	}

	return false
}

func (e *Error) Chain(err error) *Error {
	// fixme: handle e == nil
	switch v := err.(type) {
	case nil:
		if e.stack == nil {
			e.stack = interr.Callers()
		}
		return e
	case *Error:
		e.Next = v
	default:
		e.Next = &Error{
			Type: e.Type,
			Op:   e.Op,
			Key:  e.Key,
			Err:  v,
			Next: nil,
		}
	}

	return e.Next
}
