package command

import (
	"github.com/rjeczalik/bigstruct"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/pflag"
)

type Ref bigstruct.Ref

var (
	_ pflag.Value = (*Ref)(nil)
)

func (r Ref) Ref() bigstruct.Ref {
	return bigstruct.Ref(r)
}

func (r Ref) String() string {
	return r.Ref().String()
}

func (r Ref) Type() string {
	return "string"
}

func (r *Ref) Set(ref string) error {
	name, prop, err := model.ParseRef(ref)
	if err != nil {
		return err
	}

	r.Name = name
	r.Prop = prop

	return nil
}
