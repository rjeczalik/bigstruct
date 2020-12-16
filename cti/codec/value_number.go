package codec

import (
	"fmt"
	"reflect"
)

var _ = DefaultValue.
	Register("number", Value{
		Type:    "number",
		Convert: numberConvert,
	})

var numberType = reflect.TypeOf((*int)(nil)).Elem()

func numberConvert(v interface{}) (interface{}, error) {
	if v == nil {
		return 0, nil
	}

	if reflect.TypeOf(v).ConvertibleTo(numberType) {
		return reflect.ValueOf(v).Convert(numberType).Interface(), nil
	}

	s, err := stringConvert(v)
	if err != nil {
		return nil, err
	}

	if s == "" {
		return 0, nil
	}

	var n int
	if _, err := fmt.Sscanf(s.(string), "%v", &n); err != nil {
		return nil, err
	}

	return n, nil
}
