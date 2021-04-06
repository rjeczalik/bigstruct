package model

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/big"
)

type Schema struct {
	Model           `yaml:",inline"`
	Overlay         *Overlay `yaml:"-" json:"-"`
	OverlayID       uint64   `gorm:"column:overlay_id;type:bigint;not null;index" yaml:"overlay_id,omitempty" json:"overlay_id,omitempty"`
	OverlayProperty string   `gorm:"column:overlay_property;type:tinytext;not null" yaml:"overlay_property,omitempty" json:"overlay_property,omitempty"`
	Key             string   `gorm:"column:key;type:text;not null" yaml:"key,omitempty" json:"key,omitempty"`
	Type            string   `gorm:"column:type;type:tinytext;not null" yaml:"type,omitempty" json:"type,omitempty"`
	Encoding        string   `gorm:"column:encoding;type:tinytext;not null" yaml:"encoding,omitempty" json:"encoding,omitempty"`
	Metadata        Object   `gorm:"column:metadata;type:text" yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

func (*Schema) TableName() string {
	return TablePrefix + "_schema"
}

func (s *Schema) Codec() string {
	if s.Encoding != "" {
		return path.Join(s.Type, s.Encoding)
	}
	return s.Type
}

type Schemas []*Schema

func MakeSchemas(o *Overlay, f big.Fields) Schemas {
	schemas := make(Schemas, 0, len(f))

	for _, f := range f {
		if f.Type == "" {
			continue // skip empty entries, they will get recreated from the tree either way
		}

		s := &Schema{
			Key:  f.Key,
			Type: f.Type,
		}

		if o != nil {
			s.Overlay = o
			s.OverlayID = o.ID
			s.OverlayProperty = o.Property
		}

		if i := strings.IndexRune(f.Type, '/'); i != -1 {
			s.Type = f.Type[:i]
			s.Encoding = f.Type[i+1:]
		}

		schemas = append(schemas, s)
	}

	return schemas
}

func (s Schemas) SetOverlay(o *Overlay) {
	for _, s := range s {
		s.Overlay = o
		s.OverlayID = o.ID
		s.OverlayProperty = o.Property
	}
}

func (s Schemas) SetMeta(meta Object) {
	for _, s := range s {
		s.Metadata = meta
	}
}

func (s Schemas) Fields() big.Fields {
	f := make(big.Fields, 0, len(s))

	for _, s := range s {
		f = append(f, big.Field{
			Key:  s.Key,
			Type: path.Join(s.Type, s.Encoding),
		})
	}

	return f
}

func (s Schemas) WriteTab(w io.Writer) (int64, error) {
	var n int64

	m, err := fmt.Fprint(w, "ID\tOVERLAY\tKEY\tTYPE\tENCODING\tMETADATA\n")
	if err != nil {
		return int64(m), err
	}

	n += int64(m)

	for _, s := range s {
		m, err = fmt.Fprintf(w, "%v\t%s\t%s\t%s\t%s\t%s\n",
			s.ID,
			nonempty(s.Overlay.Ref(), "-"),
			s.Key,
			s.Type,
			nonempty(s.Encoding, "-"),
			nonempty(s.Metadata.String(), "-"),
		)

		n += int64(m)

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (s Schemas) WriteTo(w io.Writer) (int64, error) {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	n, err := s.WriteTab(tw)
	if err != nil {
		return n, err
	}

	if err := tw.Flush(); err != nil {
		return n, err
	}

	return n, err

}

func (s Schemas) String() string {
	var buf bytes.Buffer

	if _, err := s.WriteTo(&buf); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return buf.String()
}
