package codec

import (
	"fmt"
	"reflect"
)

var _ = DefaultField.
	Register("string", 10, Field{
		Type:    "string",
		Convert: stringConvert,
	})

var stringType = reflect.TypeOf((*string)(nil)).Elem()

func stringConvert(_ bool, v interface{}) (interface{}, error) {
	if isNull(v) {
		return "", nil
	}

	if t := reflect.TypeOf(v); t.ConvertibleTo(stringType) && !t.ConvertibleTo(numberType) {
		return reflect.ValueOf(v).Convert(stringType).Interface(), nil
	}

	return nil, fmt.Errorf("unable to convert %T to string", v)
}
