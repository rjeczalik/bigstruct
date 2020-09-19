package cti_test

import (
	"fmt"
	"log"

	"github.com/glaucusio/confetti/cti"

	"gopkg.in/yaml.v3"
)

func ExampleObject_Put() {
	obj := make(cti.Object)

	obj.Put("/foo/bar", cti.Value(`{"key":[{"name":"foo"},{"name":"bar"}]}`))

	if err := obj.Expand(); err != nil {
		log.Fatalf("obj.Expand()=%s", err)
	}

	obj.Put("/foo/bar/key/1/flag", cti.Value("--foo=bar --key=value"))

	if err := obj.Expand(); err != nil {
		log.Fatalf("obj.Expand()=%s", err)
	}

	p, err := yaml.Marshal(obj.Value())
	if err != nil {
		log.Fatalf("yaml.Marshal()=%s", err)
	}

	fmt.Printf("%s\n", p)
	// Output:
	// foo:
	//     bar:
	//         key:
	//             - name: foo
	//             - flag:
	//                 --foo: bar
	//                 --key: value
	//               name: bar
}
