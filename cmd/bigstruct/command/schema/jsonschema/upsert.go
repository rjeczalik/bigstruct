package jsonschema

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewUpsertCommand(app *command.App) *cobra.Command {
	m := &upsertCmd{
		App:     app,
		Builder: new(command.Builder),
	}

	cmd := &cobra.Command{
		Use:          "upsert",
		Aliases:      []string{"up"},
		Short:        "(wip)",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type upsertCmd struct {
	*command.App
	*command.Builder
	namespace string
}

func (m *upsertCmd) register(cmd *cobra.Command) {
	m.Builder.Register(cmd)

	f := cmd.Flags()

	f.StringVarP(&m.namespace, "namespace", "N", "", "")

	cmd.MarkFlagRequired("namespace")
}

func (m *upsertCmd) run(cmd *cobra.Command, args []string) error {
	f, err := m.Builder.Build(m.Context)
	if err != nil {
		return err
	}

	ns, err := m.Storage.Namespace(m.namespace)
	if err != nil {
		return err
	}

	s := model.MakeSchemas(ns, f)

	if err := m.Storage.UpsertSchemas(s); err != nil {
		return err
	}

	return m.Render(s)
}
