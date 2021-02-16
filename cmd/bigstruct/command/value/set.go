package value

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewSetCommand(app *command.App) *cobra.Command {
	m := &setCmd{
		App:     app,
		Builder: new(command.Builder),
		Meta:    new(command.Meta),
	}

	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Creates or updates values",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type setCmd struct {
	*command.App
	*command.Builder
	*command.Meta
	namespace string
	schema    bool
}

func (m *setCmd) register(cmd *cobra.Command) {
	m.Builder.Register(cmd)
	m.Meta.Register(cmd)

	f := cmd.Flags()

	f.StringVarP(&m.namespace, "namespace", "N", "", "")
	f.BoolVarP(&m.schema, "schema", "x", false, "")

	cmd.MarkFlagRequired("namespace")
}

func (m *setCmd) run(*cobra.Command, []string) error {
	return m.Storage.Transaction(m.txRun)
}

func (m *setCmd) txRun(g storage.Gorm) error {
	f, err := m.Builder.Build()
	if err != nil {
		return err
	}

	ns, err := g.Namespace(m.namespace)
	if err != nil {
		return err
	}

	if m.schema {
		s := model.MakeSchemas(ns, f)

		if err := g.UpsertSchemas(s); err != nil {
			return err
		}
	}

	v := model.MakeValues(ns, f)
	v.SetMeta(m.Meta.Object())

	if err := g.UpsertValues(v); err != nil {
		return err
	}

	return m.Render(v)
}
