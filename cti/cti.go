package cti

import (
	"encoding/json"
	"path"
)

func Field(a Attr, aux, value interface{}) Func {
	return func(key string, o Object) bool {
		var (
			k = path.Base(key)
			n = o[k]
		)

		n.Attr = a
		n.Aux = aux
		n.Value = value
		o[k] = n

		return true
	}
}

func Value(v interface{}) Func {
	return Field(0, nil, v)
}

func reencode(in, out interface{}) error {
	p, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(p, out)
}

func nonempty(s ...string) string {
	for _, s := range s {
		if s != "" {
			return s
		}
	}
	return ""
}

func nonil(v ...interface{}) interface{} {
	for _, v := range v {
		if v != nil {
			return v
		}
	}
	return nil
}
