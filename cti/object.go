package cti

import (
	"path"
	"sort"

	"github.com/glaucusio/confetti/internal/objects"
)

type Func func(key string, parent Object) (ok bool)

type Object map[string]struct {
	Attr     Attr        `json:"a,omitempty" yaml:"a,omitempty"`
	Children Object      `json:"c,omitempty" yaml:"c,omitempty"`
	Value    interface{} `json:"v,omitempty" yaml:"v,omitempty"`
}

func (o Object) Copy() Object {
	u := make(Object, len(o))

	for k, n := range o {
		n.Children = n.Children.Copy()
		u[k] = n
	}

	return u
}

func (o Object) Attr() Object {
	u := make(Object, len(o))

	for k, n := range o {
		m, _ := u[k] // zero value
		m.Attr = n.Attr
		m.Children = n.Children.Attr()
		u[k] = m
	}

	return u.Shake()
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

	return obj
}

func (o Object) Shake() Object {
	for k, n := range o {
		if len(n.Children) == 0 && n.Value == nil && n.Attr == 0 {
			delete(o, k)
		} else {
			n.Children = n.Children.Shake()
		}
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
	}

	var (
		it    elm
		queue = []elm{{parent: o, key: "/"}}
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]

		for _, k := range it.parent.Keys() {
			key := path.Join(it.key, k)

			if ok := fn(key, it.parent); !ok {
				return
			}

			if parent := it.parent[k].Children; len(parent) != 0 {
				queue = append(queue, elm{parent: parent, key: key})
			}
		}
	}
}

func (o Object) ReverseWalk(fn Func) {
	type elm struct {
		parent Object
		key    string
	}

	var (
		it    elm
		queue = []elm{{parent: o, key: "/"}}
		rev   []elm
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]

		for _, k := range it.parent.Keys() {
			key := path.Join(it.key, k)

			rev = append(rev, elm{parent: it.parent, key: key})

			if parent := it.parent[k].Children; len(parent) != 0 {
				queue = append(queue, elm{parent: parent, key: key})
			}
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

func (o Object) Expand() error {
	return defaultSerializer.Expand(o)
}

func (o Object) Compact() error {
	return defaultSerializer.Compact(o)
}
