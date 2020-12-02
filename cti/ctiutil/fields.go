package ctiutil

import (
	"path"
	"sort"

	"github.com/glaucusio/confetti/cti"
)

type Field struct {
	Key      string
	Encoding cti.Encoding
	Value    interface{}
}

func (f Field) Set(key string, o cti.Object) bool {
	var (
		k = path.Base(key)
		n = o[k]
	)

	n.Encoding = f.Encoding
	n.Value = f.Value
	o[k] = n

	return true
}

type Fields []Field

var (
	_ cti.Func       = (*Fields)(nil).Append
	_ sort.Interface = (*Fields)(nil)
)

func (f *Fields) Append(key string, o cti.Object) bool {
	var (
		k = path.Base(key)
		n = o[k]
	)

	*f = append(*f, Field{
		Key:      key,
		Encoding: n.Encoding,
		Value:    n.Value,
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
		o.Put(f.Key, f.Set)
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
