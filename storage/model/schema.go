package model

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/isr"
	"github.com/rjeczalik/bigstruct/internal/types"
)

type Schema struct {
	Model       `yaml:",inline"`
	Key         string     `gorm:"column:key;type:text;not null;index:idx_schema_key_namespace_id" yaml:"key,omitempty" json:"key,omitempty"`
	Type        string     `gorm:"column:type;type:tinytext;not null" yaml:"type,omitempty" json:"type,omitempty"`
	Encoding    string     `gorm:"column:encoding;type:tinytext;not null" yaml:"encoding,omitempty" json:"encoding,omitempty"`
	RawSchema   string     `gorm:"column:schema;type:text" yaml:"schema,omitempty" json:"schema,omitempty"`
	Namespace   *Namespace `yaml:"-" json:"-"`
	NamespaceID uint64     `gorm:"column:namespace_id;type:bigint;not null;index:idx_schema_key_namespace_id" yaml:"namespace_id,omitempty" json:"namespace_id,omitempty"`
}

func (*Schema) TableName() string {
	return Prefix + "_schema"
}

func (s *Schema) Schema() interface{} {
	return types.JSON(s.RawSchema).Value()
}

func (s *Schema) SetSchema(v interface{}) {
	s.RawSchema = types.MakeJSON(s).String()
}

type Schemas []*Schema

func MakeSchemas(ns *Namespace, f isr.Fields) Schemas {
	values := make(Schemas, 0, len(f))

	for _, f := range f {
		if f.Type == "" {
			continue // skip empty entries, they will get recreated from the tree either way
		}

		v := &Schema{
			Key:         f.Key,
			Type:        f.Type,
			Namespace:   ns,
			NamespaceID: ns.ID,
		}

		if i := strings.IndexRune(f.Type, '/'); i != -1 {
			v.Type = f.Type[:i]
			v.Encoding = f.Type[i+1:]
		}

		values = append(values, v)
	}

	return values
}

func (s Schemas) Fields() isr.Fields {
	f := make(isr.Fields, 0, len(s))

	for _, s := range s {
		f = append(f, isr.Field{
			Key:  s.Key,
			Type: path.Join(s.Type, s.Encoding),
		})
	}

	return f
}

func (s Schemas) SetNamespace(ns *Namespace) {
	for _, s := range s {
		s.Namespace = ns
	}
}

func (s Schemas) WriteTab(w io.Writer) (int64, error) {
	var n int64

	m, err := fmt.Fprint(w, "ID\tNAMESPACE\tKEY\tTYPE\tENCODING\tSCHEMA\n")
	if err != nil {
		return int64(m), err
	}

	n += int64(m)

	for _, s := range s {
		m, err = fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			s.ID,
			s.Namespace.Namespace(),
			s.Key,
			s.Type,
			s.Encoding,
			nonempty(s.RawSchema, "-"),
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
