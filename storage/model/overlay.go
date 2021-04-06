package model

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/internal/types"
)

type Overlay struct {
	Model    `yaml:",inline"`
	Name     string `gorm:"column:name;type:tinytext;not null" yaml:"name,omitempty" json:"name,omitempty"`
	Property string `gorm:"-" yaml:"property,omitempty" json:"property,omitempty"`
	Priority int    `gorm:"column:priority;type:smallint;not null" yaml:"priority,omitempty" json:"priority,omitempty"`
	Metadata Object `gorm:"column:metadata;type:text" yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

func (*Overlay) TableName() string {
	return TablePrefix + "_overlay"
}

func (o *Overlay) Ref() string {
	if o.Property != "" {
		return o.Name + "=" + o.Property
	}
	return o.Name
}

func (o *Overlay) SetRef(ref string) error {
	name, prop, err := ParseRef(ref)
	if err != nil {
		return err
	}

	o.Name = name
	o.Property = prop

	return nil
}

func (o *Overlay) Copy() *Overlay {
	if o == nil {
		return nil
	}
	oCopy := *o
	return &oCopy
}

func (o *Overlay) Meta() *OverlayMeta {
	var om OverlayMeta

	if err := o.Metadata.Unmarshal(&om); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return &om
}

func (o *Overlay) SetProperty(prop string) error {
	switch ok := !o.Meta().NoProperty; {
	case !ok && prop != "":
		return fmt.Errorf("property %q not supported for %q overlay", prop, o.Name)
	case ok && prop == "":
		return fmt.Errorf("property required for %q overlay", o.Name)
	default:
		o.Property = prop
	}

	return nil
}

type Overlays []*Overlay

func (os Overlays) ByName(name string) *Overlay {
	for _, o := range os {
		if o.Name == name {
			return o
		}
	}
	return nil
}

func (os Overlays) ByRef(ref string) *Overlay {
	name, prop, err := ParseRef(ref)
	if err != nil {
		return nil
	}

	o := os.ByName(name)
	if o == nil {
		return nil
	}

	if err := o.SetProperty(prop); err != nil {
		return nil
	}

	return o
}

func (os Overlays) WriteTab(w io.Writer) (int64, error) {
	var n int64

	m, err := fmt.Fprint(w, "ID\tNAME\tPRIORITY\tMETADATA\n")
	if err != nil {
		return int64(m), err
	}

	n += int64(m)

	for _, o := range os {
		m, err = fmt.Fprintf(w, "%d\t%s\t%d\t%s\n",
			o.ID,
			o.Name,
			o.Priority,
			nonempty(o.Metadata.String(), "-"),
		)

		n += int64(m)

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (os Overlays) WriteTo(w io.Writer) (int64, error) {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	n, err := os.WriteTab(tw)
	if err != nil {
		return n, err
	}

	if err := tw.Flush(); err != nil {
		return n, err
	}

	return n, err
}

func (os Overlays) String() string {
	var buf bytes.Buffer

	if _, err := os.WriteTo(&buf); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return buf.String()
}

type OverlayMeta struct {
	NoProperty bool `json:"no_property,omitempty"`
	Schema     bool `json:"schema,omitempty"`
}

func (om *OverlayMeta) JSON() types.JSON {
	return types.MakeJSON(om)
}

func (om *OverlayMeta) Metadata() Object {
	return Object(om.JSON())
}
