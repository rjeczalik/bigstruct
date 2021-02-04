package codec

import (
	"fmt"
	"reflect"
)

var _ = DefaultField.
	Register("string", Field{
		Type:    "string",
		Convert: stringConvert,
	})

var stringType = reflect.TypeOf((*string)(nil)).Elem()

func stringConvert(v interface{}) (interface{}, error) {
	if v == nil {
		return "", nil
	}

	if t := reflect.TypeOf(v); t.ConvertibleTo(stringType) && !t.ConvertibleTo(numberType) {
		return reflect.ValueOf(v).Convert(stringType).Interface(), nil
	}

	return nil, fmt.Errorf("unable to convert %T to string", v)
}
