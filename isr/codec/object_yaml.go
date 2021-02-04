package codec

import (
	"gopkg.in/yaml.v3"
)

var _ = DefaultObject.
	Register("yaml", Object{
		Type:      "yaml",
		Marshal:   yaml.Marshal,
		Unmarshal: yaml.Unmarshal,
	})
