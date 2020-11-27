package cti

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/glaucusio/confetti/internal/objects"
)

func Make(obj map[string]interface{}) Object {
	type elm struct {
		obj   map[string]interface{}
		nodes Object
	}

	root := make(Object)

	it, queue := elm{}, []elm{{obj, root}}

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

func MakeDir(dir string) (Object, error) {
	root := make(Object)

	err := filepath.Walk(dir, func(key string, fi os.FileInfo, err error) error {
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

		_ = root.Put(key, Field(AttrFile, newFile(fi), p))

		return nil
	})

	return root, err
}
