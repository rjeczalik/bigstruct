package xq

import (
	"errors"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewCommand(app *command.App) *cobra.Command {
	m := &xqCmd{App: app}

	cmd := &cobra.Command{
		Use:           "xq",
		Short:         "(wip)",
		Args:          cobra.NoArgs,
		RunE:          m.run,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	m.register(cmd)

	return cmd
}

type xqCmd struct {
	*command.App
}

func (m *xqCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	_ = f
}

func (m *xqCmd) run(_ *cobra.Command, args []string) error {
	return errors.New("not implemented")
}
