package isr_test

import (
	"fmt"
	"log"
	"os"

	"github.com/rjeczalik/bigstruct/isr"
	"github.com/rjeczalik/bigstruct/isr/codec"

	"gopkg.in/yaml.v3"
)

func ExampleObject_Put() {
	obj := make(isr.Object)

	obj.Put("/foo/bar", isr.Value(`{"key":[{"name":true},{"name":1}]}`))

	if err := obj.Decode(codec.Default); err != nil {
		log.Fatalf("obj.Decode()=%+v", err)
	}

	obj.Put("/foo/bar/key/1/args", isr.Value("--foo=bar --key=20"))

	if err := obj.Decode(codec.Default); err != nil {
		log.Fatalf("obj.Decode()=%+v", err)
	}

	obj.WriteTo(os.Stderr)

	p, err := yaml.Marshal(obj.Value())
	if err != nil {
		log.Fatalf("yaml.Marshal()=%+v", err)
	}

	fmt.Printf("%s\n", p)
	// Output:
	// foo:
	//     bar:
	//         key:
	//             - name: true
	//             - args:
	//                 --foo: bar
	//                 --key: 20
	//               name: 1
}
