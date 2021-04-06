package overlay

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewListCommand(app *command.App) *cobra.Command {
	m := &listCmd{App: app}

	cmd := &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "Lists all overlays",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type listCmd struct {
	*command.App
}

func (m *listCmd) register(cmd *cobra.Command) {
	_ = cmd.Flags()
}

func (m *listCmd) run(*cobra.Command, []string) error {
	n, err := m.Storage.ListOverlays()
	if err != nil {
		return err
	}

	return m.Render(n)
}
