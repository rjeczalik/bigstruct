package pak

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
)

type Func func(path string, rc io.ReadCloser) error

type FS func(fn Func) error

func Dir(dir string) FS {
	return func(fn Func) error {
		return filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if fi.IsDir() {
				return nil // skip
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			if err := fn(strings.TrimPrefix(path, dir), f); err != nil {
				return err
			}

			return nil
		})
	}
}

func Packr(prefix string, b *packr.Box) FS {
	return func(fn Func) error {
		return b.WalkPrefix(prefix, func(path string, f packd.File) error {
			defer f.Close()
			return fn(strings.TrimPrefix(prefix, path), f)
		})
	}
}
