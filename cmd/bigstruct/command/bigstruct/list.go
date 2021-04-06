package bigstruct

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewListCommand(app *command.App) *cobra.Command {
	m := &listCmd{
		App: app,
	}

	cmd := &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all structs",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type listCmd struct {
	*command.App
	index command.Ref
}

func (m *listCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.VarP(&m.index, "index", "z", "")

	cmd.MarkFlagRequired("index")
}

func (m *listCmd) run(cmd *cobra.Command, _ []string) error {
	s, err := m.App.Client.List(m.Context, m.index.Ref())
	if err != nil {
		return err
	}

	list := model.MakeSchemas(new(model.Overlay), s.Fields())

	return m.Render(list)
}
