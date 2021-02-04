package isr

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/bigstruct/internal/objects"
)

func Make(v interface{}) Object {
	type elm struct {
		obj   map[string]interface{}
		nodes Object
	}

	var (
		root  = make(Object)
		obj   = objects.Object(v)
		it    elm
		queue = []elm{{obj, root}}
	)

	if obj == nil {
		panic(fmt.Errorf("invalid type for v (%T); expected generic map or slice", v))
	}

	for len(queue) != 0 {
		it, queue = queue[0], queue[1:]

		for k, v := range it.obj {
			node, ok := it.nodes[k]
			if !ok {
				node.Children = make(Object)
			}

			jt := it
			jt.nodes = node.Children

			if obj := objects.Object(v); obj != nil {
				jt.obj = obj
				queue = append(queue, jt)
			} else {
				node.Value = v
				node.Children = nil
			}

			it.nodes[k] = node
		}
	}

	return root
}

func MakeFile(path string) (Object, error) {
	var f Fields

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

			f = append(f, Field{
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

		f = append(f, Field{
			Key:   cleanpath(filepath.Base(path)),
			Value: p,
		})
	}

	return f.Object(), nil
}
