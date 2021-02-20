package big

import (
	"fmt"

	"github.com/rjeczalik/bigstruct/internal/objects"
)

func Make(v interface{}) Struct {
	type elm struct {
		obj   map[string]interface{}
		nodes Struct
	}

	var (
		root  = make(Struct)
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
				node.Children = make(Struct)
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
