package bigutil

import (
	"io/ioutil"

	"github.com/rjeczalik/bigstruct/big"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
)

func MakeBox(box *packr.Box) (big.Struct, error) {
	var f big.Fields

	err := box.Walk(func(path string, r packd.File) error {
		p, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		f = append(f, big.Field{
			Key:   cleanpath(path),
			Value: string(p),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return f.Struct(), nil
}
