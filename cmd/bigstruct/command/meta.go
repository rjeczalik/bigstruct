package command

import (
	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

type Meta []string

func (m *Meta) Register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.StringSliceVarP((*[]string)(m), "meta", "m", nil, "List of key=value pairs of metadata")
}

func (m Meta) Object() model.Object {
	return new(model.Object).Set(types.MakeObject(m...))
}

func (m Meta) Strings() []string {
	return m
}
