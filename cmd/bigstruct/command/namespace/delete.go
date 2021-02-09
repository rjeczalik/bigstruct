package namespace

import (
	"errors"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewDeleteCommand(app *command.App) *cobra.Command {
	m := &deleteCmd{
		App:       app,
		Namespace: new(model.Namespace),
	}

	cmd := &cobra.Command{
		Use:          "delete",
		Aliases:      []string{"del"},
		Short:        "Deletes a namespace",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type deleteCmd struct {
	*command.App
	*model.Namespace
}

func (m *deleteCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.Uint64Var(&m.Namespace.ID, "id", 0, "")
	f.StringVarP(&m.Namespace.Name, "name", "n", "", "")
}

func (m *deleteCmd) run(*cobra.Command, []string) error {
	if m.Namespace.ID == 0 && m.Namespace.Name == "" {
		return errors.New("either --id or --name flag is required")
	}

	if err := m.Storage.DeleteNamespace(m.Namespace); err != nil {
		return err
	}

	return nil
}
