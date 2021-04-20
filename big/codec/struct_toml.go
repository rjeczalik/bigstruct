package codec

import (
	"github.com/pelletier/go-toml"
)

var _ = DefaultStruct.
	Register("toml", 20, Struct{
		Type:      "toml",
		Marshal:   toml.Marshal,
		Unmarshal: toml.Unmarshal,
	})
