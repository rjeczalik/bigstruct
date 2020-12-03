package cti

import (
	"errors"
)

type Error struct {
	Encoding string
	Op       string
	Key      string
	Err      error
	Next     *Error
}

var _ error = (*Error)(nil)

func (e *Error) Error() string {
	return e.Encoding + `: failed to ` + e.Op + `"` + e.Key + `": ` + e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
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
	switch v := err.(type) {
	case nil:
		// skip
	case *Error:
		e.Next = v
	default:
		e.Next = &Error{
			Encoding: e.Encoding,
			Op:       e.Op,
			Key:      e.Key,
			Err:      v,
			Next:     nil,
		}
	}

	return e
}
