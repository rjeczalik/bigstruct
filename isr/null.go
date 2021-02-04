package isr

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"gopkg.in/yaml.v2"
)

var NoValue null

var (
	_ driver.Valuer    = null{}
	_ sql.Scanner      = null{}
	_ yaml.Unmarshaler = null{}
	_ yaml.Marshaler   = null{}
	_ json.Marshaler   = null{}
	_ json.Unmarshaler = null{}
)

type null struct{}

func (null) Scan(interface{}) error                      { return nil }
func (null) MarshalJSON() ([]byte, error)                { return []byte("null"), nil }
func (null) UnmarshalJSON(p []byte) error                { return nil }
func (null) MarshalYAML() (interface{}, error)           { return nil, nil }
func (null) UnmarshalYAML(func(interface{}) error) error { return nil }
func (null) Value() (driver.Value, error)                { return nil, nil }
