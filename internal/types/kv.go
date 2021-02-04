package types

import (
	"sort"
	"strings"
)

type KV map[string]string

func MakeKV(kv ...string) KV {
	m := make(KV, len(kv))

	for _, kv := range kv {
		var k, v string

		if i := strings.IndexRune(kv, '='); i != -1 {
			k, v = kv[:i], kv[i+1:]
		} else {
			k = kv
		}

		if k != "" {
			m[k] = v
		}
	}

	return m
}

func MergeKV(m ...map[string]string) map[string]string {
	var kv = make(KV)

	for _, m := range m {
		kv = kv.Merge(KV(m))
	}

	return kv.Map()
}

func (m KV) JSON() JSON {
	return MakeJSON(m)
}

func (m KV) YAML() YAML {
	return MakeYAML(m)
}

func (m KV) Copy() KV {
	n := make(KV, len(m))

	for k, v := range m {
		n[k] = v
	}

	return n
}

func (m KV) Merge(n KV) KV {
	for k, v := range n {
		if v == "" {
			delete(m, k)
		} else {
			m[k] = v
		}
	}

	return m
}

func (m KV) Keys() []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (m KV) ReverseKeys() []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	return keys

}

func (m KV) Map() map[string]string {
	return m
}

func (m KV) Slice() []string {
	kv := m.Keys()

	for i, k := range kv {
		kv[i] = k + "=" + m[k]
	}

	return kv
}
