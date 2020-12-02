package cti

import (
	"errors"
	"path"

	"github.com/glaucusio/confetti/internal/objects"
)

type ObjectEncoder struct {
	Name      string
	Ext       []string
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
}

var _ Encoder = (*ObjectEncoder)(nil)

func (oe ObjectEncoder) String() string { return oe.Name }

func (oe ObjectEncoder) FileExt() []string { return oe.Ext }

func (oe ObjectEncoder) Encode(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := oe.Marshal(n.Children.Value())
	if err != nil {
		return &EncodingError{
			Encoding: oe.String(),
			Key:      key,
			Err:      err,
		}
	}

	n.Value = string(p)
	n.Children = nil
	o[k] = n

	return nil
}

func (oe ObjectEncoder) Decode(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := tobytes(n.Value)
	if err != nil {
		return &EncodingError{
			Encoding: oe.String(),
			Key:      key,
			Err:      err,
		}
	}

	var v interface{}

	if err := oe.Unmarshal(p, &v); err != nil {
		return &EncodingError{
			Encoding: oe.String(),
			Key:      key,
			Err:      err,
		}
	}

	obj := objects.Object(v)
	if len(obj) == 0 {
		return &EncodingError{
			Encoding: oe.String(),
			Key:      key,
			Err:      errors.New("not a struct or non-empty map"),
		}
	}

	n.Value = nil
	n.Children = Make(obj)

	o[k] = n

	return nil
}
