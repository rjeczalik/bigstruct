package codec

import (
	"errors"
	"path"

	"github.com/rjeczalik/bigstruct/isr"
)

type Bytes struct {
	Type      string
	Marshal   func([]byte) ([]byte, error)
	Unmarshal func([]byte) ([]byte, error)
}

var _ isr.Codec = (*Bytes)(nil)

func (b Bytes) Encode(key string, o isr.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := tobytes(n.Value)
	if err != nil {
		return &isr.Error{
			Type: b.Type,
			Op:   "encode",
			Key:  key,
			Err:  err,
		}
	}

	q, err := b.Marshal(p)
	if err != nil {
		return &isr.Error{
			Type: b.Type,
			Op:   "encode",
			Key:  key,
			Err:  err,
		}
	}

	n.Value = q
	o[k] = n

	return nil
}

func (b Bytes) Decode(key string, o isr.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	if len(n.Children) != 0 {
		return &isr.Error{
			Type: b.Type,
			Op:   "encode",
			Key:  key,
			Err:  errors.New("not a leaf node"),
		}
	}

	p, err := tobytes(n.Value)
	if err != nil {
		return &isr.Error{
			Type: b.Type,
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	q, err := b.Unmarshal(p)
	if err != nil {
		return &isr.Error{
			Type: b.Type,
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	n.Value = q
	o[k] = n

	return nil
}
