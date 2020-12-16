package codec

import (
	"errors"
	"reflect"
	"strconv"
)

var _ = DefaultValue.
	Register("bool", Value{
		Type:    "bool",
		Convert: boolConvert,
	})

var boolType = reflect.TypeOf((*bool)(nil)).Elem()

func boolConvert(v interface{}) (interface{}, error) {
	if v == nil {
		return false, nil
	}

	if reflect.TypeOf(v).ConvertibleTo(boolType) {
		return reflect.ValueOf(v).Convert(boolType).Interface(), nil
	}

	s, err := stringConvert(v)
	if err != nil {
		return nil, err
	}

	switch s {
	case "":
		return false, nil
	case "0", "1":
		return nil, errors.New("value is a number")
	}

	b, err := strconv.ParseBool(s.(string))
	if err != nil {
		return nil, err
	}

	return b, nil
}
