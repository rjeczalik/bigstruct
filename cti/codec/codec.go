package codec

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/glaucusio/confetti/cti"
)

var Default = make(Map)

func init() {
	cti.DefaultCodec = Default
}

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
