package cti

import "path"

type BytesEncoder struct {
	Name      string
	Ext       []string
	Marshal   func([]byte) ([]byte, error)
	Unmarshal func([]byte) ([]byte, error)
}

var _ Encoder = (*BytesEncoder)(nil)

func (be BytesEncoder) String() string { return be.Name }

func (be BytesEncoder) FileExt() []string { return be.Ext }

func (be BytesEncoder) Encode(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := tobytes(n.Value)
	if err != nil {
		return &Error{
			Encoding: be.String(),
			Op:       "encode",
			Key:      key,
			Err:      err,
		}
	}

	q, err := be.Marshal(p)
	if err != nil {
		return &Error{
			Encoding: be.String(),
			Op:       "encode",
			Key:      key,
			Err:      err,
		}
	}

	n.Value = q
	o[k] = n

	return nil
}

func (be BytesEncoder) Decode(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := tobytes(n.Value)
	if err != nil {
		return &Error{
			Encoding: be.String(),
			Op:       "decode",
			Key:      key,
			Err:      err,
		}
	}

	q, err := be.Unmarshal(p)
	if err != nil {
		return &Error{
			Encoding: be.String(),
			Op:       "decode",
			Key:      key,
			Err:      err,
		}
	}

	n.Value = q
	o[k] = n

	return nil
}
