package bigstruct

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewCommands(app *command.App) []*cobra.Command {
	return []*cobra.Command{
		NewListCommand(app),
		NewGetCommand(app),
		NewSetCommand(app),
		NewDebugCommand(app),
		NewHistoryCommand(app),
	}
}
