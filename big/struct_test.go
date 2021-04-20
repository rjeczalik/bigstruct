package big_test

import (
	"context"
	"testing"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/big/codec"

	"github.com/google/go-cmp/cmp"
)

func TestObject(t *testing.T) {
	o := make(big.Struct)

	o.Put("/foo/bar", big.Value("[\"qux\",\"baz\"]", "struct", "json"))
	o.Put("/ascii/48", big.Value(int('a')))
	o.Put("/ascii/49", big.Value(int('b')))
	o.Put("/ascii/50", big.Value(int('c')))
	o.Put("/yaml", big.Value("json: '{\"ini\":\"k=\\\"v\\\"\\nkey=\\\"value\\\"\\n\"}'\n"))
	o.Put("/yaml/json/ini", big.Value(nil, "struct", "ini"))

	got := o.Fields()
	want := big.Fields{{
		Key:   "/ascii/48",
		Value: int('a'),
	}, {
		Key:   "/ascii/49",
		Value: int('b'),
	}, {
		Key:   "/ascii/50",
		Value: int('c'),
	}, {
		Key:   "/foo/bar",
		Type:  "struct/json",
		Value: "[\"qux\",\"baz\"]",
	}, {
		Key:  "/yaml/json/ini",
		Type: "struct/ini",
	}, {
		Key:   "/yaml",
		Value: "json: '{\"ini\":\"k=\\\"v\\\"\\nkey=\\\"value\\\"\\n\"}'\n",
	}}

	if !cmp.Equal(want, got) {
		t.Fatalf("got != want:\n%s", cmp.Diff(want, got))
	}

	if err := o.Decode(context.TODO(), codec.Default); err != nil {
		t.Fatalf("o.Decode(nil)=%+v", err)
	}

	egot := o.Fields()
	ewant := big.Fields{{
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
		Key:   "/foo/bar/0",
		Type:  "field/string",
		Value: "qux",
	}, {
		Key:   "/foo/bar/1",
		Type:  "field/string",
		Value: "baz",
	}, {
		Key:   "/yaml/json/ini/k",
		Type:  "field/string",
		Value: "v",
	}, {
		Key:   "/yaml/json/ini/key",
		Type:  "field/string",
		Value: "value",
	}, {
		Key:  "/yaml/json/ini",
		Type: "struct/ini",
	}, {
		Key:  "/yaml/json",
		Type: "struct/json",
	}, {
		Key:  "/yaml",
		Type: "struct/yaml",
	}, {
		Key:  "/foo/bar",
		Type: "struct/json",
	}}

	if !cmp.Equal(ewant, egot) {
		t.Fatalf("egot != ewant:\n%s", cmp.Diff(ewant, egot))
	}

	var rgot big.Fields
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

	var igot big.Fields
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

	var cgot big.Fields

	if err := o.Encode(context.TODO(), codec.Default); err != nil {
		t.Fatalf("o.Encode()=%+v", err)
	}

	o.Walk(cgot.Append)

	cwant := big.Fields{{
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
		Type:  "struct/json",
		Value: "[\"qux\",\"baz\"]",
	}, {
		Key:  "/foo/bar/0",
		Type: "field/string",
	}, {
		Key:  "/foo/bar/1",
		Type: "field/string",
	}, {
		Key:   "/yaml",
		Type:  "struct/yaml",
		Value: "json: '{\"ini\":\"k=\\\"v\\\"\\nkey=\\\"value\\\"\\n\"}'\n",
	}, {
		Key:  "/yaml/json",
		Type: "struct/json",
	}, {
		Key:  "/yaml/json/ini",
		Type: "struct/ini",
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
