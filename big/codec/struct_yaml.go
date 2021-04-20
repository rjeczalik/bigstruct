package codec

import (
	"gopkg.in/yaml.v3"
)

var _ = DefaultStruct.
	Register("yaml", 10, Struct{
		Type:      "yaml",
		Marshal:   yaml.Marshal,
		Unmarshal: yaml.Unmarshal,
	})
