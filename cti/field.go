package cti

import (
	"path"
	"sort"
)

type Field struct {
	Key      string
	Encoding string
	Value    interface{}
}

func Value(v interface{}, encoding ...string) Func {
	return Field{
		Encoding: path.Join(encoding...),
		Value:    v,
	}.Put
}

func (f Field) Put(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	if f.Encoding != "" {
		n.Encoding = f.Encoding
	}

	if f.Value != nil {
		n.Value = f.Value
	}

	o[k] = n

	return nil
}

type Fields []Field

var (
	_ Func           = (*Fields)(nil).Append
	_ sort.Interface = (*Fields)(nil)
)

func (f *Fields) Append(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	*f = append(*f, Field{
		Key:      key,
		Encoding: n.Encoding,
		Value:    n.Value,
	})

	return nil
}

func (f Fields) Keys() []string {
	keys := make([]string, 0, len(f))

	for _, f := range f {
		keys = append(keys, f.Key)
	}

	return keys
}

func (f Fields) Object() Object {
	o := make(Object)

	for _, f := range f {
		o.Put(f.Key, f.Put)
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

func (f Fields) Sort() Fields {
	sort.Stable(f)
	return f
}
