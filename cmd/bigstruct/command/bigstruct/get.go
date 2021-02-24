package bigstruct

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewGetCommand(app *command.App) *cobra.Command {
	m := &getCmd{
		App:     app,
		Printer: new(command.Printer),
	}

	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Reads a struct",
		Args:         cobra.ExactArgs(1),
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type getCmd struct {
	*command.App
	*command.Printer
	index command.Ref
}

func (m *getCmd) register(cmd *cobra.Command) {
	m.Printer.Register(cmd)

	f := cmd.Flags()

	f.VarP(&m.index, "index", "z", "")

	cmd.MarkFlagRequired("index")
}

func (m *getCmd) run(cmd *cobra.Command, args []string) error {
	s, err := m.Client.Get(m.Context, m.index.Ref(), args[0])
	if err != nil {
		return err
	}

	return m.Printer.Print(m.App, cmd, s, args[0])
}
