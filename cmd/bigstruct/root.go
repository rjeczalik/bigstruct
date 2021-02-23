package main

import (
	"os"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/index"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/namespace"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/pak"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command/query"
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
		namespace.NewCommand(app),
		index.NewCommand(app),
		value.NewCommand(app),
		schema.NewCommand(app),
		pak.NewCommand(app),
		command.Hidden(xq.NewCommand(app)),
	)

	cmd.AddCommand(
		query.NewCommands(app)...,
	)

	app.Register(cmd.PersistentFlags())

	cmd.PersistentPreRunE = app.Init

	return cmd
}
