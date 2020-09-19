package cti

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

type encoding struct {
	*Encoding
	name  string
	id    Attr
	match func(string) bool
}

type Serializer struct {
	encodings []encoding
	err       []error
}

var defaultSerializer = NewSerializer()

func NewSerializer() *Serializer {
	var (
		s = &Serializer{}
	)

	s.mustRegisterEncoding(
		"json", AttrJSON, suffix(".json"),
		&Encoding{
			Marshal:   json.Marshal,
			Unmarshal: json.Unmarshal,
		},
	)

	s.mustRegisterEncoding(
		"ini", AttrINI, suffix(".ini"),
		&Encoding{
			Marshal:   ini.Marshal,
			Unmarshal: ini.Unmarshal,
		},
	)

	s.mustRegisterEncoding(
		"flag", AttrFlag, suffix(),
		&Encoding{
			Marshal:   flag.Marshal,
			Unmarshal: flag.Unmarshal,
		},
	)

	s.mustRegisterEncoding(
		"toml", AttrTOML, suffix(".toml"),
		&Encoding{
			Marshal:   toml.Marshal,
			Unmarshal: toml.Unmarshal,
		},
	)

	s.mustRegisterEncoding(
		"yaml", AttrYAML, suffix(".yaml", ".yml"),
		&Encoding{
			Marshal:   yaml.Marshal,
			Unmarshal: yaml.Unmarshal,
		},
	)

	return s
}

func (s *Serializer) RegisterEncoding(name string, id Attr, match func(string) bool, e *Encoding) error {
	for _, e := range s.encodings {
		if e.name == name {
			return fmt.Errorf("%q encoding already registered", name)
		}
		if e.id&id != 0 {
			return fmt.Errorf("%q encoding (%b) conflicts with %q (%b)", name, id, e.name, e.id)
		}
	}

	s.encodings = append(s.encodings, encoding{
		Encoding: e,
		name:     name,
		match:    match,
		id:       id,
	})

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
			enc *encoding
		)

		if enc = s.lookup(n.Attr); enc != nil {
			if *err = enc.Decode(key, o); *err != nil {
				return false
			}
		}

		if len(n.Children) != 0 {
			return true // not a leaf node, ignore
		}

		if enc = s.guess(k); enc == nil || enc.Decode(key, o) != nil {
			var ok bool

			for i := range s.encodings {
				enc = &s.encodings[i]

				if enc.Decode(key, o) == nil {
					ok = true
					break
				}
			}

			if !ok {
				return true // not decoded, ignore but don't stop
			}
		}

		n = o[k]
		n.Attr |= enc.id
		o[k] = n

		return true
	}
}

func (s *Serializer) compact(err *error) Func {
	return func(key string, o Object) bool {
		var (
			k   = path.Base(key)
			n   = o[k]
			enc *encoding
		)

		if enc = s.lookup(n.Attr); enc == nil {
			return true
		}

		if *err = enc.Encode(key, o); *err != nil {
			return false
		}

		return true
	}
}

func (s *Serializer) mustRegisterEncoding(name string, id Attr, match func(string) bool, e *Encoding) {
	if err := s.RegisterEncoding(name, id, match, e); err != nil {
		panic("unexpected error registering " + name + ": " + err.Error())
	}
}

func (s *Serializer) lookup(attr Attr) *encoding {
	for i := range s.encodings {
		enc := &s.encodings[i]

		if attr.Has(enc.id) {
			return enc
		}
	}

	return nil
}

func (s *Serializer) guess(key string) *encoding {
	for i := range s.encodings {
		enc := &s.encodings[i]

		if enc.match(key) {
			return enc
		}
	}

	return nil
}

func suffix(s ...string) func(string) bool {
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
