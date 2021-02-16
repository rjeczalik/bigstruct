package index

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewSetCommand(app *command.App) *cobra.Command {
	m := &setCmd{
		App:   app,
		Meta:  new(command.Meta),
		Index: new(model.Index),
	}

	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Creates or updates static indexes",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type setCmd struct {
	*command.App
	*command.Meta
	*model.Index
	values []string
}

func (m *setCmd) register(cmd *cobra.Command) {
	m.Meta.Register(cmd)

	f := cmd.Flags()

	f.StringVarP(&m.Index.Name, "name", "n", "", "")
	f.StringVarP(&m.Index.Property, "property", "p", "", "")
	f.StringSliceVarP(&m.values, "value", "v", nil, "")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("property")
}

func (m *setCmd) run(*cobra.Command, []string) error {
	m.Index.Metadata = m.Meta.Object()

	if o := types.MakeObject(m.values...); len(o) != 0 {
		m.Index.Index.Set(o)
	}

	if err := m.Storage.UpsertIndex(m.Index); err != nil {
		return err
	}

	return m.Render(model.Indexes{m.Index})
}
