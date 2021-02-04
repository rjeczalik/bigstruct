package codec

import (
	"encoding/json"
	"errors"
	"path"

	"github.com/rjeczalik/bigstruct/isr"
)

type Confetti struct{}

var _ isr.Codec = (*Confetti)(nil)

func (cc Confetti) Encode(key string, o isr.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	if len(n.Children) == 0 {
		return &isr.Error{
			Type: "isr",
			Op:   "encode",
			Key:  key,
			Err:  errors.New("nothing to encode"),
		}
	}

	if n.Value != nil {
		return &isr.Error{
			Type: "isr",
			Op:   "encode",
			Key:  key,
			Err:  errors.New("value in a non-leaf node"),
		}
	}

	// fixme(rjeczalik): remove the json.Marshal hack and rewrite the tree instead
	p, err := json.Marshal(n.Children)
	if err != nil {
		return &isr.Error{
			Type: "isr",
			Op:   "encode",
			Key:  key,
			Err:  err,
		}
	}

	var v map[string]interface{}

	if err := json.Unmarshal(p, &v); err != nil {
		return &isr.Error{
			Type: "isr",
			Op:   "encode",
			Key:  key,
			Err:  err,
		}
	}

	n.Children = isr.Make(v)
	o[k] = n

	return nil
}

func (cc Confetti) Decode(key string, o isr.Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := json.Marshal(n.Children.Value())
	if err != nil {
		return &isr.Error{
			Type: "isr",
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	var obj isr.Object

	if err := json.Unmarshal(p, &obj); err != nil {
		return &isr.Error{
			Type: "isr",
			Op:   "decode",
			Key:  key,
			Err:  err,
		}
	}

	n.Children = obj
	o[k] = n

	return nil
}
