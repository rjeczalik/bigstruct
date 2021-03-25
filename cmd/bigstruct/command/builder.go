package command

import (
	"context"
	"fmt"
	"path"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/big/bigutil"
	"github.com/rjeczalik/bigstruct/big/codec"
	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/spf13/cobra"
)

type Builder struct {
	Import string
	Prefix string
	Values []string
	Types  []string
	Codec  big.Codec
}

func (b *Builder) Register(cmd *cobra.Command) {
	f := cmd.Flags()

	f.StringVarP(&b.Import, "import", "i", "", "")
	f.StringVarP(&b.Prefix, "prefix", "p", "/", "")
	f.StringArrayVarP(&b.Values, "value", "v", nil, "")
	f.StringArrayVarP(&b.Types, "type", "t", nil, "")
}

func (b *Builder) Build(ctx context.Context) (big.Fields, error) {
	var (
		f, fields big.Fields
		err       error
	)

	if f, err = b.buildFromImport(ctx); err != nil {
		return nil, err
	}

	fields = append(fields, f...)

	if f, err = b.buildFromValues(); err != nil {
		return nil, err
	}

	fields = append(fields, f...)

	return fields, nil
}

func (b *Builder) buildFromImport(ctx context.Context) (big.Fields, error) {
	if b.Import == "" {
		return nil, nil
	}

	obj, err := bigutil.MakeFile(b.Import)
	if err != nil {
		return nil, err
	}

	if err := obj.Decode(ctx, b.codec()); err != nil {
		return nil, err
	}

	if b.Prefix != "" {
		obj = big.Move(b.Prefix, obj)
	}

	return obj.Fields(), nil
}

func (b *Builder) buildFromValues() (big.Fields, error) {
	if len(b.Values) == 0 && len(b.Types) == 0 {
		return nil, nil
	}

	var (
		f  big.Fields
		kv = types.MakeKV(b.Values...)
		kt = types.MakeKV(b.Types...)
	)

	for _, k := range kv.ReverseKeys() {
		var (
			key               = path.Join("/", b.Prefix, k)
			value interface{} = big.NoValue
			err   error
		)

		if v := kv[k]; v != "" {
			if value, err = types.YAML(v).TryValue(); err != nil {
				value = v
			}
		}

		f = append(f, big.Field{
			Key:   key,
			Value: value,
		})
	}

	for _, k := range kt.ReverseKeys() {
		var (
			key = path.Join("/", b.Prefix, k)
			typ = kt[k]
		)

		if typ == "" {
			return nil, fmt.Errorf("type information for %q key is empty", key)
		}

		f = append(f, big.Field{
			Key:  key,
			Type: typ,
		})
	}

	return f, nil
}

func (b *Builder) codec() big.Codec {
	if b.Codec != nil {
		return b.Codec
	}
	return codec.Default
}
