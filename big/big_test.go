package big_test

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/big/bigutil"
	"github.com/rjeczalik/bigstruct/big/codec"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

var updateGolden = flag.Bool("update-golden", false, "")

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func TestCti(t *testing.T) {
	orig, err := bigutil.MakeFile("testdata/docker")
	if err != nil {
		t.Fatalf("FileTree()=%s", err)
	}

	var (
		exp = orig.Copy()
	)

	if err := exp.Decode(codec.Default); err != nil {
		t.Fatalf("exp.Decode()=%s", err)
	}

	var (
		meta = reencode(exp.Schema())
		vgot = reencode(exp.Value())
	)

	if *updateGolden {
		if err := writeFile("testdata/docker.big.yaml.golden", exp); err != nil {
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

	if err := cpt.Encode(codec.Default); err != nil {
		t.Fatalf("cpt.Encode()=%s", err)
	}

	objwant, err := objReadFile("testdata/docker.big.yaml.golden")
	if err != nil {
		t.Fatalf("objReadFile()=%s", err)
	}

	if err := objwant.Encode(codec.Default); err != nil {
		t.Fatalf("objwant.Encode()=%s", err)
	}

	if got, want := objwant.Value(), cpt.Value(); !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	if err := objwant.Decode(codec.Default); err != nil {
		t.Fatalf("objwant.Decode()=%s", err)
	}

	if got, want := objwant.Schema(), exp.Schema(); !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	mgot := big.Make(exp.Value()).Merge(exp.Schema())

	if got, want := mgot.Value(), exp.Value(); !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	gdd := cpt.Copy().Merge(exp.Schema())

	gdd.WriteTo(os.Stdout)

	if err := gdd.Decode(codec.Default); err != nil {
		t.Fatalf("gdd.Decode()=%s", err)
	}

	if got, want := gdd.Value(), exp.Value(); !cmp.Equal(got, want) {
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

func objReadFile(file string) (big.Struct, error) {
	p, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var obj big.Struct

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
