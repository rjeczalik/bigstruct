package model_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

var ignoreSpace = cmp.Comparer(func(lhs []byte, rhs []byte) bool {
	return bytes.Equal(bytes.TrimSpace(lhs), bytes.TrimSpace(rhs))
})

type codec struct {
	marshal   func(interface{}) ([]byte, error)
	unmarshal func([]byte, interface{}) error
}

func TestObject(t *testing.T) {
	codecs := []codec{
		0: {json.Marshal, json.Unmarshal},
		1: {yaml.Marshal, yaml.Unmarshal},
	}

	cases := []struct {
		p     []byte
		want  model.Object
		mwant map[string]interface{}
		codec codec
	}{{
		[]byte(`{"foo":"bar"}`),
		`{"foo":"bar"}`,
		map[string]interface{}{"foo": "bar"},
		codecs[0],
	}, {
		[]byte(`foo: bar`),
		`{"foo":"bar"}`,
		map[string]interface{}{"foo": "bar"},
		codecs[1],
	}}

	for _, cas := range cases {
		t.Run("", func(t *testing.T) {
			var obj model.Object

			if err := cas.codec.unmarshal(cas.p, &obj); err != nil {
				t.Fatalf("unmarshal()=%+v", err)
			}

			if !cmp.Equal(cas.want, obj) {
				t.Errorf("want != got:\n%s", cmp.Diff(cas.want, obj))
			}

			if got := obj.Map(); !cmp.Equal(cas.mwant, got) {
				t.Errorf("want != got:\n%s", cmp.Diff(cas.mwant, got))
			}

			p, err := cas.codec.marshal(obj)
			if err != nil {
				t.Errorf("marshal()=%+v", err)
			}

			if !cmp.Equal(cas.p, p, ignoreSpace) {
				t.Errorf("want != got:\n%s", cmp.Diff(string(cas.p), string(p)))
			}
		})
	}
}
