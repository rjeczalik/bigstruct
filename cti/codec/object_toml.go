package codec

import (
	"github.com/pelletier/go-toml"
)

var _ = DefaultObject.
	Register("toml", Object{
		Type:      "toml",
		Marshal:   toml.Marshal,
		Unmarshal: toml.Unmarshal,
	})
