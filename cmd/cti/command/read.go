package command

import (
	"fmt"
	"os"
	"regexp"

	"github.com/glaucusio/confetti/cti"
	_ "github.com/glaucusio/confetti/cti/codec"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewRead(app *App) *cobra.Command {
	m := &readCmd{App: app}

	cmd := &cobra.Command{
		Use:   "read",
		Short: "Read a file / directory",
		Args:  cobra.MaximumNArgs(2),
		RunE:  m.run,
	}

	m.register(cmd.Flags())

	return cmd
}

func (m *readCmd) register(f *pflag.FlagSet) {
	f.StringVarP(&m.format, "format", "f", "", "Target encoding format to use")
}

type readCmd struct {
	*App
	format string
}

func (m *readCmd) run(_ *cobra.Command, args []string) error {
	var (
		name   = "."
		filter = ".*"
	)

	if len(args) > 0 {
		name = args[0]
	}

	if len(args) > 1 {
		filter = args[1]
	}

	r, err := regexp.Compile(filter)
	if err != nil {
		return err
	}

	obj, err := cti.MakeFile(name)
	if err != nil {
		return err
	}

	if err := obj.Decode(nil); err != nil {
		return err
	}

	var fields cti.Fields

	for _, f := range obj.Fields() {
		if r.MatchString(f.Key) {
			fields = append(fields, f)
		}
	}

	obj = fields.Object()

	if m.format == "" {
		obj.WriteTo(os.Stdout)
		return nil
	}

	envelope := cti.Object{
		"root": {
			Encoding: m.format,
			Children: obj,
		},
	}

	if err := envelope.Encode(nil); err != nil {
		return err
	}

	switch v := envelope["root"].Value.(type) {
	case []byte:
		os.Stdout.Write(v)
	default:
		fmt.Printf("%s\n", v)
	}

	return nil
}
