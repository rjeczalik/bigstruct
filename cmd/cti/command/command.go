package command

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

func New(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cti",
		Short: "Command-line interface for working with cti.Object",
		Args:  cobra.NoArgs,
		RunE:  PrintHelp(os.Stderr),
	}

	cmd.AddCommand(
		NewRead(app),
	)

	app.Register(cmd.PersistentFlags())

	return cmd
}

type CobraFunc func(*cobra.Command, []string) error

func PrintHelp(w io.Writer) CobraFunc {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOutput(w)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
