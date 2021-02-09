package schema

import (
	"os"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/schema/jsonschema"

	"github.com/spf13/cobra"
)

func NewCommand(app *command.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "schema",
		Aliases: []string{"x"},
		Short:   "Manages the schema for structs",
		Args:    cobra.NoArgs,
		RunE:    command.PrintHelp(os.Stderr),
	}

	cmd.AddCommand(
		NewSetCommand(app),
		NewDeleteCommand(app),
		NewListCommand(app),
		jsonschema.NewCommand(app),
	)

	return cmd
}
