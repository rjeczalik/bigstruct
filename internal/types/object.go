package types

import (
	"sort"
	"strings"
)

type Object map[string]interface{}

func MakeObject(kv ...string) Object {
	root := make(Object, len(kv))

	for _, kv := range kv {
		var (
			key = kv
			val interface{}
		)

		if i := strings.IndexRune(kv, '='); i != -1 {
			key = kv[:i]
			val = YAML(kv[i+1:]).Value()
		}

		var (
			o  = root
			it = o
			ok bool
			k  string
		)

		for _, k = range stringsSplit(key, ".") {
			o = it

			if it, ok = o[k].(Object); !ok {
				it = make(Object)
				o[k] = it
			}
		}

		if k != "" {
			o[k] = val
		}
	}

	return root
}

func (o Object) JSON() JSON {
	return MakeJSON(o)
}

func (o Object) YAML() YAML {
	return MakeYAML(o)
}

func (o Object) Copy() Object {
	n := make(Object, len(o))

	for k, v := range o {
		n[k] = v
	}

	return n
}

func (o Object) Merge(u Object) Object {
	return MakeObject(append(o.Slice(), u.Slice()...)...)
}

func (o Object) Slice() []string {
	type elm struct {
		key []string
		o   Object
	}

	var (
		it    elm
		queue = []elm{{o: o}}
		slice []string
	)

	for len(queue) != 0 {
		it, queue = queue[0], queue[1:]

		for k, v := range it.o {
			switch v := v.(type) {
			case Object:
				queue = append(queue, elm{o: v, key: append(it.key, k)})
			default:
				if v != nil {
					slice = append(slice, strings.Join(append(it.key, k), ".")+"="+MakeYAML(v).String())
				} else {
					slice = append(slice, strings.Join(append(it.key, k), "."))
				}
			}
		}
	}

	sort.Strings(slice)

	return slice
}

func (o Object) Map() map[string]interface{} {
	return o
}

func stringsSplit(s, sep string) (slice []string) {
	for _, s := range strings.Split(s, sep) {
		if s = strings.TrimSpace(s); s != "" {
			slice = append(slice, s)
		}
	}
	return slice
}
