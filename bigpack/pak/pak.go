package pak

import (
	"github.com/rjeczalik/bigstruct/storage/model"
)

type Pak struct {
	Name       string           `yaml:"name"`
	URL        string           `yaml:"url"`
	Version    string           `yaml:"version"`
	Namespaces model.Namespaces `yaml:"namespaces"`
	Indexes    model.Indexes    `yaml:"indexes,omitempty"`
	Schemas    model.Schemas    `yaml:"schemas,omitempty"`
	Values     model.Values     `yaml:"values,omitempty"`
}
