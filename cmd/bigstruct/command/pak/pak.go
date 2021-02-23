package pak

import (
	"os"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewCommand(app *command.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pak",
		Short: "Manages bigpacks",
		Args:  cobra.NoArgs,
		RunE:  command.PrintHelp(os.Stderr),
	}

	cmd.AddCommand(
		NewImportCommand(app),
	)

	return cmd
}
