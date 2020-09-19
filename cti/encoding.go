package cti

import (
	"errors"
	path "path"

	"github.com/glaucusio/confetti/internal/objects"
)

type EncodingError struct {
	Key string
	Err error
}

var _ error = (*EncodingError)(nil)

func (ee *EncodingError) Error() string {
	return `failed to encode "` + ee.Key + `": ` + ee.Err.Error()
}

func (ee *EncodingError) Unwrap() error {
	return ee.Err
}

type Encoding struct {
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
}

func (e *Encoding) Encode(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := e.Marshal(n.Children.Value())
	if err != nil {
		return &EncodingError{
			Key: key,
			Err: err,
		}
	}

	n.Value = p
	n.Children = nil
	o[k] = n

	return nil
}

func (e *Encoding) Decode(key string, o Object) error {
	var (
		p []byte
		k = path.Base(key)
		n = o[k]
	)

	switch v := n.Value.(type) {
	case []byte:
		p = v
	case string:
		p = []byte(v)
	case nil:
		if len(n.Children) == 0 {
			return &EncodingError{
				Key: key,
				Err: errors.New("nil value in leaf node"),
			}
		}
		return nil
	default:
		var err error

		if p, err = e.Marshal(v); err != nil {
			return &EncodingError{
				Key: key,
				Err: err,
			}
		}
	}

	var v interface{}

	if err := e.Unmarshal(p, &v); err != nil {
		return &EncodingError{
			Key: key,
			Err: err,
		}
	}

	obj := objects.Object(v)
	if len(obj) == 0 {
		return &EncodingError{
			Key: key,
			Err: errors.New("not a struct or non-empty map"),
		}
	}

	n.Value = nil
	n.Children = Make(obj)
	o[k] = n

	return nil
}
