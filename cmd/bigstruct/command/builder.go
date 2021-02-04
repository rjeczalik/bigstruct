package command

import (
	"path"

	"github.com/rjeczalik/bigstruct/isr"
	"github.com/rjeczalik/bigstruct/isr/codec"
	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/spf13/cobra"
)

type Builder struct {
	Import string
	Prefix string
	Values []string
	Codec  isr.Codec
}

func (b *Builder) Register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.StringVarP(&b.Import, "import", "i", "", "")
	f.StringVarP(&b.Prefix, "prefix", "p", "/", "")
	f.StringArrayVarP(&b.Values, "value", "v", nil, "")
}

func (b *Builder) Build() (isr.Fields, error) {
	var (
		f, fields isr.Fields
		err       error
	)

	if f, err = b.buildFromImport(); err != nil {
		return nil, err
	}

	fields = append(fields, f...)

	if f, err = b.buildFromValues(); err != nil {
		return nil, err
	}

	fields = append(fields, f...)

	return fields, nil
}

func (b *Builder) buildFromImport() (isr.Fields, error) {
	if b.Import == "" {
		return nil, nil
	}

	obj, err := isr.MakeFile(b.Import)
	if err != nil {
		return nil, err
	}

	if err := obj.Decode(b.codec()); err != nil {
		return nil, err
	}

	if b.Prefix != "" {
		obj = isr.Move(b.Prefix, obj)
	}

	return obj.Fields(), nil
}

func (b *Builder) buildFromValues() (isr.Fields, error) {
	if len(b.Values) == 0 {
		return nil, nil
	}

	var (
		f  isr.Fields
		kv = types.MakeKV(b.Values...)
	)

	for _, k := range kv.ReverseKeys() {
		var (
			key               = path.Join("/", b.Prefix, k)
			value interface{} = isr.NoValue
		)

		if v := kv[k]; v != "" {
			value = types.YAML(v).Value()
		}

		f = append(f, isr.Field{
			Key:   key,
			Value: value,
		})
	}

	return f, nil
}

func (b *Builder) codec() isr.Codec {
	if b.Codec != nil {
		return b.Codec
	}
	return codec.Default
}
