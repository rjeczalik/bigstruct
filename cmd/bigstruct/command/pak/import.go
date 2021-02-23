package pak

import (
	"fmt"

	"github.com/rjeczalik/bigstruct/bigpack"
	"github.com/rjeczalik/bigstruct/bigpack/pak"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewImportCommand(app *command.App) *cobra.Command {
	m := &importCmd{App: app}

	cmd := &cobra.Command{
		Use:          "import",
		Aliases:      []string{"im"},
		Short:        "Imports a bigpack",
		Args:         cobra.ExactArgs(1),
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type importCmd struct {
	*command.App
	dry bool
}

func (m *importCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.BoolVarP(&m.dry, "dry", "d", false, "")
}

func (m *importCmd) run(cmd *cobra.Command, args []string) error {
	m.setDefaults(cmd)

	pk, err := bigpack.Read(pak.Dir(args[0]))
	if err != nil {
		return fmt.Errorf("failed to read %q buildpack: %w", args[0], err)
	}

	if m.dry {
		return m.Render(pk)
	}

	if err := m.Storage.Transaction(pk.Store); err != nil {
		return fmt.Errorf("failed to store %q buildpack: %w", args[0], err)
	}

	return nil
}

func (m *importCmd) setDefaults(cmd *cobra.Command) {
	f := cmd.Flags()

	if !f.Changed("format") {
		m.Format = "yaml"
	}
}
