package schema

import (
	"fmt"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewListCommand(app *command.App) *cobra.Command {
	m := &listCmd{
		App: app,
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
	overlay string
	prefix  string
}

func (m *listCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.StringVarP(&m.overlay, "overlay", "L", "", "")
	f.StringVarP(&m.prefix, "prefix", "p", "/", "")

	cmd.MarkFlagRequired("overlay")
}

func (m *listCmd) run(*cobra.Command, []string) error {
	ns, err := m.Storage.Overlay(m.overlay)
	if err != nil {
		return fmt.Errorf("error loading %q overlay: %w", m.overlay, err)
	}

	s, err := m.Storage.ListSchemas(ns, m.prefix)
	if err != nil {
		return fmt.Errorf("error listing %q values for %q overlay: %w", m.prefix, ns.Ref(), err)
	}

	return m.Render(s)
}
