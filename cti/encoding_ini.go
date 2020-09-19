package cti

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var ini = encini{
	key: regexp.MustCompile(`^[a-zA-Z_0-9]+$`),
}

type encini struct {
	key *regexp.Regexp
}

func (encini) Marshal(v interface{}) ([]byte, error) {
	obj, keys, err := toobj(v)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	for _, k := range keys {
		fmt.Fprintf(&buf, "%s=%q\n", k, fmt.Sprint(obj[k]))
	}

	return buf.Bytes(), nil
}

func (e encini) Unmarshal(p []byte, v interface{}) error {
	s := bufio.NewScanner(bytes.NewReader(p))

	obj := make(map[string]string)

	for i := 0; s.Scan(); i++ {
		line := strings.TrimSpace(s.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		j := strings.IndexRune(line, '=')
		if j == -1 {
			return fmt.Errorf("unable to parse line %d: %q", i, line)
		}

		k := strings.TrimSpace(line[:j])
		v := strings.TrimSpace(line[j+1:])

		if !e.key.MatchString(k) {
			return fmt.Errorf("invalid key at line %d: %q", i, k)
		}

		if s, err := strconv.Unquote(v); err == nil {
			v = s
		}

		obj[k] = v
	}

	return reencode(obj, v)
}

func toobj(v interface{}) (m map[string]interface{}, keys []string, err error) {
	switch v := v.(type) {
	case map[string]interface{}:
		m = v
	case map[interface{}]interface{}:
		m = make(map[string]interface{}, len(v))

		for k, v := range v {
			m[fmt.Sprint(k)] = v
		}
	default:
		if err := reencode(v, &m); err != nil {
			return nil, nil, err
		}

		if len(m) == 0 {
			return nil, nil, fmt.Errorf("value is neither struct nor non-empty map: %T", v)
		}
	}

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return m, keys, nil
}

func reencode(in, out interface{}) error {
	p, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(p, out)
}
