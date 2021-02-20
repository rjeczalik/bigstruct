package codec

import (
	"encoding/json"
)

var _ = DefaultStruct.
	Register("json", Struct{
		Type:      "json",
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	})
