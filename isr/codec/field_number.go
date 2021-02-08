package codec

import (
	"reflect"
	"strconv"
)

var _ = DefaultField.
	Register("number", Field{
		Type:    "number",
		Convert: numberConvert,
	})

var numberType = reflect.TypeOf((*int)(nil)).Elem()

func numberConvert(v interface{}) (interface{}, error) {
	if isNull(v) {
		return 0, nil
	}

	if reflect.TypeOf(v).ConvertibleTo(numberType) {
		return reflect.ValueOf(v).Convert(numberType).Interface(), nil
	}

	var err error
	if v, err = stringConvert(v); err != nil {
		return nil, err
	}

	s, _ := v.(string)
	if s == "" {
		return 0, nil
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
