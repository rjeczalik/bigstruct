package isrutil

import (
	"io/ioutil"
	"os"
	stdpath "path"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/isr"
)

func MakeFile(path string) (isr.Object, error) {
	type index map[string]string

	var f isr.Fields

	switch fi, err := os.Stat(path); {
	case err != nil:
		return nil, err
	case fi.IsDir():
		err := filepath.Walk(path, func(key string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if fi.IsDir() {
				return nil
			}

			p, err := ioutil.ReadFile(key)
			if err != nil {
				return err
			}

			if filepath.Base(key) == ".bigstruct.index" {
				var idx index

				if err := types.YAML(p).Unmarshal(&idx); err != nil {
					return err
				}

				for k, typ := range idx {
					f = append(f, isr.Field{
						Key:  stdpath.Join(cleanpath(strings.TrimPrefix(key, path)), k),
						Type: typ,
					})
				}

				return nil
			}

			f = append(f, isr.Field{
				Key:   cleanpath(strings.TrimPrefix(key, path)),
				Value: string(p),
			})

			return nil
		})
		if err != nil {
			return nil, err
		}
	default:
		p, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		f = append(f, isr.Field{
			Key:   cleanpath(filepath.Base(path)),
			Value: p,
		})
	}

	return f.Object(), nil
}
