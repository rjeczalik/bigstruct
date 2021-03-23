package xq

import (
	"github.com/rjeczalik/bigstruct/big/bigutil"
	"github.com/rjeczalik/bigstruct/big/codec"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewCommand(app *command.App) *cobra.Command {
	m := &xqCmd{App: app}

	cmd := &cobra.Command{
		Use:           "xq",
		Short:         "",
		Args:          cobra.ExactArgs(1),
		RunE:          m.run,
		SilenceUsage:  true,
		SilenceErrors: true,
		Hidden:        true,
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
	s, err := bigutil.MakeFile(args[0])
	if err != nil {
		return err
	}

	if err := s.Decode(m.Context, codec.Default); err != nil {
		return err
	}

	return m.Render(s)
}
