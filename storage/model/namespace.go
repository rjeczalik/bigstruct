package model

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/rjeczalik/bigstruct/internal/types"
)

type Namespace struct {
	Model    `yaml:",inline"`
	Name     string   `gorm:"column:name;type:tinytext;not null" yaml:"name,omitempty" json:"name,omitempty"`
	Priority int      `gorm:"column:priority;type:smallint;not null" yaml:"priority,omitempty" json:"priority,omitempty"`
	Property Property `gorm:"column:property;type:tinytext;not null" yaml:"property,omitempty" json:"property,omitempty"`
	Metadata Metadata `gorm:"column:metadata;type:text" yaml:"metadata,omityempty" json:"metadata,omitempty"`
}

func (*Namespace) TableName() string {
	return TablePrefix + "_namespace"
}

func (n *Namespace) Ref() string {
	return Ref(n.Name, n.Property.Get())
}

type Namespaces []*Namespace

func (ns Namespaces) WriteTab(w io.Writer) (int64, error) {
	var n int64

	m, err := fmt.Fprint(w, "ID\tNAME\tPROPERTY\tPRIORITY\tMETADATA\n")
	if err != nil {
		return int64(m), err
	}

	n += int64(m)

	for _, ns := range ns {
		m, err = fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%s\n",
			ns.ID,
			ns.Name,
			ns.Property,
			ns.Priority,
			nonempty(ns.Metadata.String(), "-"),
		)

		n += int64(m)

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (ns Namespaces) WriteTo(w io.Writer) (int64, error) {
	tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)

	n, err := ns.WriteTab(tw)
	if err != nil {
		return n, err
	}

	if err := tw.Flush(); err != nil {
		return n, err
	}

	return n, err

}

func (ns Namespaces) String() string {
	var buf bytes.Buffer

	if _, err := ns.WriteTo(&buf); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return buf.String()
}

type Property string

func (p Property) String() string {
	return string(p)
}

func (p Property) Get() interface{} {
	return types.YAML(p).Value()
}

func (p *Property) Set(v interface{}) error {
	switch prop := p.Get().(type) {
	case bool:
		if !prop && v != nil {
			return errors.New("property not supported")
		}

		if prop && v == nil {
			return errors.New("property required")
		}

		if v != nil {
			*p = Property(types.MakeYAML(v))
		} else {
			*p = ""
		}
	case string:
		if prop == "" && v == nil {
			return errors.New("property required")
		}

		if v != nil {
			*p = Property(types.MakeYAML(v))
		}
	case nil:
		if v == nil {
			return errors.New("property required")
		}

		*p = Property(types.MakeYAML(v))
	default:
		if prop == nil && v == nil {
			return errors.New("property required")
		}

		if v != nil {
			*p = Property(types.MakeYAML(v))
		}
	}

	return nil
}
