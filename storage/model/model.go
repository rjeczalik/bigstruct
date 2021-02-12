package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rjeczalik/bigstruct/internal/random"
	"github.com/rjeczalik/bigstruct/internal/types"

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

type Metadata string

func (m *Metadata) Set(v interface{}) {
	*m = Metadata(types.MakeJSON(v))
}

func (m Metadata) Unmarshal(v interface{}) error {
	return types.JSON(m).Unmarshal(v)
}

func (m Metadata) String() string {
	return string(m)
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

func Ref(name string, prop interface{}) string {
	switch v := prop.(type) {
	case bool:
		return name
	case nil:
		return name
	case string:
		return name + string(RefSeparator) + v
	default:
		return name + string(RefSeparator) + fmt.Sprint(v)
	}
}

func ParseRef(ref string) (name string, prop interface{}, err error) {
	switch parts := strings.Split(ref, string(RefSeparator)); len(parts) {
	case 0:
		return "", nil, errors.New("ref is empty or missing")
	case 1:
		return parts[0], nil, nil
	case 2:
		return parts[0], types.YAML(parts[1]).Value(), nil
	default:
		return "", nil, fmt.Errorf("invalid ref: %q", name)
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

func reencode(in, out interface{}) error {
	p, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(p, out)
}
