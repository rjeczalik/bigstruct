package cti

import path "path"

func Value(v interface{}) Func {
	return func(key string, o Object) bool {
		var (
			k = path.Base(key)
			n = o[k]
		)

		n.Value = v
		o[k] = n

		return true
	}
}
