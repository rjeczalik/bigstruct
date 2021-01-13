package cti_test

import (
	"testing"

	"github.com/glaucusio/confetti/cti"
	_ "github.com/glaucusio/confetti/cti/codec"

	"github.com/google/go-cmp/cmp"
)

func TestObject(t *testing.T) {
	o := make(cti.Object)

	o.Put("/foo/bar", cti.Value("[\"qux\",\"baz\"]", "object", "json"))
	o.Put("/ascii/48", cti.Value(int('a')))
	o.Put("/ascii/49", cti.Value(int('b')))
	o.Put("/ascii/50", cti.Value(int('c')))
	o.Put("/yaml", cti.Value("json: '{\"ini\":\"k=\\\"v\\\"\\nkey=\\\"value\\\"\\n\"}'\n"))
	o.Put("/yaml/json/ini", cti.Value(nil, "object", "ini"))

	got := o.Fields()
	want := cti.Fields{{
		Key: "/ascii",
	}, {
		Key:   "/ascii/48",
		Value: int('a'),
	}, {
		Key:   "/ascii/49",
		Value: int('b'),
	}, {
		Key:   "/ascii/50",
		Value: int('c'),
	}, {
		Key: "/foo",
	}, {
		Key:   "/foo/bar",
		Type:  "object/json",
		Value: "[\"qux\",\"baz\"]",
	}, {
		Key:   "/yaml",
		Value: "json: '{\"ini\":\"k=\\\"v\\\"\\nkey=\\\"value\\\"\\n\"}'\n",
	}, {
		Key: "/yaml/json",
	}, {
		Key:  "/yaml/json/ini",
		Type: "object/ini",
	}}

	if !cmp.Equal(want, got) {
		t.Fatalf("got != want:\n%s", cmp.Diff(want, got))
	}

	if err := o.Decode(nil); err != nil {
		t.Fatalf("o.Decode(nil)=%+v", err)
	}

	egot := o.Fields()
	ewant := cti.Fields{{
		Key: "/ascii",
	}, {
		Key:   "/ascii/48",
		Type:  "value/number",
		Value: int('a'),
	}, {
		Key:   "/ascii/49",
		Type:  "value/number",
		Value: int('b'),
	}, {
		Key:   "/ascii/50",
		Type:  "value/number",
		Value: int('c'),
	}, {
		Key: "/foo",
	}, {
		Key:  "/foo/bar",
		Type: "object/json",
	}, {
		Key:   "/foo/bar/0",
		Type:  "value/string",
		Value: "qux",
	}, {
		Key:   "/foo/bar/1",
		Type:  "value/string",
		Value: "baz",
	}, {
		Key:  "/yaml",
		Type: "object/yaml",
	}, {
		Key:  "/yaml/json",
		Type: "object/json",
	}, {
		Key:  "/yaml/json/ini",
		Type: "object/ini",
	}, {
		Key:   "/yaml/json/ini/k",
		Type:  "value/string",
		Value: "v",
	}, {
		Key:   "/yaml/json/ini/key",
		Type:  "value/string",
		Value: "value",
	}}

	if !cmp.Equal(ewant, egot) {
		t.Fatalf("egot != ewant:\n%s", cmp.Diff(ewant, egot))
	}

	var rgot cti.Fields
	rwant := []string{
		"/yaml/json/ini/key",
		"/yaml/json/ini/k",
		"/yaml/json/ini",
		"/yaml/json",
		"/yaml",
		"/foo/bar/1",
		"/foo/bar/0",
		"/foo/bar",
		"/foo",
		"/ascii/50",
		"/ascii/49",
		"/ascii/48",
		"/ascii",
	}

	o.ReverseWalk(rgot.Append)

	if !cmp.Equal(rwant, rgot.Keys()) {
		t.Fatalf("rgot != rwant:\n%s", cmp.Diff(rwant, rgot.Keys()))
	}

	var igot cti.Fields
	iwant := []string{
		"/ascii/48",
		"/ascii/49",
		"/ascii/50",
		"/foo/bar/0",
		"/foo/bar/1",
		"/yaml/json/ini/k",
		"/yaml/json/ini/key",
	}

	o.ForEach(igot.Append)

	if !cmp.Equal(iwant, igot.Keys()) {
		t.Fatalf("igot != iwant:\n%s", cmp.Diff(iwant, igot.Keys()))
	}

	var cgot cti.Fields
	want[6].Type = "object/yaml"

	if err := o.Encode(nil); err != nil {
		t.Fatalf("o.Encode()=%+v", err)
	}

	o.Walk(cgot.Append)

	cwant := cti.Fields{{
		Key: "/ascii",
	}, {
		Key:   "/ascii/48",
		Type:  "value/number",
		Value: int('a'),
	}, {
		Key:   "/ascii/49",
		Type:  "value/number",
		Value: int('b'),
	}, {
		Key:   "/ascii/50",
		Type:  "value/number",
		Value: int('c'),
	}, {
		Key: "/foo",
	}, {
		Key:   "/foo/bar",
		Type:  "object/json",
		Value: "[\"qux\",\"baz\"]",
	}, {
		Key:  "/foo/bar/0",
		Type: "value/string",
	}, {
		Key:  "/foo/bar/1",
		Type: "value/string",
	}, {
		Key:   "/yaml",
		Type:  "object/yaml",
		Value: "json: '{\"ini\":\"k=\\\"v\\\"\\nkey=\\\"value\\\"\\n\"}'\n",
	}, {
		Key:  "/yaml/json",
		Type: "object/json",
	}, {
		Key:  "/yaml/json/ini",
		Type: "object/ini",
	}, {
		Key:  "/yaml/json/ini/k",
		Type: "value/string",
	}, {
		Key:  "/yaml/json/ini/key",
		Type: "value/string",
	}}

	if !cmp.Equal(cwant, cgot) {
		t.Fatalf("cgot != want:\n%s", cmp.Diff(cwant, cgot))
	}
}
