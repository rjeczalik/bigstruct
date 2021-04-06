package bigstruct

import (
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewDiffCommand(app *command.App) *cobra.Command {
	m := &diffCmd{
		App: app,
	}

	cmd := &cobra.Command{
		Use:          "diff",
		Short:        "Reads a struct",
		Args:         cobra.ExactArgs(1),
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type diffCmd struct {
	*command.App
	index      command.Ref
	start, end command.Ref
}

func (m *diffCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.VarP(&m.index, "index", "z", "")
	f.VarP(&m.start, "start", "s", "")
	f.VarP(&m.end, "end", "e", "")

	cmd.MarkFlagRequired("index")
	cmd.MarkFlagRequired("start")
	cmd.MarkFlagRequired("end")
}

func (m *diffCmd) run(cmd *cobra.Command, args []string) error {
	s, err := m.Client.Diff(m.Context, m.index.Ref(), m.start.Ref(), m.end.Ref(), args[0])
	if err != nil {
		return err
	}

	return m.Render(s)
}
