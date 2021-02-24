package codec

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/internal/objects"
)

var DefaultStruct = Default.
	RegisterMap("struct", 100, make(Map))

type Struct struct {
	Type      string
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
}

var _ big.Codec = (*Struct)(nil)

func (s Struct) Encode(_ context.Context, key string, o big.Struct) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	if len(n.Children) == 0 {
		return &big.Error{
			Type: s.Type,
			Op:   "encode",
			Key:  key,
			Err:  errors.New("nothing to encode"),
		}
	}

	if n.Value != nil {
		return &big.Error{
			Type: s.Type,
			Op:   "encode",
			Key:  key,
			Err:  errors.New("value in a non-leaf node"),
		}
	}

	p, err := s.Marshal(n.Children.Value())
	if err != nil {
		return &big.Error{
			Type: s.Type,
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

func (s Struct) Decode(_ context.Context, key string, o big.Struct) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := tobytes(n.Value)
	if err != nil {
		return &big.Error{
			Type: s.Type,
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	var v interface{}

	if err := s.Unmarshal(p, &v); err != nil {
		return &big.Error{
			Type: s.Type,
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	obj := objects.Object(v)
	if len(obj) == 0 {
		return &big.Error{
			Type: s.Type,
			Op:   "decode",
			Key:  key,
			Err:  errors.New("not a struct or non-empty map"),
		}
	}

	n.Value = nil
	n.Children = n.Children.Schema().Merge(big.Make(obj)) // todo: reduce allocs
	if n.Type == "" {
		n.Type = s.Type
	}
	o[k] = n

	return nil
}

func (s Struct) GoString() string {
	return fmt.Sprintf("codec.Struct{Type: %q}", s.Type)
}
