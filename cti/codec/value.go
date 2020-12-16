package codec

import (
	"errors"
	"fmt"
	"path"

	"github.com/glaucusio/confetti/cti"
)

var DefaultValue = Default.RegisterMap("value", make(Map))

type Value struct {
	Type    string
	Convert func(interface{}) (interface{}, error)
}

var _ cti.Codec = (*Value)(nil)

func (v Value) Encode(key string, o cti.Object) error {
	return v.convert("encode", key, o)
}

func (v Value) Decode(key string, o cti.Object) error {
	return v.convert("decode", key, o)
}

func (v Value) convert(op, key string, o cti.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	if len(n.Children) != 0 && n.Value == nil {
		return &cti.Error{
			Encoding: v.Type,
			Op:       op,
			Key:      key,
			Err:      errors.New("unable to convert value in non-leaf node"),
		}
	}

	w, err := v.Convert(n.Value)
	if err != nil {
		return &cti.Error{
			Encoding: v.Type,
			Op:       op,
			Key:      key,
			Err:      err,
		}
	}

	n.Value = w
	if n.Encoding == "" {
		n.Encoding = v.Type
	}
	o[k] = n

	return nil
}

func (v Value) GoString() string {
	return fmt.Sprintf("codec.Value{Type: %q}", v.Type)
}
