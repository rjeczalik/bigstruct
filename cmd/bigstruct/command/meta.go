package command

import (
	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/pflag"
)

type Meta []string

func (m *Meta) Register(f *pflag.FlagSet) {
	f.StringSliceVarP((*[]string)(m), "meta", "m", nil, "List of key=value pairs of metadata")
}

func (m Meta) Object() types.Object {
	return types.MakeObject(m...)
}

func (m Meta) Metadata() model.Metadata {
	return new(model.Metadata).Set(m.Object())
}

func (m Meta) Slice() []string {
	return m
}

func (m Meta) JSON() types.JSON {
	return m.Object().JSON()
}
