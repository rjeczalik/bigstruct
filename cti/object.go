package cti

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"sort"
	"text/tabwriter"

	"github.com/glaucusio/confetti/internal/objects"
)

type Func func(key string, parent Object) (ok bool)

type Object map[string]struct {
	Encoding Encoding    `json:"e,omitempty" yaml:"e,omitempty"`
	Children Object      `json:"c,omitempty" yaml:"c,omitempty"`
	Value    interface{} `json:"v,omitempty" yaml:"v,omitempty"`
}

func (o Object) Copy() Object {
	u := make(Object, len(o))

	for k, n := range o {
		n.Encoding = n.Encoding.Copy()
		n.Children = n.Children.Copy()
		u[k] = n
	}

	return u
}

func (o Object) Meta() Object {
	u := make(Object, len(o))

	for k, n := range o {
		m, _ := u[k] // zero value
		m.Encoding = n.Encoding.Copy()
		m.Children = n.Children.Meta()
		u[k] = m
	}

	return u.Shake()
}

func (o Object) Fields() Fields {
	var f Fields
	o.Walk(f.Append)
	return f
}

func (o Object) Value() interface{} {
	obj := make(map[string]interface{})

	for name, node := range o {
		if len(node.Children) != 0 {
			obj[name] = node.Children.Value()
		} else {
			obj[name] = node.Value
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

func (o Object) Shake() Object {
	for k, n := range o {
		if len(n.Children) == 0 && n.Value == nil && len(n.Encoding) == 0 {
			delete(o, k)
		} else {
			n.Children = n.Children.Shake()
		}
	}

	if len(o) == 0 {
		return nil
	}

	return o
}

func (o Object) Put(key string, fn Func) (ok bool) {
	if key == "" || key == "/" {
		return false
	}

	var (
		parent = o
		dir    = path.Dir(key)
	)

	for _, k := range objects.Split(dir) {
		node, ok := parent[k]
		if !ok || node.Children == nil {
			node.Children = make(Object)
			parent[k] = node
		}
		parent = node.Children
	}

	return fn(key, parent)
}

func (o Object) At(key string, fn Func) (ok bool) {
	var (
		parent = o
		dir    = path.Dir(key)
	)

	for _, k := range objects.Split(dir) {
		node, ok := parent[k]
		if !ok || len(node.Children) == 0 {
			return false
		}
		parent = node.Children
	}

	return fn(key, parent)
}

func (o Object) Walk(fn Func) {
	type elm struct {
		parent Object
		key    string
		left   []string
	}

	if len(o) == 0 {
		return
	}

	var (
		it    elm
		k     string
		queue = []elm{{parent: o, key: "/", left: o.Keys()}}
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]
		k, it.left = it.left[0], it.left[1:]

		key := path.Join(it.key, k)

		if ok := fn(key, it.parent); !ok {
			return
		}

		if len(it.left) != 0 {
			queue = append(queue, it)
		}

		if parent := it.parent[k].Children; len(parent) != 0 {
			queue = append(queue, elm{parent: parent, key: key, left: parent.Keys()})
		}
	}
}

func (o Object) ReverseWalk(fn Func) {
	type elm struct {
		parent Object
		key    string
		left   []string
	}

	if len(o) == 0 {
		return
	}

	var (
		it    elm
		k     string
		queue = []elm{{parent: o, key: "/", left: o.Keys()}}
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

		fn(it.key, it.parent)
	}
}

func (o Object) ForEach(fn Func) {
	_ = o.forEach("/", fn)
}

func (o Object) forEach(key string, fn Func) (ok bool) {
	for _, k := range o.Keys() {
		if n, p := o[k], path.Join(key, k); len(n.Children) != 0 {
			ok = n.Children.forEach(p, fn)
		} else {
			ok = fn(p, o)
		}

		if !ok {
			return false
		}
	}

	return true
}

func (o Object) Keys() []string {
	keys := make([]string, 0, len(o))

	for k := range o {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (o Object) WriteTab(w io.Writer) (n int64, err error) {
	var m int

	if m, err = fmt.Fprintln(w, "KEY\tENCODING\tVALUE"); err != nil {
		return int64(m), err
	}

	o.Walk(func(key string, o Object) bool {
		var (
			k = path.Base(key)
			u = o[k]
		)

		if u.Value == nil && len(u.Children) != 0 && len(u.Encoding) == 0 {
			return true
		}

		m, err = fmt.Fprintf(w, "%s\t%s\t%+v\n",
			key,
			nonempty(u.Encoding.String(), "-"),
			nonil(u.Value, "-"),
		)

		n += int64(m)

		return err == nil
	})

	return n, err
}

func (o Object) String() string {
	var buf bytes.Buffer

	if _, err := o.WriteTo(&buf); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return buf.String()
}

func (o Object) WriteTo(w io.Writer) (int64, error) {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	n, err := o.WriteTab(tw)
	if err != nil {
		return n, err
	}

	if err := tw.Flush(); err != nil {
		return n, err
	}

	return n, err
}

func (o Object) Expand() error {
	return defaultSerializer.Expand(o)
}

func (o Object) Compact() error {
	return defaultSerializer.Compact(o)
}

func (o Object) Validate() error {
	return defaultSerializer.Validate(o)
}
