package query

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewCommands(app *command.App) []*cobra.Command {
	return []*cobra.Command{
		NewGetCommand(app),
		NewSetCommand(app),
	}
}
