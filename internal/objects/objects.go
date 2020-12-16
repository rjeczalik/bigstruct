package objects

import (
	"fmt"
	"strconv"
	"strings"
)

func Split(dir string) (v []string) {
	for _, s := range strings.Split(strings.Trim(dir, "/"), "/") {
		if s != "" {
			v = append(v, s)
		}
	}
	return v
}

func Object(v interface{}) map[string]interface{} {
	var m map[string]interface{}

	switch v := v.(type) {
	case map[string]interface{}:
		if m = v; m == nil {
			m = make(map[string]interface{})
		}
	case map[interface{}]interface{}: // go-yaml/yaml#139
		m = make(map[string]interface{}, len(v))

		for k, v := range v {
			m[fmt.Sprint(k)] = v
		}
	case []interface{}:
		m = make(map[string]interface{}, len(v))

		for i, v := range v {
			m[fmt.Sprint(i)] = v
		}
	}

	return m
}

func Slice(obj map[string]interface{}) []interface{} {
	var max int

	for k := range obj {
		if n, err := strconv.Atoi(k); err != nil {
			return nil
		} else if n > max {
			max = n
		}
	}

	a := make([]interface{}, max+1)

	for k, v := range obj {
		i, _ := strconv.Atoi(k)
		a[i] = v
	}

	return a
}
