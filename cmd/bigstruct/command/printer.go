package command

import (
	"fmt"
	"os"
	"path"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/big/codec"

	"github.com/spf13/cobra"
)

type Fielder interface {
	Fields() big.Fields
}

type Printer struct {
	Encode     bool
	Decode     bool
	SchemaOnly bool
	Raw        bool
}

func (p *Printer) Register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.BoolVarP(&p.Encode, "encode", "e", p.Encode, "")
	f.BoolVarP(&p.Decode, "decode", "d", p.Decode, "")
	f.BoolVarP(&p.Raw, "raw", "r", p.Raw, "")
	f.BoolVarP(&p.SchemaOnly, "schema-only", "s", p.SchemaOnly, "")
}

func (p *Printer) Print(app *App, cmd *cobra.Command, f Fielder, prefix string) error {
	obj := f.Fields().Merge()

	if p.Encode {
		if err := obj.Encode(app.Context, codec.Default); err != nil {
			return err
		}
	} else if p.Decode {
		if err := obj.Decode(app.Context, codec.Default); err != nil {
			return err
		}
	}

	switch {
	case p.Raw && !p.SchemaOnly:
		return app.Render(obj.At(prefix))
	case p.Raw:
		return app.Render(obj.At(prefix).Schema())
	case p.Decode:
		app.DefaultFormat(cmd, "yaml")
		return app.Render(obj.At(prefix).Value())
	default:
		return obj.Walk(func(key string, o big.Struct) error {
			var (
				k = path.Base(key)
				n = o[k]
			)

			if n.Value != nil {
				fmt.Fprintf(os.Stderr, "# %s\n", key)
			}

			return app.Render(n.Value)
		})
	}
}
