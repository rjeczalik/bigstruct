package isr

import (
	"path"
	"path/filepath"
	"strings"
)

func cleanpath(s string) string {
	s = strings.TrimLeft(s, `/.\`)
	s = filepath.FromSlash(s)
	s = filepath.Clean(s)
	s = filepath.ToSlash(s)
	s = path.Join("/", s)

	return s
}

func nonempty(s ...string) string {
	for _, s := range s {
		if s != "" {
			return s
		}
	}
	return ""
}

func nonil(v ...interface{}) interface{} {
	for _, v := range v {
		if v != nil {
			return v
		}
	}
	return nil
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
