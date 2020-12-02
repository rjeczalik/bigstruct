package cti

import (
	"encoding/json"
	"strings"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

type Encoding []string

var (
	_ json.Marshaler   = (*Encoding)(nil)
	_ json.Unmarshaler = (*Encoding)(nil)
	_ yaml.Marshaler   = (*Encoding)(nil)
	_ yaml.Unmarshaler = (*Encoding)(nil)
)

func (e Encoding) Has(enc string) bool {
	for _, e := range e {
		if e == enc {
			return true
		}
	}
	return false
}

func (e Encoding) String() string {
	return strings.Join(e, "/")
}

func (e Encoding) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.String() + `"`), nil
}

func (e *Encoding) UnmarshalJSON(p []byte) error {
	var s string

	if err := json.Unmarshal(p, &s); err != nil {
		return err
	}

	*e = strings.Split(s, "/")

	return nil
}

func (e Encoding) MarshalYAML() (interface{}, error) {
	return strings.Join(e, "/"), nil
}

func (e *Encoding) UnmarshalYAML(n *yaml.Node) error {
	var s string

	if err := n.Decode(&s); err != nil {
		return err
	}

	*e = strings.Split(s, "/")

	return nil
}

func (e Encoding) Copy() Encoding {
	eCopy := make(Encoding, len(e))
	copy(eCopy, e)
	return eCopy
}

type Encoder interface {
	Encode(key string, o Object) error
	Decode(key string, o Object) error
	FileExt() []string
	String() string
}

var DefaultEncoders = []Encoder{
	ObjectEncoder{
		Name:      "json",
		Ext:       []string{"json"},
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	},
	ObjectEncoder{
		Name:      "ini",
		Ext:       []string{"conf"},
		Marshal:   ini.Marshal,
		Unmarshal: ini.Unmarshal,
	},
	ObjectEncoder{
		Name:      "flag",
		Marshal:   flag.Marshal,
		Unmarshal: flag.Unmarshal,
	},
	ObjectEncoder{
		Name:      "toml",
		Ext:       []string{"toml"},
		Marshal:   toml.Marshal,
		Unmarshal: toml.Unmarshal,
	},
	ObjectEncoder{
		Name:      "yaml",
		Ext:       []string{"yml", "yaml"},
		Marshal:   yaml.Marshal,
		Unmarshal: yaml.Unmarshal,
	},
	BytesEncoder{
		Name:      "gzip",
		Ext:       []string{"gz", "gzip"},
		Marshal:   gzip.Marshal,
		Unmarshal: gzip.Unmarshal,
	},
	TarEncoder{},
	ZipEncoder{},
}

type EncodingError struct {
	Encoding string
	Key      string
	Err      error
}

var _ error = (*EncodingError)(nil)

func (ee *EncodingError) Error() string {
	return ee.Encoding + `: failed to encode "` + ee.Key + `": ` + ee.Err.Error()
}

func (ee *EncodingError) Unwrap() error {
	return ee.Err
}
