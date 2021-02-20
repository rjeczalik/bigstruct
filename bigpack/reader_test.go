package bigpack_test

import (
	"testing"

	"github.com/rjeczalik/bigstruct/bigpack"
	"github.com/rjeczalik/bigstruct/bigpack/pak"
)

func TestReader(t *testing.T) {
	p, err := bigpack.Read(pak.Dir("testdata/bigpack-scylla"))
	if err != nil {
		t.Fatalf("Read()=%+v", err)
	}

	_ = p
}
