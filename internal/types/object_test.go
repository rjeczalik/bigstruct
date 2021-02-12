package types_test

import (
	"testing"

	"github.com/rjeczalik/bigstruct/internal/types"

	"github.com/google/go-cmp/cmp"
)

func TestMakeObject(t *testing.T) {
	kv := []string{
		"a.b.c=1",
		"a.d.e=2",
		"a.d.f.g=3",
		"a.f.g.h=4",
		"a.f.g=5",
		"....=6",
		"a.d.g",
		"b=7",
		"c.d.e=8",
		"c.foo",
	}

	want := types.Object{
		"a": types.Object{
			"b": types.Object{
				"c": 1,
			},
			"d": types.Object{
				"e": 2,
				"f": types.Object{
					"g": 3,
				},
				"g": nil,
			},
			"f": types.Object{
				"g": 5,
			},
		},
		"b": 7,
		"c": types.Object{
			"d": types.Object{
				"e": 8,
			},
			"foo": nil,
		},
	}

	got := types.MakeObject(kv...)

	if !cmp.Equal(want, got) {
		t.Fatalf("want != got:\n%s", cmp.Diff(want, got))
	}

	w := []string{
		"a.b.c=1",
		"a.d.e=2",
		"a.d.f.g=3",
		"a.d.g",
		"a.f.g=5",
		"b=7",
		"c.d.e=8",
		"c.foo",
	}

	g := got.Slice()

	if !cmp.Equal(w, g) {
		t.Fatalf("w != g:\n%s", cmp.Diff(w, g))
	}
}

func TestObjectMerge(t *testing.T) {
	o := types.Object{
		"a": types.Object{
			"b": types.Object{
				"c": 1,
			},
			"d": 2,
		},
		"b": 3,
		"e": types.Object{
			"f": 4,
		},
		"g": types.Object{
			"g": types.Object{
				"f": types.Object{
					"f": types.Object{
						"foo": "bar",
					},
				},
			},
		},
	}

	u := types.Object{
		"a": types.Object{
			"b": types.Object{
				"d": -1,
			},
			"d": types.Object{
				"foo": "bar",
			},
		},
		"b": -2,
		"e": -3,
		"g": types.Object{
			"g": nil,
		},
		"h": -4,
		"i": types.Object{
			"j": -5,
		},
	}

	want := types.Object{
		"a": types.Object{
			"b": types.Object{
				"c": 1,
				"d": -1,
			},
			"d": types.Object{
				"foo": "bar",
			},
		},
		"b": -2,
		"e": -3,
		"g": types.Object{
			"g": nil,
		},
		"h": -4,
		"i": types.Object{
			"j": -5,
		},
	}

	got := o.Merge(u)

	if !cmp.Equal(want, got) {
		t.Fatalf("want != got:\n%s", cmp.Diff(want, got))
	}
}
