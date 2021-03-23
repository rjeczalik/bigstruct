package codec

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/bigstruct/big"
)

var Default = make(Map)

type Recursive struct {
	Codec big.Codec
}

var _ big.Codec = Recursive{}

func (r Recursive) Encode(ctx context.Context, key string, o big.Struct) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	err := n.Children.ReverseWalk(func(key string, o big.Struct) error {
		return r.Codec.Encode(ctx, key, o)
	})
	if err != nil {
		return fmt.Errorf("recursive error: %w", err)
	}

	if err := r.Codec.Encode(ctx, key, o); err != nil {
		return fmt.Errorf("recursive error: %w", err)
	}

	return nil
}

func (r Recursive) Decode(ctx context.Context, key string, o big.Struct) error {
	if err := r.Codec.Decode(ctx, key, o); err != nil {
		return fmt.Errorf("recursive error: %w", err)
	}

	var (
		k = path.Base(key)
		n = o[k]
	)

	err := n.Children.Walk(func(key string, o big.Struct) error {
		return r.Codec.Encode(ctx, key, o)
	})
	if err != nil {
		return fmt.Errorf("recursive error: %w", err)
	}

	return nil
}

func (r Recursive) GoString() string {
	if s, ok := r.Codec.(fmt.GoStringer); ok {
		return "codec.Recursive{" + s.GoString() + "}"
	}
	return "codec.Recursive{...}"
}

type contextKey struct{ string }

func cleanpath(s string) string {
	s = strings.TrimLeft(s, `/.\`)
	s = filepath.FromSlash(s)
	s = filepath.Clean(s)
	s = filepath.ToSlash(s)
	s = path.Join("/", s)

	return s
}

func tobytes(v interface{}) (p []byte, err error) {
	switch v := v.(type) {
	case nil:
		// ignore
	case []byte:
		p = v
	case string:
		p = []byte(v)
	case encoding.TextMarshaler:
		p, err = v.MarshalText()
	case encoding.BinaryMarshaler:
		p, err = v.MarshalBinary()
	default:
		err = fmt.Errorf("value is neither string nor []byte: %T", v)
	}

	switch {
	case err != nil:
		return nil, err
	case len(p) == 0:
		return nil, errors.New("byte slice is empty")
	}

	return p, nil
}

func nonil(v ...interface{}) interface{} {
	for _, v := range v {
		if v != nil {
			return v
		}
	}
	return nil
}

func reencode(in, out interface{}) error {
	p, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(p, out)
}
