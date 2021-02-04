package jsonschema

import (
	"os"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewCommand(app *command.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "jsonschema",
		Aliases: []string{"j"},
		Short:   "(wip)",
		Args:    cobra.NoArgs,
		RunE:    command.PrintHelp(os.Stderr),
	}

	cmd.AddCommand(
		NewUpsertCommand(app),
		NewDeleteCommand(app),
		NewListCommand(app),
	)

	return cmd
}
