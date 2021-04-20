package codec

import (
	"errors"
	"reflect"
	"strconv"
)

var _ = DefaultField.
	Register("bool", 30, Field{
		Type:    "bool",
		Convert: boolConvert,
	})

var boolType = reflect.TypeOf((*bool)(nil)).Elem()

func boolConvert(allowEmpty bool, v interface{}) (interface{}, error) {
	if allowEmpty && isNull(v) {
		return false, nil
	}

	if reflect.TypeOf(v).ConvertibleTo(boolType) {
		return reflect.ValueOf(v).Convert(boolType).Interface(), nil
	}

	s, err := stringConvert(allowEmpty, v)
	if err != nil {
		return nil, err
	}

	switch s {
	case "":
		return nil, errors.New("string is empty")
	case "0", "1":
		return nil, errors.New("value is a number")
	}

	b, err := strconv.ParseBool(s.(string))
	if err != nil {
		return nil, err
	}

	return b, nil
}
