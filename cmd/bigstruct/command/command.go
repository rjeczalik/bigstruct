package command

import (
	"io"

	"github.com/spf13/cobra"
)

type CobraFunc func(*cobra.Command, []string) error

func Hidden(cmd *cobra.Command) *cobra.Command {
	cmd.Hidden = true
	return cmd
}

func PrintHelp(w io.Writer) CobraFunc {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOutput(w)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
