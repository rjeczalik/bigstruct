package isr_test

import (
	"testing"

	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/isr"
)

func TestNull(t *testing.T) {
	v := types.MakeYAML(isr.NoValue).Value()

	if v != nil {
		t.Fatalf("got %#v, want nil", v)
	}
}
