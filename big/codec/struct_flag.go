package codec

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var _ = DefaultStruct.
	Register("flag", Struct{
		Type:      "flag",
		Marshal:   flag.Marshal,
		Unmarshal: flag.Unmarshal,
	})

var flag flagCodec

type flagCodec struct{}

func (flagCodec) Marshal(v interface{}) ([]byte, error) {
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
			if s := fmt.Sprint(v); s != "" {
				fmt.Fprintf(&buf, "%s %s ", k, s)
			} else {
				fmt.Fprintf(&buf, "%s %q ", k, s)
			}
		}
	}

	return bytes.TrimSpace(buf.Bytes()), nil
}

func (flagCodec) Unmarshal(p []byte, v interface{}) error {
	var (
		args = split(string(p), " ") // fixme: support escaped values
		obj  = make(map[string]string)
	)

	for i := 0; i < len(args); {
		arg := args[i]

		if !strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unexpected argument: %q", arg)
		}

		i = i + 1

		if j := strings.IndexRune(arg, '='); j != -1 {
			if s, err := strconv.Unquote(arg[j+1:]); err == nil {
				obj[arg[:j]] = s
			} else {
				obj[arg[:j]] = arg[j+1:]
			}
		} else {
			obj[arg] = ""
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
