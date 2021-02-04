package query

import (
	"errors"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewSetCommand(app *command.App) *cobra.Command {
	m := &setCmd{App: app}

	cmd := &cobra.Command{
		Use:          "set",
		Short:        "(wip)",
		Args:         cobra.ExactArgs(1),
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type setCmd struct {
	*command.App
}

func (m *setCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	_ = f
}

func (m *setCmd) run(cmd *cobra.Command, args []string) error {
	return errors.New("not implemented")
}
