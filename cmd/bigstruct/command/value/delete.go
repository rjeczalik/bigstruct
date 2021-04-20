package value

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
		Short:        "Deletes a value",
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

	_ = f
}

func (m *deleteCmd) run(*cobra.Command, []string) error {
	return errors.New("not implemented")
}
