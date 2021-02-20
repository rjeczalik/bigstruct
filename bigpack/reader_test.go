package bigpack_test

import (
	"testing"

	"github.com/rjeczalik/bigstruct/bigpack"
	"github.com/rjeczalik/bigstruct/bigpack/pak"
	"github.com/rjeczalik/bigstruct/storage"
)

func TestReader(t *testing.T) {
	const uri = "sqlite://file::memory:?cache=shared"

	pk, err := bigpack.Read(pak.Dir("testdata/bigpack-scylla"))
	if err != nil {
		t.Fatalf("Read()=%+v", err)
	}

	g, err := storage.Open(uri)
	if err != nil {
		t.Fatalf("Open()=%+v", err)
	}
	defer g.Close()

	if err := g.Transaction(pk.Store); err != nil {
		t.Fatalf("Transaction()=%+v", err)
	}
}
