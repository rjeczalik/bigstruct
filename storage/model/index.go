package model

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/internal/types"
)

type Index struct {
	Model       `yaml:",inline"`
	Name        string         `gorm:"column:name;type:tinytext;not null" yaml:"name,omitempty" json:"name,omitempty"`
	Property    Property       `gorm:"column:property;type:tinytext" yaml:"property,omitempty" json:"property,omitempty"`
	ValueIndex  NamespaceIndex `gorm:"column:value_index;type:text;not null" yaml:"value_index,omitempty" json:"value_index,omitempty"`
	SchemaIndex NamespaceIndex `gorm:"column:schema_index;type:text;not null" yaml:"schema_index,omitempty" json:"schema_index,omitempty"`
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

type NamespaceIndex string

func (ni *NamespaceIndex) SetValue(v interface{}) NamespaceIndex {
	var o types.Object

	switch v := v.(type) {
	case types.Object:
		o = v
	case map[string]interface{}:
		o = v
	default:
		if err := reencode(v, &o); err != nil {
			panic("unexpected error: " + err.Error())
		}
	}

	if len(o) != 0 {
		*ni = NamespaceIndex(o.JSON().String())
	} else {
		*ni = ""
	}

	return *ni
}

func (ni *NamespaceIndex) Set(kv ...interface{}) NamespaceIndex {
	if len(kv)%2 != 0 {
		panic("odd number of arguments")
	}

	o := make(types.Object, len(kv)/2)

	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprint(kv[i])
		v := kv[i+1]

		o[k] = v
	}

	return ni.SetValue(o)
}

func (ni NamespaceIndex) Get() (v map[string]interface{}) {
	if err := types.JSON(ni).Unmarshal(&v); err != nil {
		panic("unexpected error: " + err.Error())
	}
	return v
}

func (ni NamespaceIndex) Equal(mi NamespaceIndex) bool {
	return reflect.DeepEqual(ni.Get(), mi.Get())
}

func (ni *NamespaceIndex) Merge(v map[string]interface{}) NamespaceIndex {
	return ni.SetValue(types.Object(ni.Get()).Merge(types.Object(v)).Map())
}

func (ni NamespaceIndex) String() string {
	return string(ni)
}
