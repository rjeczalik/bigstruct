package query

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/isr"

	"github.com/spf13/cobra"
)

func NewSetCommand(app *command.App) *cobra.Command {
	m := &setCmd{
		App:     app,
		Builder: new(command.Builder),
		Printer: new(command.Printer),
	}

	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Creates or updates a struct",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type setCmd struct {
	*command.App
	*command.Builder
	*command.Printer
	namespace string
	index     string
	schema    bool
}

func (m *setCmd) register(cmd *cobra.Command) {
	m.Builder.Register(cmd)

	f := cmd.Flags()

	f.StringVarP(&m.index, "index", "z", "", "")
	f.StringVarP(&m.namespace, "namespace", "N", m.index, "")
	f.BoolVarP(&m.schema, "schema", "x", false, "")

	cmd.MarkFlagRequired("index")
}

func (m *setCmd) setDefaults(cmd *cobra.Command) {
	if !cmd.Flags().Changed("namespace") {
		m.namespace = m.index
	}
}

func (m *setCmd) run(cmd *cobra.Command, _ []string) error {
	m.setDefaults(cmd)

	f, err := m.Builder.Build()
	if err != nil {
		return err
	}

	var (
		obj = f.Object()
	)

	if m.schema {
		obj = obj.Schema()
	} else {
		obj = obj.Raw()
	}

	if err := m.Query.Set(m.Context, m.index, m.namespace, obj); err != nil {
		return err
	}

	return m.Printer.Print(m.App, cmd, obj, isr.Prefix)
}
