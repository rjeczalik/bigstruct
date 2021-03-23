package bigutil

import (
	"io/ioutil"
	"os"
	stdpath "path"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/internal/types"
)

func MakeFile(path string) (big.Struct, error) {
	type index map[string]string

	var f big.Fields

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

			if filepath.Base(key) == ".bigpack" {
				var idx index

				if err := types.YAML(p).Unmarshal(&idx); err != nil {
					return err
				}

				for k, typ := range idx {
					f = append(f, big.Field{
						Key:  stdpath.Join(cleanpath(strings.TrimPrefix(key, path)), k),
						Type: typ,
					})
				}

				return nil
			}

			f = append(f, big.Field{
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

		f = append(f, big.Field{
			Key:   cleanpath(filepath.Base(path)),
			Value: p,
		})
	}

	return f.Struct(), nil
}
