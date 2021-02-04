package schema

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
}

func (m *setCmd) register(cmd *cobra.Command) {
	m.Builder.Register(cmd)

	f := cmd.Flags()

	f.StringVarP(&m.namespace, "namespace", "N", "", "")

	cmd.MarkFlagRequired("namespace")
}

func (m *setCmd) run(cmd *cobra.Command, args []string) error {
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

	s := model.MakeSchemas(ns, f)

	if err := g.UpsertSchemas(s); err != nil {
		return err
	}

	return m.Render(s)
}
