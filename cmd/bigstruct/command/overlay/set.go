package overlay

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewSetCommand(app *command.App) *cobra.Command {
	m := &setCmd{
		App:     app,
		Meta:    new(command.Meta),
		Overlay: new(model.Overlay),
	}

	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Creates or updates an overlay",
		Args:         cobra.NoArgs,
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type setCmd struct {
	*command.App
	*command.Meta
	*model.Overlay
}

func (m *setCmd) register(cmd *cobra.Command) {
	m.Meta.Register(cmd)

	f := cmd.Flags()

	f.StringVarP(&m.Overlay.Name, "name", "n", "", "")
	f.IntVarP(&m.Overlay.Priority, "priority", "p", 0, "")

	cmd.MarkFlagRequired("name")
}

func (m *setCmd) run(*cobra.Command, []string) error {
	m.Overlay.Metadata = m.Meta.Object()

	if err := m.Storage.UpsertOverlay(m.Overlay); err != nil {
		return err
	}

	return m.Render(model.Overlays{m.Overlay})
}
