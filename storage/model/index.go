package model

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"
)

type Index struct {
	Model       `yaml:",inline"`
	Name        string   `gorm:"column:name;type:tinytext;not null" yaml:"name,omitempty" json:"name,omitempty"`
	Property    Property `gorm:"column:property;type:tinytext" yaml:"property,omitempty" json:"property,omitempty"`
	ValueIndex  Metadata `gorm:"column:value_index;type:text;not null" yaml:"value_index,omitempty" json:"value_index,omitempty"`
	SchemaIndex Metadata `gorm:"column:schema_index;type:text;not null" yaml:"schema_index,omitempty" json:"schema_index,omitempty"`
}

func (*Index) TableName() string {
	return TablePrefix + "_index"
}

func (i *Index) SetRef(ref string) error {
	name, prop, err := ParseRef(ref)
	if err != nil {
		return err
	}

	i.Name = name

	return i.Property.Set(prop)
}

func (i *Index) Ref() string {
	return Ref(i.Name, i.Property.Get())
}

type Indexes []*Index

func (i Indexes) WriteTab(w io.Writer) (int64, error) {
	var n int64

	m, err := fmt.Fprint(w, "ID\tNAME\tPROPERTY\tVALUE NAMESPACE\tSCHEMA NAMESPACE\n")
	if err != nil {
		return int64(m), err
	}

	n += int64(m)

	for _, i := range i {
		m, err = fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			i.ID,
			i.Name,
			i.Property,
			i.ValueIndex,
			nonempty(i.SchemaIndex.String(), "-"),
		)

		n += int64(m)

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (i Indexes) WriteTo(w io.Writer) (int64, error) {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	n, err := i.WriteTab(tw)
	if err != nil {
		return n, err
	}

	if err := tw.Flush(); err != nil {
		return n, err
	}

	return n, err
}

func (i Indexes) String() string {
	var buf bytes.Buffer

	if _, err := i.WriteTo(&buf); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return buf.String()
}
