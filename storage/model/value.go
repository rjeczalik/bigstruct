package model

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/internal/types"
	"github.com/rjeczalik/bigstruct/isr"
)

type Value struct {
	Model       `yaml:",inline"`
	Namespace   *Namespace `gorm:"" yaml:"-" json:"-"`
	NamespaceID uint64     `gorm:"column:namespace_id;type:bigint;not null;index" yaml:"namespace_id,omitempty" json:"namespace_id,omitempty"`
	Key         string     `gorm:"column:key;type:text;not null" yaml:"key,omitempty" json:"key,omitempty"`
	RawValue    string     `gorm:"column:value;type:text" yaml:"value,omitempty" json:"value,omitempty"`
}

func (*Value) TableName() string {
	return Prefix + "_value"
}

func (v *Value) SetValue(w interface{}) {
	v.RawValue = types.MakeYAML(w).String()
}

func (v *Value) Value() interface{} {
	return types.YAML(v.RawValue).Value()
}

type Values []*Value

func MakeValues(ns *Namespace, f isr.Fields) Values {
	values := make(Values, 0, len(f))

	for _, f := range f {
		if f.Value == nil {
			continue // skip empty entries, they will get recreated from the tree either way
		}

		v := &Value{
			Key:         f.Key,
			Namespace:   ns,
			NamespaceID: ns.ID,
		}

		if f.Value != isr.NoValue {
			v.RawValue = types.MakeYAML(f.Value).String()
		}

		values = append(values, v)
	}

	return values
}

func (v Values) Fields() isr.Fields {
	f := make(isr.Fields, 0, len(v))

	for _, v := range v {
		f = append(f, isr.Field{
			Key:   v.Key,
			Value: v.Value(),
		})
	}

	return f
}

func (v Values) WriteTab(w io.Writer) (int64, error) {
	var n int64

	m, err := fmt.Fprint(w, "ID\tNAMESPACE\tKEY\tVALUE\n")
	if err != nil {
		return int64(m), err
	}

	n += int64(m)

	for _, v := range v {
		m, err = fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			v.ID,
			v.Namespace.Namespace(),
			v.Key,
			nonempty(v.RawValue, "-"),
		)

		n += int64(m)

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (v Values) WriteTo(w io.Writer) (int64, error) {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	n, err := v.WriteTab(tw)
	if err != nil {
		return n, err
	}

	if err := tw.Flush(); err != nil {
		return n, err
	}

	return n, err

}

func (v Values) String() string {
	var buf bytes.Buffer

	if _, err := v.WriteTo(&buf); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return buf.String()
}
