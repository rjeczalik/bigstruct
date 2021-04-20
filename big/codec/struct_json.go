package codec

import (
	"encoding/json"
)

var _ = DefaultStruct.
	Register("json", 30, Struct{
		Type:      "json",
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	})
