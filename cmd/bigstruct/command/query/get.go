package query

import (
	"fmt"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewGetCommand(app *command.App) *cobra.Command {
	m := &getCmd{
		App: app,
		Printer: &command.Printer{
			Encode: true,
		},
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
	index     string
	namespace string
}

func (m *getCmd) register(cmd *cobra.Command) {
	m.Printer.Register(cmd)

	f := cmd.Flags()

	f.StringVarP(&m.index, "index", "z", "", "")

	cmd.MarkFlagRequired("index")
}

func (m *getCmd) setDefaults(cmd *cobra.Command) {
	if !cmd.Flags().Changed("namespace") {
		m.namespace = m.index
	}
}

func (m *getCmd) run(cmd *cobra.Command, args []string) error {
	m.setDefaults(cmd)

	obj, err := m.Query.Get(m.Context, m.index, m.namespace, args[0])
	if err != nil {
		return err
	}

	if len(obj) == 0 {
		return fmt.Errorf("no values found for %q", args[0])
	}

	return m.Printer.Print(m.App, cmd, obj, args[0])
}
