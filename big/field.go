package big

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

func Children(s Struct) Func {
	return func(key string, u Struct) error {
		var (
			k = path.Base(key)
			n = u[k]
		)

		n.Children = s
		u[k] = n

		return nil
	}
}

func (f Field) Put(key string, s Struct) error {
	var (
		k = path.Base(key)
		n = s[k]
	)

	if f.Type != "" {
		n.Type = f.Type
	}

	if f.Value != nil {
		n.Value = f.Value
	}

	s[k] = n

	return nil
}

func (f Field) Set(key string, s Struct) error {
	var (
		k = path.Base(key)
		n = s[k]
	)

	n.Type = f.Type
	n.Value = f.Value
	s[k] = n

	return nil
}

type Fields []Field

var (
	_ Func           = (*Fields)(nil).Append
	_ sort.Interface = (*Fields)(nil)
)

func (f *Fields) Append(key string, s Struct) error {
	var (
		k = path.Base(key)
		n = s[k]
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

func (f Fields) Struct() Struct {
	s := make(Struct)

	for _, f := range f {
		s.Put(f.Key, f.Put)
	}

	return s
}

func (f Fields) Merge() Struct {
	s := make(Struct)

	for _, f := range f {
		s.Put(f.Key, f.Set)
	}

	return s.Shake()
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
