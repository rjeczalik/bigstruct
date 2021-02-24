package codec

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/rjeczalik/bigstruct/big"
)

var DefaultField = Default.
	RegisterMap("field", 50, make(Map))

type Field struct {
	Type    string
	Convert func(interface{}) (interface{}, error)
}

var _ big.Codec = (*Field)(nil)

func (f Field) Encode(_ context.Context, key string, o big.Struct) error {
	return f.convert("encode", key, o)
}

func (f Field) Decode(_ context.Context, key string, o big.Struct) error {
	return f.convert("decode", key, o)
}

func (f Field) convert(op, key string, o big.Struct) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	if len(n.Children) != 0 && n.Value == nil {
		return &big.Error{
			Type: f.Type,
			Op:   op,
			Key:  key,
			Err:  errors.New("unable to convert value in non-leaf node"),
		}
	}

	w, err := f.Convert(n.Value)
	if err != nil {
		return &big.Error{
			Type: f.Type,
			Op:   op,
			Key:  key,
			Err:  err,
		}
	}

	n.Value = w
	if n.Type == "" {
		n.Type = f.Type
	}
	o[k] = n

	return nil
}

func (f Field) GoString() string {
	return fmt.Sprintf("codec.Field{Type: %q}", f.Type)
}

func isNull(v interface{}) bool {
	return v == nil || v == big.NoValue
}
