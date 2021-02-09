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
	*model.Index
	values  []string
	schemas []string
}

func (m *setCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.StringVarP(&m.Index.Name, "name", "n", "", "")
	f.StringVarP(&m.Index.Property, "property", "p", "", "")
	f.StringSliceVarP(&m.values, "value", "v", nil, "")
	f.StringSliceVarP(&m.schemas, "schema", "x", nil, "")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("property")
}

func (m *setCmd) run(*cobra.Command, []string) error {
	if kv := types.MakeKV(m.values...); len(kv) != 0 {
		m.Index.ValueIndex.Set(kv)
	}

	if kv := types.MakeKV(m.schemas...); len(kv) != 0 {
		m.Index.SchemaIndex.Set(kv)
	}

	if err := m.Storage.UpsertIndex(m.Index); err != nil {
		return err
	}

	return m.Render(model.Indexes{m.Index})
}
