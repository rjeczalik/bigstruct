package isrutil

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
