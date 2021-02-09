package schema

import (
	"fmt"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewListCommand(app *command.App) *cobra.Command {
	m := &listCmd{
		App:     app,
		Printer: new(command.Printer),
	}

	cmd := &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "Describes a schema",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type listCmd struct {
	*command.App
	*command.Printer
	namespace string
	prefix    string
}

func (m *listCmd) register(cmd *cobra.Command) {
	m.Printer.Register(cmd)

	f := cmd.Flags()

	f.StringVarP(&m.namespace, "namespace", "N", "", "")
	f.StringVarP(&m.prefix, "prefix", "p", "/", "")

	cmd.MarkFlagRequired("namespace")
}

func (m *listCmd) run(cmd *cobra.Command, args []string) error {
	ns, err := m.Storage.Namespace(m.namespace)
	if err != nil {
		return fmt.Errorf("error loading %q namespace: %w", m.namespace, err)
	}

	s, err := m.Storage.ListSchemas(ns, m.prefix)
	if err != nil {
		return fmt.Errorf("error listing %q values for %q namespace: %w", m.prefix, ns.Namespace(), err)
	}

	return m.Printer.Print(m.App, cmd, s, m.prefix)
}
