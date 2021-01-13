package codec

import (
	"errors"
	"fmt"
	"path"

	"github.com/glaucusio/confetti/cti"
	"github.com/glaucusio/confetti/internal/objects"
)

var DefaultObject = Default.
	RegisterMap("object", make(Map))

type Object struct {
	Type      string
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
}

var _ cti.Codec = (*Object)(nil)

func (obj Object) Encode(key string, o cti.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	if len(n.Children) == 0 {
		return &cti.Error{
			Type: obj.Type,
			Op:   "encode",
			Key:  key,
			Err:  errors.New("nothing to encode"),
		}
	}

	if n.Value != nil {
		return &cti.Error{
			Type: obj.Type,
			Op:   "encode",
			Key:  key,
			Err:  errors.New("value in a non-leaf node"),
		}
	}

	p, err := obj.Marshal(n.Children.Value())
	if err != nil {
		return &cti.Error{
			Type: obj.Type,
			Op:   "encode",
			Key:  key,
			Err:  err,
		}
	}

	n.Value = string(p)
	n.Children = n.Children.Schema()
	o[k] = n

	return nil
}

func (obj Object) Decode(key string, o cti.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := tobytes(n.Value)
	if err != nil {
		return &cti.Error{
			Type: obj.Type,
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	var v interface{}

	if err := obj.Unmarshal(p, &v); err != nil {
		return &cti.Error{
			Type: obj.Type,
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	ubj := objects.Object(v)
	if len(ubj) == 0 {
		return &cti.Error{
			Type: obj.Type,
			Op:   "decode",
			Key:  key,
			Err:  errors.New("not a struct or non-empty map"),
		}
	}

	n.Value = nil
	n.Children = n.Children.Schema().Merge(cti.Make(ubj))
	if n.Type == "" {
		n.Type = obj.Type
	}
	o[k] = n

	return nil
}

func (obj Object) GoString() string {
	return fmt.Sprintf("codec.Object{Type: %q}", obj.Type)
}
