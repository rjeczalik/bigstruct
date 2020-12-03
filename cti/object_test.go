package cti_test

import (
	"testing"

	"github.com/glaucusio/confetti/cti"

	"github.com/google/go-cmp/cmp"
)

func TestObject(t *testing.T) {
	o := make(cti.Object)

	o.Put("/foo/bar", cti.Value("[\"qux\",\"baz\"]", "json"))
	o.Put("/ascii/48", cti.Value('a'))
	o.Put("/ascii/49", cti.Value('b'))
	o.Put("/ascii/50", cti.Value('c'))
	o.Put("/yaml", cti.Value("json: '{\"ini\":\"k=\\\"v\\\"\\nkey=\\\"value\\\"\\n\"}'\n"))

	got := o.Fields()
	want := cti.Fields{{
		Key: "/ascii",
	}, {
		Key:   "/ascii/48",
		Value: 'a',
	}, {
		Key:   "/ascii/49",
		Value: 'b',
	}, {
		Key:   "/ascii/50",
		Value: 'c',
	}, {
		Key: "/foo",
	}, {
		Key:      "/foo/bar",
		Encoding: cti.Encoding{"json"},
		Value:    "[\"qux\",\"baz\"]",
	}, {
		Key:   "/yaml",
		Value: "json: '{\"ini\":\"k=\\\"v\\\"\\nkey=\\\"value\\\"\\n\"}'\n",
	}}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	if err := o.Expand(); err != nil {
		t.Fatalf("o.Expand()=%s", err)
	}

	egot := o.Fields()
	ewant := cti.Fields{{
		Key: "/ascii",
	}, {
		Key:   "/ascii/48",
		Value: 'a',
	}, {
		Key:   "/ascii/49",
		Value: 'b',
	}, {
		Key:   "/ascii/50",
		Value: 'c',
	}, {
		Key: "/foo",
	}, {
		Key:      "/foo/bar",
		Encoding: cti.Encoding{"json"},
	}, {
		Key:   "/foo/bar/0",
		Value: "qux",
	}, {
		Key:   "/foo/bar/1",
		Value: "baz",
	}, {
		Key:      "/yaml",
		Encoding: cti.Encoding{"yaml"},
	}, {
		Key:      "/yaml/json",
		Encoding: cti.Encoding{"json"},
	}, {
		Key:      "/yaml/json/ini",
		Encoding: cti.Encoding{"ini"},
	}, {
		Key:   "/yaml/json/ini/k",
		Value: "v",
	}, {
		Key:   "/yaml/json/ini/key",
		Value: "value",
	}}

	if !cmp.Equal(egot, ewant) {
		t.Fatalf("egot != ewant:\n%s", cmp.Diff(egot, ewant))
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

	if !cmp.Equal(rgot.Keys(), rwant) {
		t.Fatalf("rgot != rwant:\n%s", cmp.Diff(rgot.Keys(), rwant))
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

	if !cmp.Equal(igot.Keys(), iwant) {
		t.Fatalf("igot != iwant:\n%s", cmp.Diff(igot.Keys(), iwant))
	}

	var cgot cti.Fields
	want[6].Encoding = cti.Encoding{"yaml"}

	if err := o.Compact(); err != nil {
		t.Fatalf("o.Expand()=%s", err)
	}

	o.Walk(cgot.Append)

	if !cmp.Equal(cgot, want) {
		t.Fatalf("cgot != want:\n%s", cmp.Diff(cgot, want))
	}
}
