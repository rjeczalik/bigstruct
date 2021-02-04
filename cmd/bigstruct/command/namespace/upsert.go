package namespace

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewSetCommand(app *command.App) *cobra.Command {
	m := &setCmd{
		App:       app,
		Namespace: new(model.Namespace),
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
	*model.Namespace
}

func (m *setCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.StringVarP(&m.Namespace.Name, "name", "n", "", "")
	f.IntVarP(&m.Namespace.Priority, "priority", "p", 0, "")
	f.StringVarP((*string)(&m.Namespace.Property), "property", "P", "true", "")

	cmd.MarkFlagRequired("name")
}

func (m *setCmd) run(*cobra.Command, []string) error {
	if err := m.Storage.UpsertNamespace(m.Namespace); err != nil {
		return err
	}

	return m.Render(model.Namespaces{m.Namespace})
}
