package overlay

import (
	"errors"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewDeleteCommand(app *command.App) *cobra.Command {
	m := &deleteCmd{
		App:     app,
		Overlay: new(model.Overlay),
	}

	cmd := &cobra.Command{
		Use:          "delete",
		Aliases:      []string{"del"},
		Short:        "Deletes an overlay",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type deleteCmd struct {
	*command.App
	*model.Overlay
}

func (m *deleteCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.Uint64Var(&m.Overlay.ID, "id", 0, "")
	f.StringVarP(&m.Overlay.Name, "name", "n", "", "")
}

func (m *deleteCmd) run(*cobra.Command, []string) error {
	if m.Overlay.ID == 0 && m.Overlay.Name == "" {
		return errors.New("either --id or --name flag is required")
	}

	if err := m.Storage.DeleteOverlay(m.Overlay); err != nil {
		return err
	}

	return nil
}
