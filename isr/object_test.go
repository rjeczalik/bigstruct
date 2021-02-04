package isr_test

import (
	"testing"

	"github.com/rjeczalik/bigstruct/isr"
	"github.com/rjeczalik/bigstruct/isr/codec"

	"github.com/google/go-cmp/cmp"
)

func TestObject(t *testing.T) {
	o := make(isr.Object)

	o.Put("/foo/bar", isr.Value("[\"qux\",\"baz\"]", "object", "json"))
	o.Put("/ascii/48", isr.Value(int('a')))
	o.Put("/ascii/49", isr.Value(int('b')))
	o.Put("/ascii/50", isr.Value(int('c')))
	o.Put("/yaml", isr.Value("json: '{\"ini\":\"k=\\\"v\\\"\\nkey=\\\"value\\\"\\n\"}'\n"))
	o.Put("/yaml/json/ini", isr.Value(nil, "object", "ini"))

	got := o.Fields()
	want := isr.Fields{{
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

	if err := o.Decode(codec.Default); err != nil {
		t.Fatalf("o.Decode(nil)=%+v", err)
	}

	egot := o.Fields()
	ewant := isr.Fields{{
		Key: "/ascii",
	}, {
		Key:   "/ascii/48",
		Type:  "field/number",
		Value: int('a'),
	}, {
		Key:   "/ascii/49",
		Type:  "field/number",
		Value: int('b'),
	}, {
		Key:   "/ascii/50",
		Type:  "field/number",
		Value: int('c'),
	}, {
		Key: "/foo",
	}, {
		Key:  "/foo/bar",
		Type: "object/json",
	}, {
		Key:   "/foo/bar/0",
		Type:  "field/string",
		Value: "qux",
	}, {
		Key:   "/foo/bar/1",
		Type:  "field/string",
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
		Type:  "field/string",
		Value: "v",
	}, {
		Key:   "/yaml/json/ini/key",
		Type:  "field/string",
		Value: "value",
	}}

	if !cmp.Equal(ewant, egot) {
		t.Fatalf("egot != ewant:\n%s", cmp.Diff(ewant, egot))
	}

	var rgot isr.Fields
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

	var igot isr.Fields
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

	var cgot isr.Fields
	want[6].Type = "object/yaml"

	if err := o.Encode(codec.Default); err != nil {
		t.Fatalf("o.Encode()=%+v", err)
	}

	o.Walk(cgot.Append)

	cwant := isr.Fields{{
		Key: "/ascii",
	}, {
		Key:   "/ascii/48",
		Type:  "field/number",
		Value: int('a'),
	}, {
		Key:   "/ascii/49",
		Type:  "field/number",
		Value: int('b'),
	}, {
		Key:   "/ascii/50",
		Type:  "field/number",
		Value: int('c'),
	}, {
		Key: "/foo",
	}, {
		Key:   "/foo/bar",
		Type:  "object/json",
		Value: "[\"qux\",\"baz\"]",
	}, {
		Key:  "/foo/bar/0",
		Type: "field/string",
	}, {
		Key:  "/foo/bar/1",
		Type: "field/string",
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
		Type: "field/string",
	}, {
		Key:  "/yaml/json/ini/key",
		Type: "field/string",
	}}

	if !cmp.Equal(cwant, cgot) {
		t.Fatalf("cgot != want:\n%s", cmp.Diff(cwant, cgot))
	}
}
