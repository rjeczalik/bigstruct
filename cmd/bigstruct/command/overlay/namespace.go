package overlay

import (
	"os"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewCommand(app *command.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "overlay",
		Aliases: []string{"ov"},
		Short:   "Manages overlays",
		Args:    cobra.NoArgs,
		RunE:    command.PrintHelp(os.Stderr),
	}

	cmd.AddCommand(
		NewSetCommand(app),
		NewDeleteCommand(app),
		NewListCommand(app),
	)

	return cmd
}
