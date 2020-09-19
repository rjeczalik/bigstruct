package cti_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	path "path"
	"testing"

	"github.com/glaucusio/confetti/cti"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

var updateGolden = flag.Bool("update-golden", false, "")

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func TestSerializer(t *testing.T) {
	orig, err := cti.MakeDir("testdata/file")
	if err != nil {
		t.Fatalf("FileTree()=%s", err)
	}

	var (
		s   = cti.NewSerializer()
		exp = orig.Copy()
	)

	if err := s.Expand(exp); err != nil {
		t.Fatalf("s.Expand()=%s", err)
	}

	exp.Attr().Walk(debug)

	var (
		attr = reencode(exp.Attr())
		got  = reencode(exp.Value())
	)

	if *updateGolden {
		if err := writeFile("testdata/file.yaml.golden", got); err != nil {
			t.Fatalf("writeFile()=%s", err)
		}

		if err := writeFile("testdata/file-attr.yaml.golden", attr); err != nil {
			t.Fatalf("writeFile()=%s", err)
		}

		return
	}

	want, err := readFile("testdata/file.yaml.golden")
	if err != nil {
		t.Fatalf("readFile()=%s", err)
	}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	awant, err := readFile("testdata/file-attr.yaml.golden")
	if err != nil {
		t.Fatalf("readFile()=%s", err)
	}

	if !cmp.Equal(attr, awant) {
		t.Fatalf("got != want:\n%s", cmp.Diff(attr, awant))
	}

	cpt := exp.Copy()

	if err := s.Compact(cpt); err != nil {
		t.Fatalf("s.Compact()=%s", err)
	}

	// cpt.Walk(debug)
}

func writeFile(file string, v interface{}) error {
	p, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, p, 0644)
}

func readFile(file string) (interface{}, error) {
	p, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var v interface{}

	if err := yaml.Unmarshal(p, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func reencode(v interface{}) (w interface{}) {
	p, err := yaml.Marshal(v)
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	if err := yaml.Unmarshal(p, &w); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return w
}

func debug(key string, o cti.Object) bool {
	var (
		k = path.Base(key)
		n = o[k]
	)

	switch {
	case len(n.Children) != 0 && n.Attr != 0:
		fmt.Printf("%s (%s)\n", key, n.Attr)
	case len(n.Children) != 0:
		fmt.Printf("%s\n", key)
	case n.Value != nil && n.Attr != 0:
		fmt.Printf("%s=%s (%s)\n", key, n.Value, n.Attr)
	case n.Attr != 0:
		fmt.Printf("%s (%s)\n", key, n.Attr)
	default:
		fmt.Printf("%s=%s\n", key, n.Value)
	}

	return true
}
