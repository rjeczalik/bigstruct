package ctiutil

import (
	"path"
	"sort"

	"github.com/glaucusio/confetti/cti"
)

type Field struct {
	Key   string
	Attr  cti.Attr
	Aux   interface{}
	Value interface{}
}

type Fields []Field

var (
	_ cti.Func       = (*Fields)(nil).Append
	_ sort.Interface = (*Fields)(nil)
)

func (f *Fields) Append(key string, u cti.Object) bool {
	var (
		k = path.Base(key)
		n = u[k]
	)

	*f = append(*f, Field{
		Key:   key,
		Attr:  n.Attr,
		Aux:   n.Aux,
		Value: n.Value,
	})

	return true
}

func (f Fields) Keys() []string {
	keys := make([]string, 0, len(f))

	for _, f := range f {
		keys = append(keys, f.Key)
	}

	return keys
}

func (f Fields) Object() cti.Object {
	o := make(cti.Object)

	for _, f := range f {
		o.Put(f.Key, cti.Field(f.Attr, f.Aux, f.Value))
	}

	return o
}

func (f Fields) Len() int {
	return len(f)
}

func (f Fields) Less(i, j int) bool {
	return f[i].Key < f[j].Key
}

func (f Fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
