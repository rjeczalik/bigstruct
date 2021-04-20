package schema

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewSetCommand(app *command.App) *cobra.Command {
	m := &setCmd{
		App:     app,
		Builder: new(command.Builder),
		Meta:    new(command.Meta),
	}

	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Creates or updates a schema",
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
	*command.Meta
	overlay string
}

func (m *setCmd) register(cmd *cobra.Command) {
	m.Builder.Register(cmd)
	m.Meta.Register(cmd)

	f := cmd.Flags()

	f.StringVarP(&m.overlay, "overlay", "L", "", "")

	cmd.MarkFlagRequired("overlay")
}

func (m *setCmd) run(cmd *cobra.Command, args []string) error {
	return m.Storage.Transaction(m.txRun)
}

func (m *setCmd) txRun(g storage.Gorm) error {
	f, err := m.Builder.Build(m.Context)
	if err != nil {
		return err
	}

	ns, err := g.Overlay(m.overlay)
	if err != nil {
		return err
	}

	s := model.MakeSchemas(ns, f)
	s.SetMeta(m.Meta.Object())

	if err := g.UpsertSchemas(s); err != nil {
		return err
	}

	return m.Render(s)
}
