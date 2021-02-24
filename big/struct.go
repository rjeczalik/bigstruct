package big

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"sort"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/internal/objects"
)

type Func func(key string, parent Struct) error

type Struct map[string]struct {
	Type     string      `json:"type,omitempty" yaml:"type,omitempty"`
	Value    interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	Children Struct      `json:"children,omitempty" yaml:"children,omitempty"`
}

func Move(key string, s Struct) Struct {
	if key == "" || key == "/" {
		return s
	}

	root := make(Struct)

	if err := root.Put(key, s.Move); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return root
}

func (s Struct) Copy() Struct {
	u := make(Struct, len(s))

	for k, n := range s {
		n.Children = n.Children.Copy()
		u[k] = n
	}

	return u
}

func (s Struct) Schema() Struct {
	u := make(Struct, len(s))

	for k, n := range s {
		m, _ := u[k] // zero value
		m.Type = n.Type
		m.Children = n.Children.Schema()
		u[k] = m
	}

	return u.Shake()
}

func (s Struct) Raw() Struct {
	s.Walk(func(key string, s Struct) error {
		var (
			k = path.Base(key)
			n = s[k]
		)

		n.Type = ""
		s[k] = n

		return nil
	})

	return s.Shake()
}

func (s Struct) Merge(u Struct) Struct {
	return append(s.Fields(), u.Fields()...).Struct()
}

func (s Struct) Fields() Fields {
	var f Fields
	s.ForEach(f.Append)
	s.ReverseWalk(f.AppendIf)
	return f
}

func (s Struct) Value() interface{} {
	obj := make(map[string]interface{})

	for name, node := range s {
		if node.Value != nil {
			obj[name] = node.Value
		} else {
			obj[name] = node.Children.Value()
		}
	}

	if a := objects.Slice(obj); len(a) != 0 {
		return a
	}

	if len(obj) == 0 {
		return nil
	}

	return obj
}

func (s Struct) Shake() Struct {
	for k, n := range s {
		if len(n.Children) == 0 && n.Value == nil && n.Type == "" {
			delete(s, k)
		} else {
			n.Children = n.Children.Shake()
		}
	}

	if len(s) == 0 {
		return nil
	}

	return s
}

func (s Struct) ShakeTypes() Struct {
	for k, n := range s {
		if len(n.Children) == 0 && n.Value == nil {
			delete(s, k)
		} else {
			n.Children = n.Children.ShakeTypes()
		}
	}

	if len(s) == 0 {
		return nil
	}

	return s
}

func (s Struct) Move(key string, u Struct) error {
	var (
		k = path.Base(key)
		n = u[k]
	)

	n.Children = s
	u[k] = n

	return nil
}

func (s Struct) Put(key string, fn Func) error {
	if key == "" || key == "/" {
		return nil // fixme: error invalid key
	}

	var (
		parent = s
		dir    = path.Dir(key)
	)

	for _, k := range objects.Split(dir) {
		node, ok := parent[k]
		if !ok || node.Children == nil {
			node.Children = make(Struct)
			parent[k] = node
		}
		parent = node.Children
	}

	return fn(key, parent)
}

func (s Struct) At(key string) Struct {
	var (
		parent = s
	)

	for _, k := range objects.Split(key) {
		node, ok := parent[k]
		if !ok || len(node.Children) == 0 {
			return nil // fixme: error not found
		}
		parent = node.Children
	}

	return parent
}

func (s Struct) ValueAt(key string) interface{} {
	var (
		dir  = path.Dir(key)
		base = path.Base(key)
	)

	return s.At(dir)[base].Value
}

func (s Struct) TypeAt(key string) string {
	var (
		dir  = path.Dir(key)
		base = path.Base(key)
	)

	return s.At(dir)[base].Type
}

func (s Struct) Walk(fn Func) error {
	type elm struct {
		parent Struct
		key    string
		left   []string
	}

	if len(s) == 0 {
		return nil
	}

	var (
		it    elm
		k     string
		queue = []elm{{parent: s, key: "/", left: s.Keys()}}
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]
		k, it.left = it.left[0], it.left[1:]

		key := path.Join(it.key, k)

		if err := fn(key, it.parent); err != nil {
			return err
		}

		if len(it.left) != 0 {
			queue = append(queue, it)
		}

		if parent := it.parent[k].Children; len(parent) != 0 {
			queue = append(queue, elm{parent: parent, key: key, left: parent.Keys()})
		}
	}

	return nil
}

func (s Struct) ReverseWalk(fn Func) error {
	type elm struct {
		parent Struct
		key    string
		left   []string
	}

	if len(s) == 0 {
		return nil
	}

	var (
		it    elm
		k     string
		queue = []elm{{parent: s, key: "/", left: s.Keys()}}
		rev   []elm
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]
		k, it.left = it.left[0], it.left[1:]

		key := path.Join(it.key, k)

		rev = append(rev, elm{parent: it.parent, key: key})

		if len(it.left) != 0 {
			queue = append(queue, it)
		}

		if parent := it.parent[k].Children; len(parent) != 0 {
			queue = append(queue, elm{parent: parent, key: key, left: parent.Keys()})
		}
	}

	for len(rev) != 0 {
		it, rev = rev[len(rev)-1], rev[:len(rev)-1]

		if err := fn(it.key, it.parent); err != nil {
			return err
		}
	}

	return nil
}

func (s Struct) ForEach(fn Func) error {
	return s.forEach("/", fn)
}

func (s Struct) forEach(key string, fn Func) (err error) {
	for _, k := range s.Keys() {
		if n, p := s[k], path.Join(key, k); len(n.Children) != 0 {
			err = n.Children.forEach(p, fn)
		} else {
			err = fn(p, s)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s Struct) Keys() []string {
	keys := make([]string, 0, len(s))

	for k := range s {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (s Struct) WriteTab(w io.Writer) (n int64, err error) {
	m, err := fmt.Fprintln(w, "KEY\tTYPE\tVALUE")
	if err != nil {
		return int64(n), err
	}

	n += int64(m)

	err = s.Walk(func(key string, s Struct) error {
		var (
			k = path.Base(key)
			u = s[k]
		)

		if u.Value == nil && len(u.Children) != 0 && len(u.Type) == 0 {
			return nil
		}

		m, err := fmt.Fprintf(w, "%s\t%s\t%+v\n",
			key,
			nonempty(u.Type, "-"),
			nonil(u.Value, "-"),
		)

		n += int64(m)

		return err
	})

	return n, err
}

func (s Struct) String() string {
	var buf bytes.Buffer

	if _, err := s.WriteTo(&buf); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return buf.String()
}

func (s Struct) WriteTo(w io.Writer) (int64, error) {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	n, err := s.WriteTab(tw)
	if err != nil {
		return n, err
	}

	if err := tw.Flush(); err != nil {
		return n, err
	}

	return n, err
}

func (s Struct) Encode(ctx context.Context, c Codec) error {
	if c == nil {
		panic("codec is nil")
	}
	return s.ReverseWalk(func(key string, s Struct) error {
		return c.Encode(ctx, key, s)
	})
}

func (s Struct) Decode(ctx context.Context, c Codec) error {
	if c == nil {
		panic("codec is nil")
	}
	return s.Walk(func(key string, s Struct) error {
		return c.Decode(ctx, key, s)
	})
}
