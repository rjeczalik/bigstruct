package big

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"sort"
	"text/tabwriter"
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

func (f Field) String() string {
	switch v := f.Value.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	default:
		panic(fmt.Errorf("unable to convert %T to string", f.Value))
	}
}

func (f Field) Bytes() []byte {
	switch v := f.Value.(type) {
	case nil:
		return nil
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		panic(fmt.Errorf("unable to convert %T to []byte", f.Value))
	}
}

type Fields []Field

var (
	_ Func           = (*Fields)(nil).Append
	_ sort.Interface = (*Fields)(nil)
)

func (f Fields) WriteTab(w io.Writer) (n int64, err error) {
	m, err := fmt.Fprintln(w, "KEY\tTYPE\tVALUE")
	if err != nil {
		return int64(n), err
	}

	n += int64(m)

	for _, f := range f {
		m, err := fmt.Fprintf(w, "%s\t%s\t%+v\n",
			f.Key,
			nonempty(f.Type, "-"),
			nonil(f.Value, "-"),
		)

		n += int64(m)

		if err != nil {
			return n, err
		}
	}

	return n, err
}

func (f Fields) String() string {
	var buf bytes.Buffer

	if _, err := f.WriteTo(&buf); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return buf.String()
}

func (f Fields) WriteTo(w io.Writer) (int64, error) {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	n, err := f.WriteTab(tw)
	if err != nil {
		return n, err
	}

	if err := tw.Flush(); err != nil {
		return n, err
	}

	return n, err
}

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

func (f *Fields) AppendIf(key string, s Struct) error {
	var (
		k = path.Base(key)
		n = s[k]
	)

	if len(n.Children) != 0 && (n.Type != "" || n.Value != nil) {
		*f = append(*f, Field{
			Key:   key,
			Type:  n.Type,
			Value: n.Value,
		})
	}

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

func (f Fields) At(i int) Field {
	if i < len(f) {
		return f[i]
	}
	return Field{}
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
