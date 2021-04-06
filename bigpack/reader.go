package bigpack

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	stdpath "path"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/big/codec"
	"github.com/rjeczalik/bigstruct/bigpack/pak"
	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/storage/model"

	"gopkg.in/yaml.v3"
)

var globalReader = &Reader{}

type Reader struct {
	Codec big.Codec // or codec.Default if nil
}

func (r *Reader) Read(ctx context.Context, fs pak.FS) (*pak.Pak, error) {
	const (
		def   = "bigpack.yaml"
		index = ".bigpack"
		sep   = string(os.PathSeparator)
	)

	var (
		pk  pak.Pak
		err error
		ok  bool
	)

	err = fs(func(path string, rc io.ReadCloser) error {
		if path != def {
			return nil // skip
		}

		ok = true

		p, err := ioutil.ReadAll(rc)
		if err != nil {
			return fmt.Errorf("error reading %q: %w", path, err)
		}

		if err := yaml.Unmarshal(p, &pk); err != nil {
			return fmt.Errorf("error parsing %q: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("did not find definition file %q", def)
	}

	var (
		mschema = make(map[string]big.Fields)
		mvalue  = make(map[string]big.Fields)
	)

	err = fs(func(path string, rc io.ReadCloser) error {
		if path == def || !(strings.HasPrefix(path, "schema"+sep) || strings.HasPrefix(path, "value"+sep)) {
			return nil // skip
		}

		// path ~ <type>/<overlay>[/<property>]/<key>
		parts := strings.SplitN(path, sep, 4)

		out := mschema
		if parts[0] == "value" {
			out = mvalue
		}

		if len(parts) < 3 {
			return fmt.Errorf("file %q has invalid layout", path)
		}

		o := pk.Overlays.ByName(parts[1])
		if o == nil {
			return fmt.Errorf("overlay %q not found for key: %q", parts[1], path)
		}

		var key string

		if !o.Meta().NoProperty {
			if err := o.SetProperty(parts[2]); err != nil {
				return fmt.Errorf(
					"error setting %q property for %q overlay and %q key: %w",
					parts[2], o.Name, path, err,
				)
			}

			key = stdpath.Join("/", filepath.ToSlash(parts[3]))
		} else if len(parts) == 4 {
			key = stdpath.Join("/", parts[2], filepath.ToSlash(parts[3]))
		} else {
			key = stdpath.Join("/", filepath.ToSlash(parts[2]))
		}

		p, err := ioutil.ReadAll(rc)
		if err != nil {
			return fmt.Errorf("error reading %q key from %q overlay: %w", key, o.Ref(), err)
		}

		if stdpath.Base(key) == index {
			var idx map[string]string

			if err := types.YAML(p).Unmarshal(&idx); err != nil {
				return fmt.Errorf("error parsing %q index from %q overlay: %w", key, o.Ref(), err)
			}

			for k, typ := range idx {
				out[o.Ref()] = append(out[o.Ref()], big.Field{
					Key:  stdpath.Join(stdpath.Dir(key), k),
					Type: typ,
				})
			}
		} else {
			out[o.Ref()] = append(out[o.Ref()], big.Field{
				Key:   key,
				Value: string(p),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for ref, f := range mschema {
		s := f.Struct()

		if err := s.Decode(ctx, r.codec()); err != nil {
			return nil, fmt.Errorf("error decoding schema for %q overlay: %w", ref, err)
		}

		fmt.Println("DEBU", s)

		var (
			o      = pk.Overlays.ByRef(ref)
			schema = model.MakeSchemas(o, s.Fields())
		)

		pk.Schemas = append(pk.Schemas, schema...)
	}

	for ref, f := range mvalue {
		s := f.Struct()

		if err := s.Decode(ctx, r.codec()); err != nil {
			return nil, fmt.Errorf("error decoding values for %q overlay: %w", ref, err)
		}

		f = s.Fields()

		var (
			o      = pk.Overlays.ByRef(ref)
			schema = model.MakeSchemas(o, f)
			value  = model.MakeValues(o, f)
		)

		pk.Schemas = append(pk.Schemas, schema...)
		pk.Values = append(pk.Values, value...)
	}

	return &pk, nil
}

func (r *Reader) codec() big.Codec {
	if r.Codec != nil {
		return r.Codec
	}
	return codec.Default
}

func Read(ctx context.Context, fs pak.FS) (*pak.Pak, error) {
	return globalReader.Read(ctx, fs)
}
