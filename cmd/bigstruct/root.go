package main

import (
	"os"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/bigstruct"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/index"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/overlay"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/pak"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/schema"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/value"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/xq"

	"github.com/spf13/cobra"
)

func NewCommand(app *command.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bigstruct",
		Short: "Command-line interface to bigstruct storage",
		Args:  cobra.NoArgs,
		RunE:  command.PrintHelp(os.Stderr),
	}

	cmd.AddCommand(
		overlay.NewCommand(app),
		index.NewCommand(app),
		value.NewCommand(app),
		schema.NewCommand(app),
		pak.NewCommand(app),
		command.Hidden(xq.NewCommand(app)),
	)

	cmd.AddCommand(
		bigstruct.NewCommands(app)...,
	)

	app.Register(cmd.PersistentFlags())

	cmd.PersistentPreRunE = app.Init

	return cmd
}
