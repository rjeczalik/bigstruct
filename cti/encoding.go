package cti

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"

	"github.com/glaucusio/confetti/internal/objects"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v1"
)

var DefaultEncodings = []Encoding{{
	ID:        AttrJSON,
	Name:      "json",
	Marshal:   json.Marshal,
	Unmarshal: json.Unmarshal,
	Match:     MatchSuffix(".json"),
}, {
	ID:        AttrINI,
	Name:      "ini",
	Marshal:   ini.Marshal,
	Unmarshal: ini.Unmarshal,
	Match:     MatchSuffix(".json"),
}, {
	ID:        AttrFlag,
	Name:      "flag",
	Marshal:   flag.Marshal,
	Unmarshal: flag.Unmarshal,
	Match:     MatchSuffix(),
}, {
	ID:        AttrTOML,
	Name:      "toml",
	Marshal:   toml.Marshal,
	Unmarshal: toml.Unmarshal,
	Match:     MatchSuffix(".toml"),
}, {
	ID:        AttrYAML,
	Name:      "yaml",
	Marshal:   yaml.Marshal,
	Unmarshal: yaml.Unmarshal,
	Match:     MatchSuffix(".yaml", ".yml"),
}}

type EncodingError struct {
	Key string
	Err error
}

var _ error = (*EncodingError)(nil)

func (ee *EncodingError) Error() string {
	return `failed to encode "` + ee.Key + `": ` + ee.Err.Error()
}

func (ee *EncodingError) Unwrap() error {
	return ee.Err
}

type Encoding struct {
	ID        Attr
	Name      string
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
	Match     func(key string) bool
}

func (e *Encoding) String() string {
	return fmt.Sprintf("%q (%b)", e.Name, e.ID)
}

func (e *Encoding) Encode(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	p, err := e.Marshal(n.Children.Value())
	if err != nil {
		return &EncodingError{
			Key: key,
			Err: err,
		}
	}

	if n.Attr.Has(AttrString) {
		n.Value = string(p)
		n.Attr &= ^AttrString
	} else {
		n.Value = p
	}
	n.Children = nil
	o[k] = n

	return nil
}

func (e *Encoding) Decode(key string, o Object) error {
	var (
		p []byte
		a Attr
		k = path.Base(key)
		n = o[k]
	)

	switch v := n.Value.(type) {
	case []byte:
		p = v
	case string:
		p = []byte(v)
		a = AttrString
	case nil:
		if len(n.Children) == 0 {
			return &EncodingError{
				Key: key,
				Err: errors.New("nil value in leaf node"),
			}
		}

		return nil
	default:
		var err error

		if p, err = e.Marshal(v); err != nil {
			return &EncodingError{
				Key: key,
				Err: err,
			}
		}
	}

	var v interface{}

	if err := e.Unmarshal(p, &v); err != nil {
		return &EncodingError{
			Key: key,
			Err: err,
		}
	}

	obj := objects.Object(v)
	if len(obj) == 0 {
		return &EncodingError{
			Key: key,
			Err: errors.New("not a struct or non-empty map"),
		}
	}

	n.Attr |= a
	n.Value = nil
	n.Children = Make(obj)
	o[k] = n

	return nil
}
