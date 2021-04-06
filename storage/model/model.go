package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/rjeczalik/bigstruct/internal/random"

	"gorm.io/gorm"
)

var (
	TablePrefix       = "bigstruct"
	RefSeparator      = '='
	ModelBeforeCreate = RandomID
)

func RandomID(m *Model, db *gorm.DB) error {
	if m.ID == 0 {
		m.ID = random.ID()
	}

	return nil
}

type Model struct {
	ID        uint64         `gorm:"column:id;type:bigint;not null;primaryKey;autoIncrement" yaml:"id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null" yaml:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null" yaml:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index" yaml:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (m *Model) BeforeCreate(db *gorm.DB) error {
	return ModelBeforeCreate(m, db)
}

func (*Model) Options() []string {
	return []string{
		"gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",
	}
}

func Ref(name, prop string) string {
	if prop != "" {
		return name + string(RefSeparator) + prop
	}
	return name
}

func ParseRef(ref string) (name, prop string, err error) {
	parts, err := splitAlphanum(ref, string(RefSeparator))
	if err != nil {
		return "", "", fmt.Errorf("invalid ref: %w", err)
	}

	switch len(parts) {
	case 0:
		return "", "", fmt.Errorf("ref is empty or missing: %q", ref)
	case 1:
		return parts[0], "", nil
	case 2:
		return parts[0], parts[1], nil
	default:
		return "", "", fmt.Errorf("invalid ref: %q", ref)
	}
}

func nonempty(s ...string) string {
	for _, s := range s {
		if s != "" {
			return s
		}
	}
	return ""
}

func splitAlphanum(s, sep string) (sl []string, err error) {
	for _, s := range strings.Split(s, sep) {
		if s = strings.TrimSpace(s); s == "" {
			continue
		}

		for _, r := range s {
			if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '.' {
				return nil, fmt.Errorf("bad char: %q", r)
			}
		}

		sl = append(sl, s)
	}

	return sl, nil
}

func reencode(in, out interface{}) error {
	p, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(p, out)
}
