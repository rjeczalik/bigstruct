package codec

import (
	"errors"
	"fmt"
	"path"

	"github.com/rjeczalik/bigstruct/isr"
)

var DefaultField = Default.
	RegisterMap("field", 50, make(Map))

type Field struct {
	Type    string
	Convert func(interface{}) (interface{}, error)
}

var _ isr.Codec = (*Field)(nil)

func (v Field) Encode(key string, o isr.Object) error {
	return v.convert("encode", key, o)
}

func (v Field) Decode(key string, o isr.Object) error {
	return v.convert("decode", key, o)
}

func (v Field) convert(op, key string, o isr.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	if len(n.Children) != 0 && n.Value == nil {
		return &isr.Error{
			Type: v.Type,
			Op:   op,
			Key:  key,
			Err:  errors.New("unable to convert value in non-leaf node"),
		}
	}

	w, err := v.Convert(n.Value)
	if err != nil {
		return &isr.Error{
			Type: v.Type,
			Op:   op,
			Key:  key,
			Err:  err,
		}
	}

	n.Value = w
	if n.Type == "" {
		n.Type = v.Type
	}
	o[k] = n

	return nil
}

func (v Field) GoString() string {
	return fmt.Sprintf("codec.Field{Type: %q}", v.Type)
}

func isNull(v interface{}) bool {
	return v == nil || v == isr.NoValue
}
