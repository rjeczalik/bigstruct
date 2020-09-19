package cti

import (
	"bytes"
	"fmt"
	"strings"
)

var flag encflag

type encflag struct{}

func (encflag) Marshal(v interface{}) ([]byte, error) {
	obj, keys, err := toobj(v)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	for _, k := range keys {
		switch v := obj[k].(type) {
		case bool:
			if v {
				fmt.Fprintf(&buf, "%s ", k)
			} else {
				fmt.Fprintf(&buf, "%s=%t ", k, v)
			}
		case nil:
			fmt.Fprintf(&buf, "%s ", k)
		default:
			fmt.Fprintf(&buf, "%s=%q ", k, fmt.Sprint(v))
		}
	}

	return bytes.TrimSpace(buf.Bytes()), nil
}

func (encflag) Unmarshal(p []byte, v interface{}) error {
	var (
		args = split(string(p), " ") // fixme: support escaped values
		obj  = make(map[string]interface{})
	)

	for i := 0; i < len(args); {
		arg := args[i]

		if !strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unexpected argument: %q", arg)
		}

		i = i + 1

		if j := strings.IndexRune(arg, '='); j != -1 {
			obj[arg[:j]] = arg[j+1:]
		} else {
			obj[arg] = nil
		}

		if i < len(args) && !strings.HasPrefix(args[i], "-") {
			obj[arg] = args[i] // decode?
			i = i + 1
		}
	}

	return reencode(obj, v)
}

func split(s, sep string) (r []string) {
	for _, s := range strings.Split(s, sep) {
		if s = strings.TrimSpace(s); s != "" {
			r = append(r, s)
		}
	}
	return r
}
