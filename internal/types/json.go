package types

import (
	"bytes"
	"encoding/json"
)

type JSON string

func MakeJSON(v interface{}) JSON {
	if v == nil {
		return ""
	}

	p, err := jsonMarshal(v, false)
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	return JSON(bytes.TrimSpace(p))
}

func MakePrettyJSON(v interface{}) JSON {
	if v == nil {
		return ""
	}

	p, err := jsonMarshal(v, true)
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	return JSON(bytes.TrimSpace(p))
}

func (s JSON) Value() interface{} {
	if s == "" {
		return nil
	}

	var v interface{}

	if err := s.Unmarshal(&v); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return v
}

func (s JSON) String() string {
	return string(s)
}

func (s JSON) Bytes() []byte {
	if s == "" {
		return nil
	}
	return []byte(s)
}

func (s JSON) Unmarshal(v interface{}) error {
	if s == "" {
		return nil
	}
	return json.Unmarshal(s.Bytes(), v)
}

func jsonMarshal(v interface{}, pretty bool) ([]byte, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if pretty {
		enc.SetIndent("", "\t")
	}

	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
