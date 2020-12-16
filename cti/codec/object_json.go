package codec

import (
	"encoding/json"
)

var _ = DefaultObject.
	Register("json", Object{
		Type:      "json",
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	})
