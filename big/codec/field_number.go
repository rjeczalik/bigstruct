package codec

import (
	"errors"
	"reflect"
	"strconv"
)

var _ = DefaultField.
	Register("number", 20, Field{
		Type:    "number",
		Convert: numberConvert,
	})

var numberType = reflect.TypeOf((*int)(nil)).Elem()

func numberConvert(allowEmpty bool, v interface{}) (interface{}, error) {
	if allowEmpty && isNull(v) {
		return 0, nil
	}

	if reflect.TypeOf(v).ConvertibleTo(numberType) {
		return reflect.ValueOf(v).Convert(numberType).Interface(), nil
	}

	var err error
	if v, err = stringConvert(allowEmpty, v); err != nil {
		return nil, err
	}

	s, _ := v.(string)
	if s == "" {
		if allowEmpty {
			return 0, nil
		}

		return nil, errors.New("string is empty")
	}

	if d, err := strconv.ParseFloat(s, 64); err == nil {
		return int(d), nil
	}

	n, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return nil, err
	}

	return n, nil
}
