package cti

import (
	"encoding/json"
	"path"
)

type ConfettiEncoder struct{}

var _ Encoder = (*ConfettiEncoder)(nil)

func (ce ConfettiEncoder) String() string { return "cti" }

func (ce ConfettiEncoder) FileExt() []string { return []string{"cti"} }

func (ce ConfettiEncoder) Encode(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	// fixme(rjeczalik): remove the json.Marshal hack and rewrite the tree instead
	p, err := json.Marshal(n.Children)
	if err != nil {
		return &Error{
			Encoding: ce.String(),
			Op:       "encode",
			Key:      key,
			Err:      err,
		}
	}

	var v map[string]interface{}

	if err := json.Unmarshal(p, &v); err != nil {
		return &Error{
			Encoding: ce.String(),
			Op:       "encode",
			Key:      key,
			Err:      err,
		}
	}

	n.Children = Make(v)
	o[k] = n

	return nil
}

func (ce ConfettiEncoder) Decode(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := json.Marshal(n.Children.Value())
	if err != nil {
		return &Error{
			Encoding: ce.String(),
			Op:       "decode",
			Key:      key,
			Err:      err,
		}
	}

	var obj Object

	if err := json.Unmarshal(p, &obj); err != nil {
		return &Error{
			Encoding: ce.String(),
			Op:       "decode",
			Key:      key,
			Err:      err,
		}
	}

	n.Children = obj
	o[k] = n

	return nil
}
