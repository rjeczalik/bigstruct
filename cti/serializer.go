package cti

import (
	"fmt"
	"path"
	"strings"
)

type Serializer struct {
	enc []Encoder
}

var defaultSerializer = NewSerializer()

func NewSerializer() *Serializer {
	s := new(Serializer)

	for _, enc := range DefaultEncoders {
		s.mustRegisterEncoding(enc)
	}

	return s
}

func (s *Serializer) RegisterEncoding(e Encoder) error {
	for _, enc := range s.enc {
		if enc.String() == e.String() {
			return fmt.Errorf("encoding %s already registered", e)
		}
	}

	s.enc = append(s.enc, e)

	return nil
}

func (s *Serializer) Expand(o Object) error {
	var err error

	o.Walk(s.expand(&err))

	return err
}

func (s *Serializer) Compact(o Object) error {
	var err error

	o.ReverseWalk(s.compact(&err))

	return err
}

func (s *Serializer) Marshal(o Object) ([]byte, error) {
	return nil, nil
}

func (s *Serializer) Unmarshal(o Object, v interface{}) error {
	return nil
}

func (s *Serializer) expand(err *error) Func {
	return func(key string, o Object) bool {
		var (
			k     = path.Base(key)
			n     = o[k]
			stack encstack
		)

		if len(n.Children) != 0 || n.Value == nil {
			return true // not a leaf node, ignore
		}

		if stack, *err = s.lookupEncstack(n.Encoding); *err == nil && len(stack) != 0 {
			if *err = stack.Decode(key, o); *err != nil {
				return false
			}
		} else if *err != nil {
			return false
		} else if stack = s.guessEncstack(k); stack == nil || stack.Decode(key, o) != nil {
			var ok bool

			for _, enc := range s.enc {
				if err := enc.Decode(key, o); err == nil {
					stack = encstack{enc}
					ok = true
					break
				}
			}

			if !ok {
				return true // not decoded, ignore but don't stop
			}
		}

		n = o[k]
		n.Encoding = stack.encoding()
		o[k] = n

		return true
	}
}

func (s *Serializer) compact(err *error) Func {
	return func(key string, o Object) bool {
		var (
			k     = path.Base(key)
			n     = o[k]
			stack encstack
		)

		if stack, *err = s.lookupEncstack(n.Encoding); len(stack) == 0 {
			return *err == nil
		}

		if *err = stack.Encode(key, o); *err != nil {
			return false
		}

		return true
	}
}

func (s *Serializer) mustRegisterEncoding(e Encoder) {
	if err := s.RegisterEncoding(e); err != nil {
		panic("unexpected error registering encoding " + e.String() + ": " + err.Error())
	}
}

func (s *Serializer) lookupEncstack(e Encoding) (stack encstack, err error) {
	for i := len(e) - 1; i >= 0; i-- {
		enc := s.lookupEnc(e[i])

		if enc == nil {
			return nil, fmt.Errorf("unsupported encoding: %q", e[i])
		}

		stack = append(stack, enc)
	}

	return stack, nil
}

func (s *Serializer) guessEncstack(key string) (stack encstack) {
	for _, ext := range filext(path.Base(key)) {
		if enc := s.lookupEncByExt(ext); enc != nil {
			stack = append(stack, enc)
		}
	}
	return stack
}

func (s *Serializer) lookupEnc(typ string) Encoder {
	for _, enc := range s.enc {
		if enc.String() == typ {
			return enc
		}
	}

	return nil
}

func (s *Serializer) lookupEncByExt(ext string) Encoder {
	for _, enc := range s.enc {
		for _, s := range enc.FileExt() {
			if s == ext {
				return enc
			}
		}
	}

	return nil
}

type encstack []Encoder

func (stack encstack) encoding() (e Encoding) {
	e = make([]string, 0, len(stack))

	for _, enc := range stack {
		e = append(e, enc.String())
	}

	return e
}

func (stack encstack) Decode(key string, o Object) error {
	for _, enc := range stack {
		if err := enc.Decode(key, o); err != nil {
			return err
		}
	}

	return nil
}

func (stack encstack) Encode(key string, o Object) error {
	for _, enc := range stack {
		if err := enc.Encode(key, o); err != nil {
			return err
		}
	}

	return nil
}

func filext(name string) (ext []string) {
	s := strings.Split(name, ".")
	if len(s) < 2 {
		return nil
	}

	for _, s := range s[1:] {
		switch s = strings.ToLower(s); s {
		case "tgz":
			ext = append(ext, "tar", "gz")
		default:
			ext = append(ext, s)
		}
	}

	return ext
}
