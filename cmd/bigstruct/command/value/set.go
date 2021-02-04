package value

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewSetCommand(app *command.App) *cobra.Command {
	m := &setCmd{
		App:     app,
		Builder: new(command.Builder),
	}

	cmd := &cobra.Command{
		Use:          "set",
		Aliases:      []string{"up"},
		Short:        "(wip)",
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
	namespace string
	schema    bool
}

func (m *setCmd) register(cmd *cobra.Command) {
	m.Builder.Register(cmd)

	f := cmd.Flags()

	f.StringVarP(&m.namespace, "namespace", "N", "", "")
	f.BoolVarP(&m.schema, "schema", "x", false, "")

	cmd.MarkFlagRequired("namespace")
}

func (m *setCmd) run(cmd *cobra.Command, args []string) error {
	f, err := m.Builder.Build()
	if err != nil {
		return err
	}

	ns, err := m.Storage.Namespace(m.namespace)
	if err != nil {
		return err
	}

	tx := m.Storage.DB.Begin()

	if m.schema {
		s := model.MakeSchemas(ns, f)

		if err := m.Storage.With(tx).UpsertSchemas(s); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	v := model.MakeValues(ns, f)

	if err := m.Storage.With(tx).UpsertValues(v); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return m.Render(v)
}
