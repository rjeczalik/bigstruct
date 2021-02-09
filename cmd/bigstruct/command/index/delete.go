package index

import (
	"errors"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewDeleteCommand(app *command.App) *cobra.Command {
	m := &deleteCmd{
		App:   app,
		Index: new(model.Index),
	}

	cmd := &cobra.Command{
		Use:          "delete",
		Aliases:      []string{"del"},
		Short:        "Deletes a static index",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type deleteCmd struct {
	*command.App
	*model.Index
}

func (m *deleteCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.Uint64Var(&m.Index.ID, "id", 0, "")
	f.StringVarP(&m.Index.Name, "name", "n", "", "")
	f.StringVarP(&m.Index.Property, "property", "p", "", "")
}

func (m *deleteCmd) run(_ *cobra.Command, args []string) error {
	if i := m.Index; i.ID == 0 || (i.Name == "" || i.Property == "") {
		return errors.New("either --id or --name and --property flags are required")
	}

	if err := m.Storage.DeleteIndex(m.Index); err != nil {
		return err
	}

	return nil
}
