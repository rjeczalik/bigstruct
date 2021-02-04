package codec

import (
	"errors"
	"fmt"
	"path"

	"github.com/rjeczalik/bigstruct/isr"
	"github.com/rjeczalik/bigstruct/internal/objects"
)

var DefaultObject = Default.
	RegisterMap("object", 100, make(Map))

type Object struct {
	Type      string
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
}

var _ isr.Codec = (*Object)(nil)

func (obj Object) Encode(key string, o isr.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	if len(n.Children) == 0 {
		return &isr.Error{
			Type: obj.Type,
			Op:   "encode",
			Key:  key,
			Err:  errors.New("nothing to encode"),
		}
	}

	if n.Value != nil {
		return &isr.Error{
			Type: obj.Type,
			Op:   "encode",
			Key:  key,
			Err:  errors.New("value in a non-leaf node"),
		}
	}

	p, err := obj.Marshal(n.Children.Value())
	if err != nil {
		return &isr.Error{
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

func (obj Object) Decode(key string, o isr.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := tobytes(n.Value)
	if err != nil {
		return &isr.Error{
			Type: obj.Type,
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	var v interface{}

	if err := obj.Unmarshal(p, &v); err != nil {
		return &isr.Error{
			Type: obj.Type,
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	ubj := objects.Object(v)
	if len(ubj) == 0 {
		return &isr.Error{
			Type: obj.Type,
			Op:   "decode",
			Key:  key,
			Err:  errors.New("not a struct or non-empty map"),
		}
	}

	n.Value = nil
	n.Children = n.Children.Schema().Merge(isr.Make(ubj))
	if n.Type == "" {
		n.Type = obj.Type
	}
	o[k] = n

	return nil
}

func (obj Object) GoString() string {
	return fmt.Sprintf("codec.Object{Type: %q}", obj.Type)
}
