package pak

import (
	"fmt"

	"github.com/rjeczalik/bigstruct/storage"
	"github.com/rjeczalik/bigstruct/storage/model"
)

type Pak struct {
	Name     string         `yaml:"name"`
	URL      string         `yaml:"url"`
	Version  string         `yaml:"version"`
	Overlays model.Overlays `yaml:"overlays"`
	Indexes  model.Indexes  `yaml:"indexes,omitempty"`
	Schemas  model.Schemas  `yaml:"schemas,omitempty"`
	Values   model.Values   `yaml:"values,omitempty"`
}

func (pk *Pak) Store(tx storage.Gorm) error {
	for _, o := range pk.Overlays {
		if err := tx.UpsertOverlay(o); err != nil {
			return fmt.Errorf("error upserting %q overlay: %w", o.Ref(), err)
		}
	}

	for _, idx := range pk.Indexes {
		if err := tx.UpsertIndex(idx); err != nil {
			return fmt.Errorf("error upserting %q index: %w", idx.Ref(), err)
		}
	}

	if len(pk.Schemas) != 0 {
		if err := tx.UpsertSchemas(pk.Schemas); err != nil {
			return fmt.Errorf("error upserting schema: %w", err)
		}
	}

	if len(pk.Values) != 0 {
		if err := tx.UpsertValues(pk.Values); err != nil {
			return fmt.Errorf("error upserting values: %w", err)
		}
	}

	return nil
}
