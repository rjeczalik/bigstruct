package value

import (
	"os"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewCommand(app *command.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "value",
		Aliases: []string{"v"},
		Short:   "Manages values",
		Args:    cobra.NoArgs,
		RunE:    command.PrintHelp(os.Stderr),
		Hidden:  true,
	}

	cmd.AddCommand(
		NewSetCommand(app),
		NewDeleteCommand(app),
		NewListCommand(app),
	)

	return cmd
}
