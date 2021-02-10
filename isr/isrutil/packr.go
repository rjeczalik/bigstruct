package isrutil

import (
	"io/ioutil"

	"github.com/rjeczalik/bigstruct/isr"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
)

func MakeBox(box *packr.Box) (isr.Object, error) {
	var f isr.Fields

	err := box.Walk(func(path string, r packd.File) error {
		p, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		f = append(f, isr.Field{
			Key:   cleanpath(path),
			Value: string(p),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return f.Object(), nil
}
