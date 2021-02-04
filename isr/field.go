package isr

import (
	"path"
	"sort"
)

type Field struct {
	Key   string
	Type  string
	Value interface{}
}

func Value(v interface{}, typ ...string) Func {
	return Field{
		Type:  path.Join(typ...),
		Value: v,
	}.Put
}

func Children(o Object) Func {
	return func(key string, u Object) error {
		var (
			k = path.Base(key)
			n = o[k]
		)

		n.Children = o
		o[k] = n

		return nil
	}
}

func (f Field) Put(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	if f.Type != "" {
		n.Type = f.Type
	}

	if f.Value != nil {
		n.Value = f.Value
	}

	o[k] = n

	return nil
}

func (f Field) Set(key string, o Object) error {
	var (
		k = path.Base(key)
		n = o[k]
	)

	n.Type = f.Type
	n.Value = f.Value
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
		Key:   key,
		Type:  n.Type,
		Value: n.Value,
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

func (f Fields) Merge() Object {
	o := make(Object)

	for _, f := range f {
		o.Put(f.Key, f.Set)
	}

	return o.Shake()
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
