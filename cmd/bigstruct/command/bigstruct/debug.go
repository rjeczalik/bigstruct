package bigstruct

import (
	"fmt"

	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"
	"github.com/rjeczalik/bigstruct/storage/model"

	"github.com/spf13/cobra"
)

func NewDebugCommand(app *command.App) *cobra.Command {
	m := &debugCmd{
		App: app,
	}

	cmd := &cobra.Command{
		Use:          "debug",
		Short:        "Prints LSM-Tree for debugging",
		Args:         cobra.ExactArgs(1),
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type debugCmd struct {
	*command.App
	index   command.Ref
	overlay command.Ref
}

func (m *debugCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.VarP(&m.index, "index", "z", "")
	f.VarP(&m.overlay, "overlay", "L", "")

	cmd.MarkFlagRequired("index")
}

func (m *debugCmd) run(cmd *cobra.Command, args []string) error {
	index, schema, values, err := m.Client.Debug(m.Context, m.index.Ref(), m.overlay.Ref(), args[0])
	if err != nil {
		return err
	}

	fmt.Printf("# INDEX\n\n%s\n# SCHEMA\n\n%s\n# VALUES\n\n%s", model.Indexes{index}, schema, values)

	return nil
}
