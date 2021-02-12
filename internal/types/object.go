package types

import (
	"fmt"
	"sort"
	"strings"
)

type Object map[string]interface{}

func MakeObject(kv ...string) Object {
	o := make(Object, len(kv))

	for _, kv := range kv {
		var (
			k string
			v interface{}
		)

		if i := strings.IndexRune(kv, '='); i != -1 {
			k = kv[:i]
			v = YAML(kv[i+1:]).Value()
		} else {
			k = kv
		}

		if k != "" {
			o[k] = v
		}
	}

	return o
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

func (o Object) Merge(n Object) Object {
	for k, v := range n {
		if v == nil {
			delete(o, k)
		} else {
			o[k] = v
		}
	}

	return o
}

func (o Object) Keys() []string {
	keys := make([]string, 0, len(o))

	for k := range o {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (o Object) ReverseKeys() []string {
	keys := make([]string, 0, len(o))

	for k := range o {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	return keys

}

func (o Object) Map() map[string]interface{} {
	return o
}

func (o Object) Slice() []string {
	kv := o.Keys()

	for i, k := range kv {
		kv[i] = k + "=" + fmt.Sprint(o[k])
	}

	return kv
}
