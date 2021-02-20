package model

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/rjeczalik/bigstruct/internal/types"
	"gopkg.in/yaml.v3"
)

type Object string

var (
	_ json.Marshaler   = Object("")
	_ json.Unmarshaler = (*Object)(nil)
	_ yaml.Marshaler   = Object("")
	_ yaml.Unmarshaler = (*Object)(nil)
)

func (obj *Object) Set(v interface{}) Object {
	var o types.Object

	switch v := v.(type) {
	case types.Object:
		o = v
	case map[string]interface{}:
		o = v
	default:
		if err := reencode(v, &o); err != nil {
			panic("unexpected error: " + err.Error())
		}
	}

	if len(o) != 0 {
		*obj = Object(o.JSON())
	} else {
		*obj = ""
	}

	return *obj
}

func (obj *Object) SetValues(kv ...interface{}) Object {
	if len(kv)%2 != 0 {
		panic("odd number of arguments")
	}

	o := make(types.Object, len(kv)/2)

	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprint(kv[i])
		v := kv[i+1]

		o[k] = v
	}

	return obj.Set(o)
}

func (obj Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(obj.Map())
}

func (obj *Object) UnmarshalJSON(p []byte) error {
	var o types.Object

	if err := json.Unmarshal(p, &o); err != nil {
		return err
	}

	obj.Set(o)

	return nil
}

func (obj Object) MarshalYAML() (interface{}, error) {
	return obj.Map(), nil
}

func (obj *Object) UnmarshalYAML(n *yaml.Node) error {
	var o types.Object

	if err := n.Decode(&o); err != nil {
		return err
	}

	obj.Set(o)

	return nil
}

func (obj *Object) Merge(v map[string]interface{}) Object {
	return obj.Set(types.Object(obj.Map()).Merge(types.Object(v)).Map())
}

func (obj Object) Map() map[string]interface{} {
	return types.JSON(obj).Object().Map()
}

func (obj Object) Unmarshal(v interface{}) error {
	return types.JSON(obj).Unmarshal(v)
}

func (obj *Object) Update(ubj Object) bool {
	if old := *obj; !obj.Merge(ubj.Map()).Equal(old) {
		return true
	}
	return false
}

func (obj Object) Equal(ubj Object) bool {
	return reflect.DeepEqual(obj.Map(), obj.Map())
}

func (obj Object) String() string {
	return string(obj)
}
