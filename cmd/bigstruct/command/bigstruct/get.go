package bigstruct

import (
	"fmt"
	"path"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/cmd/bigstruct/command"

	"github.com/spf13/cobra"
)

func NewGetCommand(app *command.App) *cobra.Command {
	m := &getCmd{
		App: app,
	}

	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Reads a struct",
		Args:         cobra.ExactArgs(1),
		RunE:         m.run,
		SilenceUsage: true,
	}

	m.register(cmd)

	return cmd
}

type getCmd struct {
	*command.App
	index command.Ref
}

func (m *getCmd) register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.VarP(&m.index, "index", "z", "")

	cmd.MarkFlagRequired("index")
}

func (m *getCmd) run(cmd *cobra.Command, args []string) error {
	s, err := m.Client.Get(m.Context, m.index.Ref(), args[0])
	if err != nil {
		return err
	}

	return s.Walk(func(key string, o big.Struct) error {
		var (
			k = path.Base(key)
			n = o[k]
		)

		if n.Value == nil {
			return nil
		}

		fmt.Printf("# %s\n%s\n", key, n.Value)

		return nil
	})
}
