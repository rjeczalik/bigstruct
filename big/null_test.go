package big_test

import (
	"testing"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/internal/types"
)

func TestNull(t *testing.T) {
	v := types.MakeYAML(big.NoValue).Value()

	if v != nil {
		t.Fatalf("got %#v, want nil", v)
	}
}
