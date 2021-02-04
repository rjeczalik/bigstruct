package model

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"reflect"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/internal/types"
)

type Index struct {
	Model       `yaml:",inline"`
	Name        string         `gorm:"column:name;type:tinytext;not null" yaml:"name,omitempty" json:"name,omitempty"`
	Property    string         `gorm:"column:property;type:tinytext" yaml:"property,omitempty" json:"property,omitempty"`
	ValueIndex  NamespaceIndex `gorm:"column:value_index;type:text;not null" yaml:"value_index,omitempty" json:"value_index,omitempty"`
	SchemaIndex NamespaceIndex `gorm:"column:schema_index;type:text;not null" yaml:"schema_index,omitempty" json:"schema_index,omitempty"`
}

func (*Index) TableName() string {
	return Prefix + "_index"
}

func (i *Index) Prefix() string {
	return path.Join("/", i.Name, i.Property)
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

type NamespaceIndex string

func (ni *NamespaceIndex) Set(m map[string]string) NamespaceIndex {
	if len(m) != 0 {
		*ni = NamespaceIndex(types.MakeJSON(m).String())
	} else {
		*ni = ""
	}
	return *ni
}

func (ni NamespaceIndex) Get() (m map[string]string) {
	if err := types.JSON(ni).Unmarshal(&m); err != nil {
		panic("unexpected error: " + err.Error())
	}
	return m
}

func (ni NamespaceIndex) Equal(mi NamespaceIndex) bool {
	return reflect.DeepEqual(ni.Get(), mi.Get())
}

func (ni *NamespaceIndex) Merge(m map[string]string) NamespaceIndex {
	return ni.Set(types.KV(ni.Get()).Merge(types.KV(m)).Map())
}

func (ni NamespaceIndex) String() string {
	return string(ni)
}
