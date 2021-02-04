package command

import (
	"fmt"
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
	Value      bool
}

func (p *Printer) Register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.BoolVarP(&p.Encode, "encode", "e", false, "")
	f.BoolVarP(&p.Decode, "decode", "d", false, "")
	f.BoolVarP(&p.Raw, "raw", "r", false, "")
	f.BoolVarP(&p.SchemaOnly, "schema-only", "x", false, "")
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
		var (
			dir  = path.Dir(prefix)
			base = path.Base(prefix)
		)

		fmt.Println("DIR:", obj.At(dir))
		_ = base

		return app.Render(fmt.Sprintf("%s", obj.At(path.Dir(prefix))[path.Base(prefix)].Value))
	}
}
