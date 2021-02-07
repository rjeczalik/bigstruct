package command

import (
	"fmt"
	"os"
	"path"

	"github.com/rjeczalik/bigstruct/isr"
	"github.com/rjeczalik/bigstruct/isr/codec"

	"github.com/spf13/cobra"
)

type Fielder interface {
	Fields() isr.Fields
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
	f.BoolVarP(&p.SchemaOnly, "schema-only", "x", p.SchemaOnly, "")
}

func (p *Printer) Print(app *App, cmd *cobra.Command, f Fielder, prefix string) error {
	if !p.Encode && !p.Decode {
		return app.Render(f)
	}

	obj := f.Fields().Merge()

	if p.Encode {
		if err := obj.Encode(codec.Default); err != nil {
			return err
		}
	} else {
		if err := obj.Decode(codec.Default); err != nil {
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
		return obj.Walk(func(key string, o isr.Object) error {
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
