package bigstruct

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewSetCommand(app *command.App) *cobra.Command {
	m := &setCmd{
		App:     app,
		Builder: new(command.Builder),
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
	namespace command.Ref
	index     command.Ref
	schema    bool
}

func (m *setCmd) register(cmd *cobra.Command) {
	m.Builder.Register(cmd)

	f := cmd.Flags()

	f.VarP(&m.index, "index", "z", "")
	f.VarP(&m.namespace, "namespace", "N", "")

	cmd.MarkFlagRequired("index")
}

func (m *setCmd) setDefaults(cmd *cobra.Command) {
	if !cmd.Flags().Changed("namespace") {
		m.namespace = m.index
	}
}

func (m *setCmd) run(cmd *cobra.Command, _ []string) error {
	m.setDefaults(cmd)

	f, err := m.Builder.Build(m.Context)
	if err != nil {
		return err
	}

	var (
		obj = f.Struct()
	)

	if err := m.Client.Set(m.Context, m.index.Ref(), m.namespace.Ref(), obj); err != nil {
		return err
	}

	return m.Render(f)
}
