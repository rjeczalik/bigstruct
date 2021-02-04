package model

import (
	"github.com/rjeczalik/bigstruct/internal/types"
)

type Meta struct {
	Model   `yaml:",inline"`
	Type    string `gorm:"column:type;type:tinytext;not null:index:idx_type_ref_id" yaml:"type,omitempty" json:"type,omitempty"`
	RawText string `gorm:"column:text:type:text:not null" yaml:"text,omitempty" json:"text,omitempty"`
	RefID   uint64 `gorm:"column:ref_id;type:bigint;not null;index:idx_type_ref_id" yaml:"id,omitempty" json:"id,omitempty"`
}

func (*Meta) TableName() string {
	return Prefix + "_meta"
}

func (m *Meta) SetText(v interface{}) {
	m.RawText = types.MakeJSON(v).String()
}

func (m *Meta) Text(v interface{}) error {
	return types.JSON(m.RawText).Unmarshal(v)
}
