package model

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/big"
	"github.com/rjeczalik/bigstruct/internal/types"
)

type Value struct {
	Model           `yaml:",inline"`
	Overlay         *Overlay `gorm:"" yaml:"-" json:"-"`
	OverlayID       uint64   `gorm:"column:overlay_id;type:bigint;not null;index" yaml:"overlay_id,omitempty" json:"overlay_id,omitempty"`
	OverlayProperty string   `gorm:"column:overlay_property;type:tinytext;not null" yaml:"overlay_property,omitempty" json:"overlay_property,omitempty"`
	Key             string   `gorm:"column:key;type:text;not null" yaml:"key,omitempty" json:"key,omitempty"`
	RawValue        string   `gorm:"column:value;type:text" yaml:"value,omitempty" json:"value,omitempty"`
	Metadata        Object   `gorm:"column:metadata;type:text" yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

func (*Value) TableName() string {
	return TablePrefix + "_value"
}

func (v *Value) SetValue(w interface{}) {
	v.RawValue = types.MakeYAML(w).String()
}

func (v *Value) Value() interface{} {
	return types.YAML(v.RawValue).Value()
}

type Values []*Value

func MakeValues(o *Overlay, f big.Fields) Values {
	values := make(Values, 0, len(f))

	for _, f := range f {
		if f.Value == nil {
			continue // skip empty entries, they will get recreated from the tree either way
		}

		v := &Value{
			Key: f.Key,
		}

		if o != nil {
			v.Overlay = o
			v.OverlayID = o.ID
			v.OverlayProperty = o.Property
		}

		if f.Value != big.NoValue {
			v.RawValue = types.MakeYAML(f.Value).String()
		}

		values = append(values, v)
	}

	return values
}

func (v Values) SetOverlay(o *Overlay) {
	for _, v := range v {
		v.Overlay = o
		v.OverlayID = o.ID
		v.OverlayProperty = o.Property
	}
}

func (v Values) SetMeta(meta Object) {
	for _, v := range v {
		v.Metadata = meta
	}
}

func (v Values) Fields() big.Fields {
	f := make(big.Fields, 0, len(v))

	for _, v := range v {
		f = append(f, big.Field{
			Key:   v.Key,
			Value: v.Value(),
		})
	}

	return f
}

func (v Values) WriteTab(w io.Writer) (int64, error) {
	var n int64

	m, err := fmt.Fprint(w, "ID\tOVERLAY\tKEY\tVALUE\tMETADATA\n")
	if err != nil {
		return int64(m), err
	}

	n += int64(m)

	for _, v := range v {
		m, err = fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			v.ID,
			v.Overlay.Ref(),
			v.Key,
			nonempty(v.RawValue, "-"),
			nonempty(v.Metadata.String(), "-"),
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
