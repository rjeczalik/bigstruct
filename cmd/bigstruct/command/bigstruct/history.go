package bigstruct

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewHistoryCommand(app *command.App) *cobra.Command {
	m := &historyCmd{
		App: app,
	}

	cmd := &cobra.Command{
		Use:          "history",
		Short:        "Shows modification history",
		Args:         cobra.ExactArgs(1),
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type historyCmd struct {
	*command.App
	index     string
	namespace string
}

func (m *historyCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.StringVarP(&m.index, "index", "z", "", "")

	cmd.MarkFlagRequired("index")
}

func (m *historyCmd) run(cmd *cobra.Command, args []string) error {
	return nil
}
