package cti

import (
	"fmt"
	"path"
	"strings"
)

type Serializer struct {
	enc []Encoding
	aux []Aux
}

var defaultSerializer = NewSerializer()

func NewSerializer() *Serializer {
	s := new(Serializer)

	for _, enc := range DefaultEncodings {
		s.mustRegisterEncoding(enc)
	}

	for _, aux := range DefaultAuxs {
		s.mustRegisterAux(aux)
	}

	return s
}

func (s *Serializer) RegisterEncoding(e Encoding) error {
	for _, enc := range s.enc {
		if enc.Name == e.Name {
			return fmt.Errorf("encoding %s already registered", &e)
		}
		if enc.ID&e.ID != 0 {
			return fmt.Errorf("encoding %s conflicts with %s", &e, &enc)
		}
	}

	s.enc = append(s.enc, e)

	return nil
}

func (s *Serializer) RegisterAux(x Aux) error {
	for _, aux := range s.aux {
		if aux.Name == x.Name {
			return fmt.Errorf("encoding %s already registered", &x)
		}
		if aux.ID&x.ID != 0 {
			return fmt.Errorf("encoding %s conflicts with %s", &x, &aux)
		}
	}

	s.aux = append(s.aux, x)

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
			k   = path.Base(key)
			n   = o[k]
			enc *Encoding
		)

		if x := s.lookupAux(n.Attr); x != nil {
			_ = x.Expand(key, o)
			n = o[k]
		}

		if len(n.Children) != 0 || n.Value == nil {
			return true // not a leaf node, ignore
		}

		if enc = s.lookupEnc(n.Attr); enc != nil {
			if *err = enc.Decode(key, o); *err != nil {
				return false
			}
		} else if enc = s.guessEnc(k); enc == nil || enc.Decode(key, o) != nil {
			var ok bool

			for i := range s.enc {
				enc = &s.enc[i]

				if err := enc.Decode(key, o); err == nil {
					ok = true
					break
				}
			}

			if !ok {
				return true // not decoded, ignore but don't stop
			}
		}

		n = o[k]
		n.Attr |= enc.ID
		o[k] = n

		return true
	}
}

func (s *Serializer) compact(err *error) Func {
	return func(key string, o Object) bool {
		var (
			k   = path.Base(key)
			n   = o[k]
			enc *Encoding
		)

		if enc = s.lookupEnc(n.Attr); enc == nil {
			return true
		}

		if *err = enc.Encode(key, o); *err != nil {
			return false
		}

		return true
	}
}

func (s *Serializer) mustRegisterEncoding(e Encoding) {
	if err := s.RegisterEncoding(e); err != nil {
		panic("unexpected error registering encoding " + e.String() + ": " + err.Error())
	}
}

func (s *Serializer) mustRegisterAux(x Aux) {
	if err := s.RegisterAux(x); err != nil {
		panic("unexpected error registering aux " + x.String() + ": " + err.Error())
	}
}

func (s *Serializer) lookupEnc(attr Attr) *Encoding {
	for i := range s.enc {
		enc := &s.enc[i]

		if attr.Has(enc.ID) {
			return enc
		}
	}

	return nil
}

func (s *Serializer) lookupAux(attr Attr) *Aux {
	for i := range s.aux {
		aux := &s.aux[i]

		if attr.Has(aux.ID) {
			return aux
		}
	}

	return nil
}

func (s *Serializer) guessEnc(key string) *Encoding {
	for i := range s.enc {
		enc := &s.enc[i]

		if enc.Match(key) {
			return enc
		}
	}

	return nil
}

func MatchSuffix(s ...string) func(string) bool {
	return func(key string) bool {
		key = strings.ToLower(key)

		for _, s := range s {
			if strings.HasSuffix(key, s) {
				return true
			}
		}

		return false
	}
}
