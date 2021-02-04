package types

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

type YAML string

func MakeYAML(v interface{}) YAML {
	if v == nil {
		return ""
	}

	p, err := yaml.Marshal(v)
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	return YAML(bytes.TrimSpace(p))
}

func (s YAML) Value() interface{} {
	if s == "" {
		return nil
	}

	var v interface{}

	if err := s.Unmarshal(&v); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return v
}

func (s YAML) String() string {
	return string(s)
}

func (s YAML) Bytes() []byte {
	if s == "" {
		return nil
	}
	return []byte(s)
}

func (s YAML) Unmarshal(v interface{}) error {
	if s == "" {
		return nil
	}

	return yaml.Unmarshal(s.Bytes(), v)
}
