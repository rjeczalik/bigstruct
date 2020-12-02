package cti_test

import (
	"flag"
	"io/ioutil"
	"os"
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

func TestCti(t *testing.T) {
	orig, err := cti.MakeDir("testdata/docker")
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

	var (
		meta = reencode(exp.Meta())
		vgot = reencode(exp.Value())
	)

	if *updateGolden {
		if err := writeFile("testdata/docker.object.yaml.golden", exp); err != nil {
			t.Fatalf("writeFile()=%s", err)
		}

		if err := writeFile("testdata/docker.yaml.golden", vgot); err != nil {
			t.Fatalf("writeFile()=%s", err)
		}

		if err := writeFile("testdata/docker.meta.yaml.golden", meta); err != nil {
			t.Fatalf("writeFile()=%s", err)
		}

		return
	}

	vwant, err := vReadFile("testdata/docker.yaml.golden")
	if err != nil {
		t.Fatalf("vReadFile()=%s", err)
	}

	if !cmp.Equal(vgot, vwant) {
		t.Fatalf("vgot != vwant:\n%s", cmp.Diff(vgot, vwant))
	}

	awant, err := vReadFile("testdata/docker.meta.yaml.golden")
	if err != nil {
		t.Fatalf("vReadFile()=%s", err)
	}

	if !cmp.Equal(meta, awant) {
		t.Fatalf("got != want:\n%s", cmp.Diff(meta, awant))
	}

	cpt := exp.Copy()

	if err := s.Compact(cpt); err != nil {
		t.Fatalf("s.Compact()=%s", err)
	}

	objwant, err := objReadFile("testdata/docker.object.yaml.golden")
	if err != nil {
		t.Fatalf("objReadFile()=%s", err)
	}

	if err := s.Compact(objwant); err != nil {
		t.Fatalf("s.Compact()=%s", err)
	}

	if got, want := objwant.Value(), cpt.Value(); !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	if err := s.Expand(objwant); err != nil {
		t.Fatalf("s.Expand()=%s", err)
	}

	if got, want := objwant.Meta(), exp.Meta(); !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}
}

func writeFile(file string, v interface{}) error {
	p, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, p, 0644)
}

func vReadFile(file string) (interface{}, error) {
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

func objReadFile(file string) (cti.Object, error) {
	p, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var obj cti.Object

	if err := yaml.Unmarshal(p, &obj); err != nil {
		return nil, err
	}

	return obj, nil
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
